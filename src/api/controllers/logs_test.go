package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	nomad "github.com/hashicorp/nomad/nomad/structs"
)

type MockNomadClientLogs struct{}

func (m *MockNomadClientLogs) Get(path string) ([]byte, error) {
	switch path {
	case "/job/valid-endpoint-prospector/allocations":
		allocs := []nomad.AllocListStub{
			{
				ID:           "valid-endpoint",
				ClientStatus: "running",
			},
		}

		allocsBytes, err := json.Marshal(allocs)
		if err != nil {
			return nil, err
		}

		return allocsBytes, nil
	case "/client/fs/logs/valid-endpoint":
		return []byte(`{"ID": "valid-endpoint"}`), nil
	}

	return nil, fmt.Errorf("error")
}

func (m *MockNomadClientLogs) Delete(path string) ([]byte, error) {
	switch path {
	case "/job/valid-endpoint-prospector/allocations":
		allocs := []nomad.AllocListStub{
			{
				ID: "valid-endpoint",
			},
		}

		allocsBytes, err := json.Marshal(allocs)
		if err != nil {
			return nil, err
		}

		return allocsBytes, nil
	case "/client/fs/logs/valid-endpoint":
		return []byte(`{"ID": "valid-endpoint"}`), nil
	}

	return nil, fmt.Errorf("error")
}

func (m *MockNomadClientLogs) Forward(ctx *gin.Context, path string) (*http.Response, error) {
	println(path)
	switch path {
	case "/client/fs/logs/valid-endpoint?follow=true&offset=5000&origin=end&plain=true&task=&type=stdout":

		// make a mock response for a reader
		// make this a read closer
		// return a response with a body

		return &http.Response{
			Body:       io.NopCloser(bytes.NewBufferString("test")),
			StatusCode: 200,
		}, nil
	}

	return nil, fmt.Errorf("error")
}

func (m *MockNomadClientLogs) Post(path string, reqBody *bytes.Buffer) ([]byte, error) {
	return nil, nil
}

func TestStreamLogs(t *testing.T) {
	c := NomadProxyController{
		Client: &MockNomadClientLogs{},
	}

	tcs := []struct {
		name   string
		path   string
		expect error
	}{
		{
			name:   "valid job",
			path:   "valid-endpoint-prospector",
			expect: nil,
		},
		{
			name:   "invalid job",
			path:   "invalid-endpoint-prospector",
			expect: nil,
		},
		{
			name:   "invalid response",
			path:   "invalid-response-prospector",
			expect: nil,
		},
		{
			name:   "no job",
			path:   "no-job-prospector",
			expect: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(w)

			// set query param ?type=stdout
			ctx.Request = httptest.NewRequest("GET", "/logs?type=stdout&task=", nil)
			ctx.Request.URL.RawQuery = "type=stdout"

			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: tc.path})

			c.StreamLogs(ctx)

			println(w.Body.String())
		})
	}

}
