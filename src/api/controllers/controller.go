package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	nomad "github.com/hashicorp/nomad/nomad/structs"
)

type Controller struct {
	Client        NomadClient
	IdentityKey   string
	JWTMiddleware *jwt.GinJWTMiddleware
}

type Message struct {
	Message string `json:"message"`
}

type Project struct {
	Name       string      `json:"name" validate:"required"`
	Type       string      `json:"type" validate:"required"`
	Components []Component `json:"components" validate:"required"`
	User       string      `json:"-"`
}

type Component struct {
	Name       string   `json:"name" validate:"required"`
	Image      string   `json:"image"`
	Volumes    []string `json:"volumes"`
	Resources  `json:"resources"`
	Network    `json:"network"`
	UserConfig `json:"user_config"`
}

type Network struct {
	Port   int    `json:"port" validate:"min=0,max=65535"`
	Expose bool   `json:"expose" validate:"optional" default:"false"`
	Mac    string `json:"-"`
}
type Resources struct {
	Cpu    int `json:"cpu" validate:"min=0,max=2000"`
	Memory int `json:"memory" validate:"min=0,max=2000"`
}

type UserConfig struct {
	User   string `json:"user" validate:"optional"`
	SSHKey string `json:"ssh_key" validate:"optional"`
}

type ShortJob struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Type    string `json:"type"`
	Created int    `json:"created"`
}

type ComponentStatus struct {
	Name         string `json:"name"`
	State        string `json:"state"`
	DateModified int    `json:"date_modified"`
	Image        string `json:"image"`
}

type Utilization struct {
	Cpu       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Timestamp int     `json:"timestamp"`
}

type NomadClient interface {
	Get(endpoint string) ([]byte, error)
	Post(endpoint string, reqBody *bytes.Buffer) ([]byte, error)
	Delete(endpoint string) ([]byte, error)
	Forward(ctx *gin.Context, endpoint string) (*http.Response, error)
}

func (j *Project) ToJson() *bytes.Buffer {
	jobJson, _ := json.Marshal(j)
	return bytes.NewBuffer(jobJson)
}

func (c *Controller) parseRunningAllocs(jobId string) ([]nomad.AllocListStub, error) {
	data, err := c.Client.Get("/job/" + jobId + "/allocations")
	if err != nil {
		return nil, err
	}

	var allocations []nomad.AllocListStub
	err = json.Unmarshal(data, &allocations)
	if err != nil {
		return nil, err
	}

	var runningAllocs []nomad.AllocListStub
	for _, alloc := range allocations {
		if alloc.ClientStatus == "running" || alloc.ClientStatus == "pending" {
			runningAllocs = append(runningAllocs, alloc)
		}
	}

	if len(runningAllocs) > 0 {
		return runningAllocs, nil
	}

	return nil, gin.Error{Err: fmt.Errorf("no running allocations found"), Meta: 404}
}

func (c *Controller) getJobFromNomad(id string) (nomad.Job, error) {
	data, err := c.Client.Get("/job/" + id)
	if err != nil {
		return nomad.Job{}, err
	}

	var job nomad.Job
	err = json.Unmarshal(data, &job)
	if err != nil {
		return nomad.Job{}, err
	}

	return job, nil
}
