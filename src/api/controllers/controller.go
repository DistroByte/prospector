package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const NOMAD_URL = "http://zeus.internal:4646/v1"

type Controller struct {
	Client        NomadClient
	IdentityKey   string
	JWTMiddleware *jwt.GinJWTMiddleware
}

type Message struct {
	Message string `json:"message"`
}

type Resources struct {
	Cpu    int `json:"cpu" validate:"min=0,max=2000"`
	Memory int `json:"memory" validate:"min=0,max=2000"`
}

type Network struct {
	Port   int    `json:"port" validate:"min=0,max=65535"`
	Expose bool   `json:"expose" validate:"optional" default:"false"`
	Mac    string `json:"mac" validate:"optional"`
}

type UserConfig struct {
	User   string `json:"user" validate:"optional"`
	SSHKey string `json:"ssh_key" validate:"optional"`
}

type Job struct {
	Name       string `json:"name" validate:"required"`
	Type       string `json:"type" validate:"required"`
	Image      string `json:"image"`
	Resources  `json:"resources"`
	Network    `json:"network"`
	UserConfig `json:"user_config"`
}

type ShortJob struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type NomadClient interface {
	Get(endpoint string) ([]byte, error)
	Post(endpoint string, reqBody *bytes.Buffer) ([]byte, error)
	Delete(endpoint string) ([]byte, error)
}

type DefaultNomadClient struct{}

func (n *DefaultNomadClient) Get(endpoint string) ([]byte, error) {
	url := NOMAD_URL + endpoint

	resp, err := http.Get(url)

	if err != nil {
		return nil, gin.Error{Err: err, Meta: resp.StatusCode}
	}

	if resp.StatusCode != 200 {
		return nil, gin.Error{Err: err, Meta: resp.StatusCode}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, gin.Error{Err: err, Meta: resp.StatusCode}
	}

	return body, nil
}

func (n *DefaultNomadClient) Post(endpoint string, reqBody *bytes.Buffer) ([]byte, error) {
	url := NOMAD_URL + endpoint

	resp, err := http.Post(url, "application/json", reqBody)
	if err != nil {
		return nil, gin.Error{Err: err, Meta: resp.StatusCode}
	}

	if resp.StatusCode != 200 {
		return nil, gin.Error{Err: err, Meta: resp.StatusCode}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, gin.Error{Err: err, Meta: resp.StatusCode}
	}

	return body, nil
}

func (n *DefaultNomadClient) Delete(endpoint string) ([]byte, error) {
	url := NOMAD_URL + endpoint

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, gin.Error{Err: err, Meta: http.StatusInternalServerError}
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, gin.Error{Err: err, Meta: resp.StatusCode}
	}

	if resp.StatusCode != 200 {
		return nil, gin.Error{Err: err, Meta: resp.StatusCode}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, gin.Error{Err: err, Meta: resp.StatusCode}
	}

	return body, nil
}

func (j *Job) ToJson() *bytes.Buffer {
	jobJson, _ := json.Marshal(j)
	return bytes.NewBuffer(jobJson)
}
