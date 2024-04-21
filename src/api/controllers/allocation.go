package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"prospector/helpers"
	"strings"

	"github.com/gin-gonic/gin"
	nomad "github.com/hashicorp/nomad/nomad/structs"
)

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

	if !helpers.CheckJobHasValidName(jobId) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid job ID"})
		return
	}

	if taskName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task name"})
		return
	}

	alloc, err := c.parseRunningAllocs(jobId)
	if err != nil {
		ctx.Error(err)
		return
	}

	for _, alloc := range alloc {
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
//	@Success		200	{object}	[]ComponentStatus
func (c *Controller) GetComponents(ctx *gin.Context) {
	id := ctx.Param("id")

	if !helpers.CheckJobHasValidName(id) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid job ID"})
		return
	}

	allocData, err := c.Client.Get("/job/" + id + "/allocations")
	if err != nil {
		ctx.Error(err)
	}

	jobData, err := c.Client.Get("/job/" + id)
	if err != nil {
		ctx.Error(err)
	}

	var allocs []nomad.AllocListStub
	err = json.Unmarshal(allocData, &allocs)
	if err != nil {
		ctx.Error(err)
	}

	var job nomad.Job
	err = json.Unmarshal(jobData, &job)
	if err != nil {
		ctx.Error(err)
	}

	if len(allocs) == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "No components found"})
		return
	}

	var components []ComponentStatus

	for _, taskGroup := range job.TaskGroups {
		for _, task := range taskGroup.Tasks {
			var component ComponentStatus
			component.Name = task.Name
			component.State = job.Status
			component.DateModified = 0
			if task.Config["image"] != nil {
				component.Image = task.Config["image"].(string)
			} else {
				component.Image = strings.Split(strings.Split(task.Artifacts[0].GetterSource, "/")[len(strings.Split(task.Artifacts[0].GetterSource, "/"))-1], ".")[0]
			}

			for _, taskState := range allocs {
				if component.DateModified < int(taskState.ModifyTime) && strings.Contains(taskState.Name, task.Name) {
					component.DateModified = int(taskState.ModifyTime)
				}
			}

			components = append(components, component)
		}
	}

	ctx.JSON(http.StatusOK, components)
}
