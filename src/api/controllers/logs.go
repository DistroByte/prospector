package controllers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"prospector/helpers"

	"github.com/gin-gonic/gin"
	nomad "github.com/hashicorp/nomad/nomad/structs"
)

type NomadProxyController struct {
	Client NomadClient
}

// StreamLogs godoc
//
//	@Summary		StreamLogs
//	@Description	Stream logs from a Nomad job allocation
//	@Tags			proxy
//	@Produce		json
//	@Param			id		path	string	true	"Project ID"
//	@Param			task	query	string	true	"Component name"
//	@Param			type	query	string	true	"Log type (stdout or stderr)"
//	@Success		200
//	@Security		BearerAuth
//	@Router			/v1/jobs/{id}/logs [get]
func (c *NomadProxyController) StreamLogs(ctx *gin.Context) {
	jobId := ctx.Param("id")
	logType := ctx.Query("type")
	if logType != "stdout" && logType != "stderr" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log type"})
		return
	}

	if !helpers.CheckJobHasValidName(jobId) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid job ID"})
		return
	}

	allocs, err := c.parseRunningAlloc(jobId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if allocs == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No running allocations found"})
		return
	}

	path := "/client/fs/logs/" + allocs.ID
	queryParams := url.Values{}
	queryParams.Add("task", ctx.Query("task"))
	queryParams.Add("type", logType)
	queryParams.Add("follow", "true")
	queryParams.Add("offset", "5000")
	queryParams.Add("plain", "true")
	queryParams.Add("origin", "end")

	url := path + "?" + queryParams.Encode()
	resp, err := c.Client.Forward(ctx, url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	fmt.Println("Started streaming logs")

	err = streamResponse(ctx, resp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Finished streaming logs")
}

func (c *NomadProxyController) parseRunningAlloc(jobId string) (*nomad.AllocListStub, error) {
	data, err := c.Client.Get(fmt.Sprintf("/job/%s/allocations", jobId))
	if err != nil {
		return nil, err
	}

	var allocs []nomad.AllocListStub
	err = json.Unmarshal(data, &allocs)
	if err != nil {
		return nil, err
	}
	for _, alloc := range allocs {
		if alloc.ClientStatus == "running" || alloc.ClientStatus == "pending" {
			return &alloc, nil
		}
	}

	return nil, err
}

func streamResponse(ctx *gin.Context, resp *http.Response) error {
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		defer gzipReader.Close()
		resp.Body = gzipReader
	}

	buf := make([]byte, 1024)
	for {
		bytesRead, err := resp.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		select {
		case <-ctx.Request.Context().Done():
			return nil
		default:
		}

		_, err = ctx.Writer.Write(buf[:bytesRead])
		if err != nil {
			return err
		}

		if f, ok := ctx.Writer.(http.Flusher); ok {
			f.Flush()
		}
	}

	return nil
}
