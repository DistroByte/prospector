package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"

	nomad "github.com/hashicorp/nomad/nomad/structs"
)

func CreateJobFromTemplate(project Project, jobSource string) (int, error) {
	// escape double quotes
	jobSource = strings.ReplaceAll(jobSource, `"`, `\"`)

	// create job spec from template using project
	jobSpec := new(bytes.Buffer)
	if err := template.Must(template.New("job").Parse(jobSource)).Execute(jobSpec, project); err != nil {
		return 500, err
	}

	// create parse body
	parseBody := new(bytes.Buffer)
	parseBodyTemplate := template.Must(template.New("parseBody").Parse(`{ "JobHCL": "{{ . }}", "Canonicalize": true }`))
	if err := parseBodyTemplate.Execute(parseBody, jobSpec.String()); err != nil {
		return 500, err
	}

	parseBodyCleaned := strings.ReplaceAll(parseBody.String(), "\n", `\n`)
	parseBodyCleaned = strings.ReplaceAll(parseBodyCleaned, "\t", ` `)

	// parse job against nomad
	parseResponse, err := http.Post("http://zeus.internal:4646/v1/jobs/parse", "application/json", strings.NewReader(parseBodyCleaned))
	if err != nil {
		return 500, err
	}

	// process parse response
	var parseResponseBuffer bytes.Buffer
	if _, err := io.Copy(&parseResponseBuffer, parseResponse.Body); err != nil {
		return 500, err
	}

	// check for error
	if parseResponse.StatusCode != http.StatusOK {
		return 500, fmt.Errorf("error parsing job: %s", parseResponseBuffer.String())
	}

	// create job run body
	jobRunTemplate := template.Must(template.New("jobRun").Parse(`{ "Job": {{ . }} }`))

	// create job run body
	var jobRun bytes.Buffer
	if err := jobRunTemplate.Execute(&jobRun, parseResponseBuffer.String()); err != nil {
		return 500, err
	}

	// run job against nomad
	jobRunResponse, err := http.Post("http://zeus.internal:4646/v1/jobs", "application/json", &jobRun)
	if err != nil {
		return 500, err
	}

	// process job run response
	var jobRegisterResponse nomad.JobRegisterResponse
	if err := json.NewDecoder(jobRunResponse.Body).Decode(&jobRegisterResponse); err != nil {
		return 500, err
	}

	return http.StatusOK, nil
}

func WriteTextFilesForVM(job Project) error {
	var cloudInitMetaData = `instance-id: prospector/{{ .Name }}
local-hostname: {{ .Name }}`

	var cloudInitUserData = `#cloud-config
groups:
  - admingroup
  - cloud-users

users:
  - default
  - name: {{ (index .Components 0).UserConfig.User }}
    groups: sudo
    shell: /bin/bash
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    ssh_authorized_keys:
      - {{ (index .Components 0).UserConfig.SSHKey }}

package_update: true
package_upgrade: true

password: {{ (index .Components 0).UserConfig.User }}{{ (index .Components 0).UserConfig.SSHKey }}
chpasswd:
  expire: False`

	userFilePath := fmt.Sprintf("./vm-config/%s-vm/user-data", job.Name)
	metaFilePath := fmt.Sprintf("./vm-config/%s-vm/meta-data", job.Name)

	userFileDir := fmt.Sprintf("./vm-config/%s-vm", job.Name)
	metaFileDir := fmt.Sprintf("./vm-config/%s-vm", job.Name)

	_, err := makesDirsAndFiles(userFilePath, userFileDir)
	if err != nil {
		return err
	}

	_, err = makesDirsAndFiles(metaFilePath, metaFileDir)
	if err != nil {
		return err
	}

	userFileDest, err := os.Create(userFilePath)
	if err != nil {
		return err
	}

	defer userFileDest.Close()

	metaFileDest, err := os.Create(metaFilePath)
	if err != nil {
		return err
	}

	defer metaFileDest.Close()

	cloudInitMetaDataTemplate, err := template.New("cloudInitMetaData").Parse(cloudInitMetaData)
	if err != nil {
		return err
	}

	var cloudInitMetaDataBuffer bytes.Buffer

	if err := cloudInitMetaDataTemplate.Execute(&cloudInitMetaDataBuffer, job); err != nil {
		return err
	}

	cloudInitUserDataTemplate, err := template.New("cloudInitUserData").Parse(cloudInitUserData)
	if err != nil {
		return err
	}

	var cloudInitUserDataBuffer bytes.Buffer

	if err := cloudInitUserDataTemplate.Execute(&cloudInitUserDataBuffer, job); err != nil {
		return err
	}

	userFileDest.WriteString(cloudInitUserDataBuffer.String())
	metaFileDest.WriteString(cloudInitMetaDataBuffer.String())

	return nil
}

func makesDirsAndFiles(metaFilePath string, metaFileDir string) (bool, error) {
	if _, err := os.Stat(metaFilePath); os.IsNotExist(err) {
		err := os.MkdirAll(metaFileDir, os.ModePerm)

		if err != nil {
			return true, err
		}
	}
	return false, nil
}
