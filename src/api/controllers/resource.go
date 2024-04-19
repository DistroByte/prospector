package controllers

import (
	"encoding/json"
	"net/http"
	"prospector/helpers"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	nomad "github.com/hashicorp/nomad/nomad/structs"
)

// GetUserAllocatedResources gets the total CPU and memory allocated to a user's running projects
//
//	@Summary		Get allocated resources
//	@Description	Get the total CPU and memory allocated to a user's running projects
//	@Tags			resources
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/resources/allocated [get]
//	@Success		200	{object}	Resources
func (c *Controller) GetUserAllocatedResources(ctx *gin.Context) {
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
			runningJobs = append(runningJobs, ShortJob{ID: job.ID, Status: job.Status})
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
				resources = append(resources, Resources{Cpu: task.Resources.CPU, Memory: task.Resources.MemoryMB})
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

// GetUserUsedResources gets the current CPU and memory usage of a user's running projects
//
//	@Summary		Get all used resources
//	@Description	Get the current CPU and memory usage of all a user's running projects
//	@Tags			resources
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/resources [get]
//	@Success		200	{object}	Utilization
//	@Code			204 "No running jobs found"
func (c *Controller) GetUserUsedResources(ctx *gin.Context) {
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
			runningJobs = append(runningJobs, ShortJob{ID: job.ID, Status: job.Status})
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

// GetJobAllocatedResources gets the total CPU and memory allocated to a project
//
//	@Summary		Get project allocated resources
//	@Description	Get the total CPU and memory allocated to a user's project
//	@Tags			resources
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/resources/{id}/allocated [get]
//	@Param			id	path		string	true	"Project ID"
//	@Success		200	{object}	Resources
//	@Code			204 "No running projects found"
func (c *Controller) GetJobAllocatedResources(ctx *gin.Context) {
	id := ctx.Param("id")

	if !helpers.CheckJobHasValidName(id) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid job ID"})
		return
	}

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

	var resources []Resources

	for _, group := range job.TaskGroups {
		for _, task := range group.Tasks {
			resources = append(resources, Resources{Cpu: task.Resources.CPU, Memory: task.Resources.MemoryMB})
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

// GetJobUsedResources gets the current CPU and memory usage of a project
//
//	@Summary		Get project resource usage
//	@Description	Get the current CPU and memory usage of a user's project
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

	if !helpers.CheckJobHasValidName(id) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid job ID"})
		return
	}

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

// GetComponentAllocatedResources gets the total CPU and memory allocated to a component
//
//	@Summary		Get component allocated resources
//	@Description	Get the total CPU and memory allocated to a project's component
//	@Tags			resources
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/resources/{id}/{component}/allocated [get]
//	@Param			id			path		string	true	"Project ID"
//	@Param			component	path		string	true	"Component name"
//	@Success		200			{object}	Resources
//	@Code			204 "No running projects found"
func (c *Controller) GetComponentAllocatedResources(ctx *gin.Context) {
	id := ctx.Param("id")
	component := ctx.Param("component")

	if !helpers.CheckJobHasValidName(id) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid job ID"})
		return
	}

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

	var resources Resources

	for _, group := range job.TaskGroups {
		for _, task := range group.Tasks {
			if task.Name == component {
				resources = Resources{Cpu: task.Resources.CPU, Memory: task.Resources.MemoryMB}
			}
		}
	}

	ctx.JSON(http.StatusOK, resources)
}

// GetComponentUsedResources gets the current CPU and memory usage of a component
//
//	@Summary		Get component resource usage
//	@Description	Get the current CPU and memory usage of a project's component
//	@Tags			resources
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/v1/resources/{id}/{component} [get]
//	@Param			id			path		string	true	"Project ID"
//	@Param			component	path		string	true	"Component name"
//	@Success		200			{object}	Utilization
//	@Code			204 "No running projects found"
func (c *Controller) GetComponentUsedResources(ctx *gin.Context) {
	id := ctx.Param("id")
	component := ctx.Param("component")

	if !helpers.CheckJobHasValidName(id) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid job ID"})
		return
	}

	var job nomad.Job
	job, err := c.getJobFromNomad(id)
	if err != nil {
		ctx.Error(err)
	}
	var allocations []nomad.AllocListStub
	allocations, err = c.parseRunningAllocs(job.ID)
	if err != nil {
		ctx.Error(err)
		return
	}

	var util Utilization

	// parse json into map
	var stats map[string]interface{}
	for _, alloc := range allocations {
		if alloc.TaskGroup == component {
			data, err := c.Client.Get("/client/allocation/" + alloc.ID + "/stats")
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
