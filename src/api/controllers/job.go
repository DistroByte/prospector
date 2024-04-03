package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"text/template"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	nomad "github.com/hashicorp/nomad/nomad/structs"
)

var dockerSource = `job "{{ .User }}-{{ .Name }}-prospector" {
	datacenters = ["dc1"]
	type = "service"
	
	{{ range .Components }}
	group "{{ .Name }}" {
		count = 1

		network {
			port "port" {
				to = {{ .Network.Port }}
			}
		}
		
		task "{{ .Name }}" {
			driver = "docker"
			
			config {
				image = "{{ .Image }}"
				ports = ["port"]
			}

			resources {
				cpu    = {{ .Resources.Cpu }}
				memory = {{ .Resources.Memory }}
			}

			service {
				name = "{{ .Name }}"
				port = "port"

				check {
					name = "{{ .Name }}-health"
					type = "http"
					path = "/"
					interval = "10s"
					timeout = "2s"
				}

				{{ if .Network.Expose }}
				tags = [
					"traefik.enable=true",
					"traefik.http.routers.{{ .Name }}-{{ .UserConfig.User }}-prospector.rule=Host(` + "`" + `{{ .Name }}-{{ .UserConfig.User }}.prospector.ie` + "`" + `)",
					"traefik.http.routers.{{ .Name }}-{{ .UserConfig.User }}-prospector.entrypoints=websecure",
					"traefik.http.routers.{{ .Name }}-{{ .UserConfig.User }}-prospector.tls=true",
					"traefik.http.routers.{{ .Name }}-{{ .UserConfig.User }}-prospector.tls.certresolver=lets-encrypt"
				]
				{{ end }}
			}
		}
		{{ end }}
	}
}
`

//lint:ignore U1000 Unused template for now
var vmSource = `job "{{ .User }}-{{ .Name }}-prospector" {
  datacenters = ["dc1"]

  group "{{ .User }}-{{ .Name }}" {

    network {
      mode = "host"
    }

    service {
      name = "{{ .Name }}-vm"
    }

    task "{{ .Name }}" {
      constraint {
        attribute = "${attr.unique.hostname}"
        value     = "hermes"
      }

      resources {
        cpu    = {{ .Cpu }}
        memory = {{ .Memory }}
      }

      artifact {
        source      = "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-genericcloud-amd64.qcow2"
        destination = "local/{{ .Name }}-vm.qcow2"
        mode        = "file"
      }

      driver = "qemu"

      config {
        image_path = "local/{{ .Name }}-vm.qcow2"
        accelerator = "kvm"
        drive_interface = "virtio"

        args = [
          "-netdev",
          "bridge,id=hn0",
          "-device",
          "virtio-net-pci,netdev=hn0,id=nic1,mac={{ .Mac }}",
          "-smbios",
          "type=1,serial=ds=nocloud-net;s=https://prospector.ie/api/vm-config/{{ .Name }}-vm/",
        ]
      }
    }
  }
}
`

// GetJobs gets all the jobs from nomad that have the word "prospector" in their name
//
//	@Summary		Get all jobs
//	@Description	Get all jobs from nomad
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/jobs [get]
//	@Param			long	query	boolean	false	"Get long job details"
//	@Param			running	query	boolean	false	"Get running jobs"
//	@Code			204 "No jobs found"
func (c *Controller) GetJobs(ctx *gin.Context) {
	claims := jwt.ExtractClaims(ctx)
	ctx.Set(c.IdentityKey, claims[c.IdentityKey])

	println(claims[c.IdentityKey].(string))

	data, err := c.Client.Get("/jobs?meta=true")
	if err != nil {
		ctx.Error(err)
	}

	var jobs []nomad.JobListStub = []nomad.JobListStub{}
	err = json.Unmarshal(data, &jobs)
	if err != nil {
		ctx.Error(err)
	}

	var filteredJobs []nomad.JobListStub = []nomad.JobListStub{}
	for _, job := range jobs {
		if strings.Contains(job.Name, "-prospector") && strings.Contains(job.Name, claims[c.IdentityKey].(string)) {
			filteredJobs = append(filteredJobs, job)
		}
	}

	var jobSummaries []ShortJob
	for _, job := range filteredJobs {
		jobSummaries = append(jobSummaries, ShortJob{
			ID:     job.ID,
			Status: job.Status,
		})
	}

	var runningJobs []ShortJob = []ShortJob{}
	for _, job := range jobSummaries {
		if job.Status == "running" {
			runningJobs = append(runningJobs, job)
		}
	}

	if ctx.Query("long") == "true" {
		if len(filteredJobs) == 0 {
			ctx.JSON(http.StatusNoContent, gin.H{"message": "No jobs found"})
			return
		}

		ctx.JSON(http.StatusOK, filteredJobs)
		return
	} else if ctx.Query("running") == "true" {
		if len(runningJobs) == 0 {
			ctx.JSON(http.StatusNoContent, gin.H{"message": "No running jobs found"})
			return
		}

		ctx.JSON(http.StatusOK, runningJobs)
		return
	} else {
		if len(jobSummaries) == 0 {
			ctx.JSON(http.StatusNoContent, gin.H{"message": "No jobs found"})
			return
		}

		ctx.JSON(http.StatusOK, jobSummaries)
		return
	}
}

// GetJob gets a job from nomad
//
//	@Summary		Get a job
//	@Description	Get a job from nomad
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/jobs/{id} [get]
//	@Param			id	path	string	true	"Job ID"
func (c *Controller) GetJob(ctx *gin.Context) {
	id := ctx.Param("id")

	data, err := c.Client.Get("/job/" + id)
	if err != nil {
		ctx.Error(err)
	}

	var job nomad.Job
	err = json.Unmarshal(data, &job)
	if err != nil {
		ctx.Error(err)
	}

	ctx.JSON(http.StatusOK, job)
}

// CreateJob creates a container in nomad
//
//	@Summary		Create a job in nomad
//	@Description	Create and submit a job for nomad to deploy
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			job	body		Job	true	"Job"
//	@Success		200	{object}	Message
//	@Router			/v1/jobs [post]
func (c *Controller) CreateJob(ctx *gin.Context) {
	var job Job

	if err := ctx.BindJSON(&job); err != nil {
		println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var res int
	var err error

	// generate random mac address

	claims := jwt.ExtractClaims(ctx)
	job.User = claims[c.IdentityKey].(string)

	for i := 0; i < len(job.Components); i++ {
		job.Components[i].Network.Mac = fmt.Sprintf("52:54:00:%02x:%02x:%02x", byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)))
		job.Components[i].UserConfig.User = job.User
	}

	if job.Type == "docker" {
		res, err = createJob(job, dockerSource)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if job.Type == "vm" {
		err = writeTextFiles(job)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		res, err = createJob(job, vmSource)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job type"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	println(res)

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Job submitted successfully"})
}

// DeleteJob deletes a job from nomad
//
//	@Summary		Delete a job
//	@Description	Delete a job from nomad
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	Message
//	@Router			/v1/jobs/{id} [delete]
//	@Param			id		path	string	true	"Job ID"
//	@Param			purge	query	bool	false	"Purge job"
func (c *Controller) DeleteJob(ctx *gin.Context) {
	id := ctx.Param("id")

	purge := ctx.Query("purge")

	if !strings.Contains(id, "-prospector") || id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	data, err := c.Client.Delete("/job/" + id + "?purge=" + purge)
	if err != nil {
		ctx.Error(err)
	}

	var message Message
	err = json.Unmarshal(data, &message)
	if err != nil {
		ctx.Error(err)
	}

	ctx.JSON(http.StatusOK, message)
}

// RestartJob restarts a job in nomad
//
//	@Summary		Restart a job
//	@Description	Restart a job in nomad
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	Message
//	@Router			/v1/jobs/{id}/restart [put]
//	@Param			id	path	string	true	"Job ID"
func (c *Controller) RestartJob(ctx *gin.Context) {
	id := ctx.Param("id")

	if !strings.Contains(id, "-prospector") || id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	alloc, err := c.parseRunningAllocs(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	body := bytes.NewBuffer([]byte{})

	data, err := c.Client.Post("/client/allocation/"+alloc.ID+"/restart", body)
	if err != nil {
		println(err.Error())
		ctx.Error(err)
	}

	var response nomad.GenericResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job restarted successfully"})
}

func (c *Controller) parseRunningAllocs(jobId string) (*nomad.AllocListStub, error) {
	data, err := c.Client.Get("/job/" + jobId + "/allocations")
	if err != nil {
		return nil, err
	}

	var allocations []nomad.AllocListStub
	err = json.Unmarshal(data, &allocations)
	if err != nil {
		return nil, err
	}

	for _, alloc := range allocations {
		if alloc.ClientStatus == "running" || alloc.ClientStatus == "pending" {
			return &alloc, nil
		}
	}

	return nil, gin.Error{Err: fmt.Errorf("no running allocations found"), Meta: 404}
}

func writeTextFiles(job Job) error {
	var cloudInitMetaData = `instance-id: prospector/{{ .Name }}
local-hostname: {{ .Name }}`

	var cloudInitUserData = `#cloud-config
groups:
  - admingroup
  - cloud-users

users:
  - default
  - name: {{ .UserConfig.User }}
    groups: sudo
    shell: /bin/bash
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    ssh_authorized_keys:
      - {{ .UserConfig.SSHKey }}

package_update: true
package_upgrade: true

password: {{ .UserConfig.User }}{{ .UserConfig.SSHKey }}
chpasswd:
  expire: False`

	userFilePath := fmt.Sprintf("./vm-config/%s-vm/user-data", job.Name)
	metaFilePath := fmt.Sprintf("./vm-config/%s-vm/meta-data", job.Name)

	userFileDir := fmt.Sprintf("./vm-config/%s-vm", job.Name)
	metaFileDir := fmt.Sprintf("./vm-config/%s-vm", job.Name)

	if _, err := os.Stat(userFilePath); os.IsNotExist(err) {
		err := os.MkdirAll(userFileDir, os.ModePerm)

		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(metaFilePath); os.IsNotExist(err) {
		err := os.MkdirAll(metaFileDir, os.ModePerm)

		if err != nil {
			return err
		}
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

func createJob(job Job, jobSource string) (int, error) {
	jobSource = strings.ReplaceAll(jobSource, `"`, `\"`)

	jobTemplate, err := template.New("job").Parse(jobSource)
	if err != nil {
		return 500, err
	}

	var jobSpec bytes.Buffer

	if err := jobTemplate.Execute(&jobSpec, job); err != nil {
		return 500, err
	}

	parseBodySource := `{ "JobHCL": "{{ . }}", "Canonicalize": true }`

	parseBodyTemplate, err := template.New("parseBody").Parse(parseBodySource)
	if err != nil {
		return 500, err
	}

	var parseBody bytes.Buffer

	if err := parseBodyTemplate.Execute(&parseBody, jobSpec.String()); err != nil {
		return 500, err
	}

	parseBodyCleaned := strings.ReplaceAll(parseBody.String(), "\n", `\n`)
	parseBodyCleaned = strings.ReplaceAll(parseBodyCleaned, "\t", ` `)

	println(parseBodyCleaned)

	parseResponse, err := http.Post("http://zeus.internal:4646/v1/jobs/parse", "application/json", strings.NewReader(parseBodyCleaned))
	if err != nil {
		return 500, err
	}

	// parse response
	var parseResponseBuffer bytes.Buffer

	if _, err := io.Copy(&parseResponseBuffer, parseResponse.Body); err != nil {
		return 500, err
	}

	jobRunSource := `{ "Job": {{ . }} }`

	jobRunTemplate, err := template.New("jobRun").Parse(jobRunSource)
	if err != nil {
		return 500, err
	}

	var jobRun bytes.Buffer

	if err := jobRunTemplate.Execute(&jobRun, parseResponseBuffer.String()); err != nil {
		return 500, err
	}

	jobRunResponse, err := http.Post("http://zeus.internal:4646/v1/jobs", "application/json", &jobRun)
	if err != nil {
		return 500, err
	}

	// parse response
	var jobRunResponseBuffer bytes.Buffer

	if _, err := io.Copy(&jobRunResponseBuffer, jobRunResponse.Body); err != nil {
		return 500, err
	}

	println(jobRunResponseBuffer.String())

	return 200, nil
}
