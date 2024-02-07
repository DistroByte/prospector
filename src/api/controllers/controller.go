package controllers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const NOMAD_URL = "http://zeus.internal:4646/v1"

type Controller struct {
	Client NomadClient
}

type Message struct {
	Message string `json:"message"`
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
