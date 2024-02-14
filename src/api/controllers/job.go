package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	nomad "github.com/hashicorp/nomad/nomad/structs"
)

// GetJobs gets all the jobs from nomad that have the word "prospector" in their name
//
//	@Summary		Get all jobs
//	@Description	Get all jobs from nomad
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		None
//	@Success		200	{object}	[]nomad.JobListStub
//	@Router			/jobs [get]
//	@Param			long	query	string	false	"Get long job details"
func (c *Controller) GetJobs(ctx *gin.Context) {
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
		if strings.Contains(job.Name, "-prospector") {
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
//	@Security		None
//	@Success		200	{object}	nomad.Job
//	@Router			/jobs/{id} [get]
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

// CreateJob creates a job in nomad
//
//	@Summary		Create a job
//	@Description	Create and submit a job to nomad to deploy
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		None
//	@Param			job	body		Job	true	"Job"
//	@Success		200	{object}	Message
//	@Router			/jobs [post]
func (c *Controller) CreateJob(ctx *gin.Context) {
	var job Job

	if err := ctx.BindJSON(&job); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := createNomadJob(job)

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
//	@Security		None
//	@Success		200	{object}	Message
//	@Router			/jobs/{id} [delete]
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

func createNomadJob(job Job) (int, error) {
	jobSource := `job "{{ .Name }}-prospector" {
	datacenters = ["dc1"]
	type = "service"

	group "{{ .Name }}-prospector" {
		count = 1

		network {
			port "web" {
				to = {{ .Port }}
			}
		}

		service {
			name = "{{ .Name }}"
			port = "web"

			check {
				name = "{{ .Name }}-health"
				type = "http"
				path = "/"
				interval = "10s"
				timeout = "2s"
			}

			tags = [
				"traefik.enable=true",
				"traefik.http.routers.{{ .Name }}-prospector.rule=Host(` + "`" + `{{ .Name }}.prospector.ie` + "`" + `)",
				"traefik.http.routers.{{ .Name }}-prospector.entrypoints=websecure",
				"traefik.http.routers.{{ .Name }}-prospector.tls=true",
				"traefik.http.routers.{{ .Name }}-prospector.tls.certresolver=lets-encrypt"
			]
		}

		task "{{ .Name }}-prospector" {
			driver = "docker"
			
			config {
				image = "{{ .Image }}"
				ports = ["web"]
			}

			resources {
				cpu    = {{ .Cpu }}
				memory = {{ .Memory }}
			}
		}
	}
}
`

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
