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

	if !helpers.CheckJobHasValidName(ctx, jobId) {
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

	if !helpers.CheckJobHasValidName(ctx, id) {
		return
	}

	data, err := c.Client.Get("/job/" + id + "/allocations")
	if err != nil {
		ctx.Error(err)
	}

	var allocs []nomad.AllocListStub
	err = json.Unmarshal(data, &allocs)
	if err != nil {
		ctx.Error(err)
	}

	if len(allocs) == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "No components found"})
		return
	}

	var taskGroups []ComponentStatus
	for _, alloc := range allocs {
		taskGroups = append(taskGroups, ComponentStatus{
			Name:  alloc.Name[strings.LastIndex(alloc.Name, ".")+1 : len(alloc.Name)-3],
			State: alloc.ClientStatus,
		})
	}

	ctx.JSON(http.StatusOK, taskGroups)
}
