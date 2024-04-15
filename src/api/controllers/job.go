package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"prospector/helpers"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	nomad "github.com/hashicorp/nomad/nomad/structs"
)

// GetJobs gets all the jobs from nomad that have the word "prospector" in their name
//
//	@Summary		Get all projects
//	@Description	Get all projects from nomad
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
			ID:      job.ID,
			Status:  job.Status,
			Type:    job.Meta["job-type"],
			Created: int(job.SubmitTime),
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
//	@Param			job	body		Project	true	"Project"
//	@Success		200	{object}	Message
//	@Router			/v1/jobs [post]
func (c *Controller) CreateJob(ctx *gin.Context) {
	var job Project
	if err := ctx.BindJSON(&job); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := jwt.ExtractClaims(ctx)
	job.User = claims[c.IdentityKey].(string)

	for i := 0; i < len(job.Components); i++ {
		job.Components[i].UserConfig.User = job.User
		// generate random mac address
		if job.Type == "vm" {
			job.Components[i].Network.Mac = fmt.Sprintf("52:54:00:%02x:%02x:%02x", byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)))
		}
	}

	if job.Type == "docker" {
		_, err := CreateJobFromTemplate(job, DockerSource)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if job.Type == "vm" {
		err := WriteTextFilesForVM(job)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_, err = CreateJobFromTemplate(job, VMSource)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job type"})
		return
	}

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
//	@Param			purge	query	bool	false	"Purge project"
func (c *Controller) DeleteJob(ctx *gin.Context) {
	id := ctx.Param("id")

	purge := ctx.Query("purge")

	if !helpers.CheckJobHasValidName(id) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid job ID"})
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

	if !helpers.CheckJobHasValidName(id) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid job ID"})
		return
	}

	alloc, err := c.parseRunningAllocs(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	for _, alloc := range alloc {
		body := bytes.NewBuffer([]byte{})

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
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Project restarted successfully"})
}

// StartJob starts a project after it has been stopped
//
//	@Summary		Start a project
//	@Description	Start a job in nomad
//	@Tags			job
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	Message
//	@Router			/v1/jobs/{id}/start [post]
//	@Param			id	path	string	true	"Project ID"
func (c *Controller) StartJob(ctx *gin.Context) {
	id := ctx.Param("id")

	if !helpers.CheckJobHasValidName(id) {
		return
	}

	body := bytes.NewBuffer([]byte(`{ "JobID": "` + id + `" }`))

	data, err := c.Client.Post("/job/"+id+"/revert", body)
	if err != nil {
		ctx.Error(err)
		return
	}

	var response nomad.GenericResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		ctx.Error(err)
	}

	ctx.JSON(http.StatusOK, response)
}
