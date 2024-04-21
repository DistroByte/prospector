package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"text/template"

	nomad "github.com/hashicorp/nomad/nomad/structs"
)

func last(i int, slice interface{}) bool {
	v := reflect.ValueOf(slice)
	return i == v.Len()-1
}

func escapeQuotes(s *bytes.Buffer) string {
	return strings.ReplaceAll(s.String(), "\"", "\\\"")
}

func CreateJobFromTemplate(project Project, jobSource string) (int, error) {
	t, err := template.New("").Funcs(template.FuncMap{
		"last":         last,
		"json":         project.ToJson,
		"escapeQuotes": escapeQuotes,
	}).Parse(jobSource)
	if err != nil {
		return 500, err
	}

	body := &bytes.Buffer{}
	err = t.Execute(body, project)
	if err != nil {
		return 500, err
	}

	println(body.String())

	// run job against nomad
	data, err := http.Post("http://zeus.internal:4646/v1/jobs", "application/json", body)
	if err != nil {
		return 500, err
	}

	body = &bytes.Buffer{}
	_, err = io.Copy(body, data.Body)
	if err != nil {
		return 500, err
	}

	var resp nomad.JobRegisterResponse
	err = json.Unmarshal(body.Bytes(), &resp)
	if err != nil {
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
  expire: False

mounts:
  - [dionysus.internal:/volume1/data/prospector/{{ (index .Components 0).UserConfig.User }}, /mnt/user-storage, nfs, "auto,nofail,noatime,nolock,intr,tcp,actimeo=1800", "0", "0"]
`

	userFilePath := fmt.Sprintf("./vm-config/%s-%s-vm/user-data", job.User, job.Components[0].Name)
	metaFilePath := fmt.Sprintf("./vm-config/%s-%s-vm/meta-data", job.User, job.Components[0].Name)

	userFileDir := fmt.Sprintf("./vm-config/%s-%s-vm", job.User, job.Components[0].Name)
	metaFileDir := fmt.Sprintf("./vm-config/%s-%s-vm", job.User, job.Components[0].Name)

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

func ConvertJobToProject(job nomad.Job) (Project, error) {
	var project Project

	meta := job.Meta
	if meta == nil {
		return project, fmt.Errorf("job has no meta")
	}

	jobDefinition, ok := meta["job-definition"]
	if !ok {
		return project, fmt.Errorf("job has no job-definition")
	}

	err := json.Unmarshal([]byte(jobDefinition), &project)
	if err != nil {
		return project, err
	}

	return project, nil
}
