package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DefaultNomadClient struct {
	URL string
}

func (n *DefaultNomadClient) Get(endpoint string) ([]byte, error) {
	url := n.URL + endpoint

	resp, err := http.Get(url)

	if err != nil {
		return nil, gin.Error{Err: errors.New("nomad error: " + err.Error()), Meta: resp.StatusCode}
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
	url := n.URL + endpoint

	resp, err := http.Post(url, "application/json", reqBody)
	if err != nil {
		return nil, gin.Error{Err: err, Meta: resp.StatusCode}
	}

	// if response is not 2xx
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, gin.Error{Err: errors.New("bad response received with error: " + fmt.Sprint(resp.StatusCode)), Meta: resp.StatusCode}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, gin.Error{Err: err, Meta: resp.StatusCode}
	}

	return body, nil
}

func (n *DefaultNomadClient) Delete(endpoint string) ([]byte, error) {
	url := n.URL + endpoint

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, gin.Error{Err: err, Meta: http.StatusInternalServerError}
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, gin.Error{Err: errors.New("nomad error: " + err.Error()), Meta: resp.StatusCode}
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
