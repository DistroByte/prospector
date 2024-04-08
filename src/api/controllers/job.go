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
	
	meta {
		job-type = "docker"
	}
	
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

  meta {
	job-type = "vm"
  }

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
//	@Summary		Get all projects
//	@Description	Get all jobs from nomad
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/jobs [get]
//	@Param			long	query	boolean	false	"Get long project details"
//	@Param			running	query	boolean	false	"Get running projects"
//	@Code			204 "No projects found"
//	@Success		200	{object}	[]ShortJob
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
			Type:   job.Meta["job-type"],
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
//	@Summary		Get a project
//	@Description	Get a job from nomad
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/jobs/{id} [get]
//	@Param			id	path	string	true	"Project ID"
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

// CreateJob creates a container or VM
//
//	@Summary		Create a project
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

// DeleteJob deletes a project
//
//	@Summary		Delete a project
//	@Description	Delete a job from nomad
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	Message
//	@Router			/v1/jobs/{id} [delete]
//	@Param			id		path	string	true	"Project ID"
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

// RestartJob restarts a project
//
//	@Summary		Restart a project
//	@Description	Restart a job in nomad
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	Message
//	@Router			/v1/jobs/{id}/restart [put]
//	@Param			id	path	string	true	"Project ID"
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

	ctx.JSON(http.StatusOK, gin.H{"message": "Project restarted successfully"})
}

// RestartAlloc restarts a component in a project
//
//	@Summary		Restart a component
//	@Description	Restart a component in a project
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	Message
//	@Router			/v1/jobs/{id}/component/{component}/restart [put]
//	@Param			id			path	string	true	"Project ID"
//	@Param			component	path	string	true	"Component name"
func (c *Controller) RestartAlloc(ctx *gin.Context) {
	taskName := ctx.Param("component")
	jobId := ctx.Param("id")

	if taskName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task name"})
		return
	}

	alloc, err := c.parseRunningAllocs(jobId)
	if err != nil {
		ctx.Error(err)
		return
	}

	// send { "TaskName": "taskName" }
	body := bytes.NewBuffer([]byte(`{ "TaskName": "` + taskName + `" }`))
	data, err := c.Client.Post("/client/allocation/"+alloc.ID+"/restart", body)
	if err != nil {
		ctx.Error(err)
	}

	var response nomad.GenericResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Componet restarted successfully"})
}

// GetComponents gets all the components in a project
//
//	@Summary		Get all components
//	@Description	Get all components in a project
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/jobs/{id}/components [get]
//	@Param			id	path	string	true	"Project ID"
//	@Code			204 "No components found"
//	@Success		200	{object}	[]string
func (c *Controller) GetComponents(ctx *gin.Context) {
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

	if job.TaskGroups == nil {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "No components found"})
		return
	}

	var taskGroups []string
	for _, task := range job.TaskGroups {
		taskGroups = append(taskGroups, task.Name)
	}

	ctx.JSON(http.StatusOK, taskGroups)
}

// GetAllocatedResources gets the total CPU and memory allocated to a user's running jobs
//
//	@Summary		Get allocated resources
//	@Description	Get the total CPU and memory allocated to a user's running jobs
//	@Tags			resources
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/resources/allocated [get]
//	@Success		200	{object}	Resources
func (c *Controller) GetAllocatedResources(ctx *gin.Context) {
	claims := jwt.ExtractClaims(ctx)

	// get all user jobs
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

	// parse running allocations from jobs
	var runningJobs []ShortJob = []ShortJob{}
	for _, job := range filteredJobs {
		if job.Status == "running" {
			runningJobs = append(runningJobs, ShortJob{
				ID:     job.ID,
				Status: job.Status,
			})
		}
	}

	if len(runningJobs) == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "No running jobs found"})
		return
	}

	var resources []Resources

	// get task groups
	for _, job := range runningJobs {
		data, err := c.Client.Get("/job/" + job.ID)
		if err != nil {
			ctx.Error(err)
		}

		var job nomad.Job
		err = json.Unmarshal(data, &job)
		if err != nil {
			ctx.Error(err)
		}

		if job.TaskGroups == nil {
			ctx.JSON(http.StatusNoContent, gin.H{"message": "No components found"})
			return
		}

		for _, group := range job.TaskGroups {
			for _, task := range group.Tasks {
				resources = append(resources, Resources{
					Cpu:    task.Resources.CPU,
					Memory: task.Resources.MemoryMB,
				})
			}
		}
	}

	var totalCpu int
	var totalMemory int

	for _, resource := range resources {
		totalCpu += resource.Cpu
		totalMemory += resource.Memory
	}

	ctx.JSON(http.StatusOK, gin.H{"cpu": totalCpu, "memory": totalMemory})
}

// GetAllUsedResources gets the current CPU and memory usage of a user's running jobs
//
//	@Summary		Get all used resources
//	@Description	Get the current CPU and memory usage of all a user's running jobs
//	@Tags			resources
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/resources [get]
//	@Success		200	{object}	Utilization
//	@Code			204 "No running jobs found"
func (c *Controller) GetAllUsedResources(ctx *gin.Context) {
	claims := jwt.ExtractClaims(ctx)

	// get all user jobs
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

	// parse running allocations from jobs
	var runningJobs []ShortJob = []ShortJob{}
	for _, job := range filteredJobs {
		if job.Status == "running" {
			runningJobs = append(runningJobs, ShortJob{
				ID:     job.ID,
				Status: job.Status,
			})
		}
	}

	if len(runningJobs) == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "No running jobs found"})
		return
	}

	var util Utilization

	// get all running allocations
	for _, job := range runningJobs {
		data, err := c.Client.Get("/job/" + job.ID + "/allocations")
		if err != nil {
			ctx.Error(err)
			return
		}

		var allocations []nomad.AllocListStub
		err = json.Unmarshal(data, &allocations)
		if err != nil {
			ctx.Error(err)
			return
		}

		// parse json into map
		var stats map[string]interface{}
		for _, alloc := range allocations {
			if alloc.ClientStatus == "running" {
				data, err = c.Client.Get("/client/allocation/" + alloc.ID + "/stats")
				if err != nil {
					ctx.Error(err)
					return
				}

				err = json.Unmarshal(data, &stats)
				if err != nil {
					ctx.Error(err)
					return
				}

				cpuStats := stats["ResourceUsage"].(map[string]interface{})["CpuStats"].(map[string]interface{})
				memoryStats := stats["ResourceUsage"].(map[string]interface{})["MemoryStats"].(map[string]interface{})

				util = Utilization{
					Cpu:       cpuStats["Percent"].(float64) + util.Cpu,
					Memory:    memoryStats["Usage"].(float64)/1024/1024 + util.Memory,
					Timestamp: int(stats["Timestamp"].(float64)),
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, util)
}

// GetJobUsedResources gets the current CPU and memory usage of a job
//
//	@Summary		Get job resource usage
//	@Description	Get the current CPU and memory usage of a user's job
//	@Tags			resources
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/resources/{id} [get]
//	@Param			id	path		string	true	"Project ID"
//	@Success		200	{object}	Utilization
//	@Code			204 "No running jobs found"
func (c *Controller) GetJobUsedResources(ctx *gin.Context) {
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

	var util Utilization

	// get all running allocations
	data, err = c.Client.Get("/job/" + job.ID + "/allocations")
	if err != nil {
		ctx.Error(err)
		return
	}

	var allocations []nomad.AllocListStub
	err = json.Unmarshal(data, &allocations)
	if err != nil {
		ctx.Error(err)
		return
	}

	// parse json into map
	var stats map[string]interface{}
	for _, alloc := range allocations {
		if alloc.ClientStatus == "running" {
			data, err = c.Client.Get("/client/allocation/" + alloc.ID + "/stats")
			if err != nil {
				ctx.Error(err)
				return
			}

			err = json.Unmarshal(data, &stats)
			if err != nil {
				ctx.Error(err)
				return
			}

			cpuStats := stats["ResourceUsage"].(map[string]interface{})["CpuStats"].(map[string]interface{})
			memoryStats := stats["ResourceUsage"].(map[string]interface{})["MemoryStats"].(map[string]interface{})

			util = Utilization{
				Cpu:       cpuStats["Percent"].(float64) + util.Cpu,
				Memory:    memoryStats["Usage"].(float64)/1024/1024 + util.Memory,
				Timestamp: int(stats["Timestamp"].(float64)),
			}
		}
	}

	ctx.JSON(http.StatusOK, util)
}

// GetComponentUsedResources gets the current CPU and memory usage of a component
//
//	@Summary		Get component resource usage
//	@Description	Get the current CPU and memory usage of a job's component
//	@Tags			resources
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/resources/{id}/{component} [get]
//	@Param			id			path		string	true	"Project ID"
//	@Param			component	path		string	true	"Component name"
//	@Success		200			{object}	Utilization
//	@Code			204 "No running jobs found"
func (c *Controller) GetComponentUsedResources(ctx *gin.Context) {
	id := ctx.Param("id")
	component := ctx.Param("component")

	data, err := c.Client.Get("/job/" + id)
	if err != nil {
		ctx.Error(err)
	}

	var job nomad.Job
	err = json.Unmarshal(data, &job)
	if err != nil {
		ctx.Error(err)
	}

	// get all running allocations
	data, err = c.Client.Get("/job/" + job.ID + "/allocations")
	if err != nil {
		ctx.Error(err)
		return
	}

	var allocations []nomad.AllocListStub
	err = json.Unmarshal(data, &allocations)
	if err != nil {
		ctx.Error(err)
		return
	}

	var util Utilization

	// parse json into map
	var stats map[string]interface{}
	for _, alloc := range allocations {
		if alloc.TaskGroup == component {
			data, err = c.Client.Get("/client/allocation/" + alloc.ID + "/stats")
			if err != nil {
				ctx.Error(err)
				return
			}

			err = json.Unmarshal(data, &stats)
			if err != nil {
				ctx.Error(err)
				return
			}

			cpuStats := stats["ResourceUsage"].(map[string]interface{})["CpuStats"].(map[string]interface{})
			memoryStats := stats["ResourceUsage"].(map[string]interface{})["MemoryStats"].(map[string]interface{})

			util = Utilization{
				Cpu:       cpuStats["Percent"].(float64) + util.Cpu,
				Memory:    memoryStats["Usage"].(float64)/1024/1024 + util.Memory,
				Timestamp: int(stats["Timestamp"].(float64)) + util.Timestamp,
			}
		}
	}

	ctx.JSON(http.StatusOK, util)
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
