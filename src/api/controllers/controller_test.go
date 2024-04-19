package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	nomad "github.com/hashicorp/nomad/nomad/structs"
)

func TestToJson(t *testing.T) {
	project := Project{
		Name: "Test Project",
		Type: "Test Type",
		Components: []Component{
			{
				Name:  "Test Component",
				Image: "Test Image",
				Resources: Resources{
					Cpu:    1,
					Memory: 1,
				},
				Network: Network{
					Port:   1,
					Expose: false,
				},
			},
		},
	}

	json := project.ToJson()
	if json == nil {
		t.Errorf("expected json to not be nil")
	}

	if json.String() != `{"name":"Test Project","type":"Test Type","components":[{"name":"Test Component","image":"Test Image","resources":{"cpu":1,"memory":1},"network":{"port":1,"expose":false},"user_config":{"user":"","ssh_key":""}}]}` {
		t.Errorf("expected json to be %s, got %s", `{"name":"Test Project","type":"Test Type","components":[{"name":"Test Component","image":"Test Image","resources":{"cpu":1,"memory":1},"network":{"port":1,"expose":false},"user_config":{"user":"","ssh_key":""}}]}`, json.String())
	}
}

func TestParseRunningAllocs(t *testing.T) {
	c := Controller{
		Client: &MockNomadClient{},
	}

	tcs := []struct {
		name string
		path string
	}{
		{
			name: "valid job",
			path: "test-valid-endpoint",
		},
		{
			name: "invalid job",
			path: "test-invalid-endpoint",
		},
		{
			name: "invalid response",
			path: "test-invalid-response",
		},
		{
			name: "no allocations",
			path: "test-no-allocations",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			_, err := c.parseRunningAllocs(tc.path)

			if err != nil {
				if tc.path != "test-valid-endpoint" {
					return
				} else {
					t.Errorf("expected error to be nil, got %v", err)
					return
				}
			}
		})
	}
}

func TestGetJobFromNomad(t *testing.T) {
	c := Controller{
		Client: &MockNomadClient{},
	}

	tcs := []struct {
		name string
		path string
	}{
		{
			name: "valid job",
			path: "test-valid-endpoint",
		},
		{
			name: "invalid job",
			path: "test-invalid-endpoint",
		},
		{
			name: "invalid response",
			path: "test-invalid-response",
		},
		{
			name: "no jobs",
			path: "test-no-jobs",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			_, err := c.getJobFromNomad(tc.path)

			if err != nil {
				if tc.path != "test-valid-endpoint" {
					return
				} else {
					t.Errorf("expected error to be nil, got %v", err)
					return
				}
			}
		})
	}
}

type MockNomadClient struct{}

func (m *MockNomadClient) Get(endpoint string) ([]byte, error) {
	if endpoint == "/job/test-valid-endpoint/allocations" {
		allocs := []nomad.AllocListStub{
			{
				ID:           "test-alloc",
				ClientStatus: "running",
			},
			{
				ID:           "test-alloc-2",
				ClientStatus: "running",
			},
			{
				ID:           "test-alloc-3",
				ClientStatus: "stopped",
			},
		}

		allocsBytes, err := json.Marshal(allocs)
		if err != nil {
			return nil, err
		}

		return allocsBytes, nil
	} else if endpoint == "/job/test-invalid-endpoint/allocations" {
		return nil, fmt.Errorf("error")
	} else if endpoint == "/job/test-no-allocations/allocations" {
		allocs := []nomad.AllocListStub{
			{
				ID:           "test-alloc-3",
				ClientStatus: "stopped",
			},
		}

		allocsBytes, err := json.Marshal(allocs)
		if err != nil {
			return nil, err
		}

		return allocsBytes, nil
	} else if endpoint == "/job/test-valid-endpoint" {
		job := nomad.Job{
			ID: "test",
		}

		jobBytes, err := json.Marshal(job)
		if err != nil {
			return nil, err
		}

		return jobBytes, nil
	} else if endpoint == "/job/test-invalid-endpoint" {
		return nil, fmt.Errorf("error")
	} else if endpoint == "/job/test-invalid-response" {
		return []byte(`{"State": "running"}`), nil
	} else if endpoint == "/job/test-no-jobs" {
		return nil, nil

	} else if endpoint == "/job/test-project-prospector/allocations" {
		allocs := []nomad.AllocListStub{
			{
				ID:           "test-alloc",
				ClientStatus: "running",
			},
			{
				ID:           "test-alloc-2",
				ClientStatus: "running",
			},
			{
				ID:           "test-alloc-3",
				ClientStatus: "stopped",
			},
		}

		allocsBytes, err := json.Marshal(allocs)
		if err != nil {
			return nil, err
		}

		return allocsBytes, nil
	} else if endpoint == "/jobs?meta=true" {
		jobs := []nomad.Job{
			{
				Name:   "test-resource-job-prospector",
				ID:     "test-resource-job-prospector",
				Status: "running",
				TaskGroups: []*nomad.TaskGroup{
					{
						Name: "test-task-group",
						Tasks: []*nomad.Task{
							{
								Name: "test-task",
								Resources: &nomad.Resources{
									CPU:      1,
									MemoryMB: 1,
								},
							},
						},
					},
				},
			},
		}

		jobsBytes, err := json.Marshal(jobs)
		if err != nil {
			return nil, err
		}

		return jobsBytes, nil
	} else if endpoint == "/job/test-resource-job-prospector" {
		println("test-resource-job-prospector")
		jobs := []nomad.Job{
			{
				Name:   "test-resource-job-prospector",
				ID:     "test-resource-job-prospector",
				Status: "running",
				TaskGroups: []*nomad.TaskGroup{
					{
						Name: "test-task-group",
						Tasks: []*nomad.Task{
							{
								Name: "test-task",
								Resources: &nomad.Resources{
									CPU:      1,
									MemoryMB: 1,
								},
							},
						},
					},
				},
			},
		}

		jobsBytes, err := json.Marshal(jobs)
		if err != nil {
			return nil, err
		}

		return jobsBytes, nil
	}

	return nil, nil
}

func (m *MockNomadClient) Post(endpoint string, reqBody *bytes.Buffer) ([]byte, error) {
	switch endpoint {
	case "/client/allocation/test-alloc/restart":
		response := nomad.GenericResponse{}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		return responseBytes, nil
	}

	return nil, nil
}

func (m *MockNomadClient) Delete(endpoint string) ([]byte, error) {
	return nil, nil
}

func (m *MockNomadClient) Forward(ctx *gin.Context, endpoint string) (*http.Response, error) {
	return nil, nil
}
