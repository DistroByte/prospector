package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	nomad "github.com/hashicorp/nomad/nomad/structs"
)

type MockNomadClientJobs struct{}

func (m *MockNomadClientJobs) Get(path string) ([]byte, error) {
	switch path {
	case "/jobs?meta=true":
		job := []nomad.JobListStub{
			{
				ID:     "test-nginx-prospector",
				Name:   "test-nginx-prospector",
				Status: "stopped",
			},
			{
				ID:     "test-nginx-prospector-2",
				Name:   "test-nginx-prospector-2",
				Status: "running",
			},
		}

		jobBytes, err := json.Marshal(job)
		if err != nil {
			return nil, err
		}

		return jobBytes, nil
	case "/jobs/test-valid-endpoint":
		return []byte(`{"ID": "test-valid-endpoint"}`), nil
	case "/jobs/test-invalid-endpoint":
		return nil, fmt.Errorf("error")
	case "/jobs/test-invalid-response":
		return nil, nil
	case "/jobs/test-no-job":
		return nil, nil
	case "/job/test-valid-endpoint-prospector":
		return []byte(`{"ID": "test-valid-endpoint"}`), nil
	case "/job/valid-endpoint-prospector":
		return []byte(`{"ID": "valid-endpoint-prospector"}`), nil
	}

	return nil, fmt.Errorf("error")
}

func (m *MockNomadClientJobs) Post(path string, reqBody *bytes.Buffer) ([]byte, error) {
	return nil, nil
}

func (m *MockNomadClientJobs) Delete(path string) ([]byte, error) {
	return nil, nil
}

func (m *MockNomadClientJobs) Forward(ctx *gin.Context, path string) (*http.Response, error) {
	return nil, nil
}

func TestGetJobs(t *testing.T) {
	c := Controller{
		Client: &MockNomadClientJobs{},
	}

	tcs := []struct {
		name   string
		path   string
		query  string
		expect error
	}{
		{
			name:   "valid jobs",
			path:   "test-valid-endpoint",
			expect: nil,
		},
		{
			name:   "invalid jobs",
			path:   "test-invalid-endpoint",
			expect: nil,
		},
		{
			name:   "invalid response",
			path:   "test-invalid-response",
			query:  "?running=true",
			expect: nil,
		},
		{
			name:   "no jobs",
			path:   "test-no-jobs",
			query:  "?long=true",
			expect: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(w)
			claims := jwt.MapClaims{
				c.IdentityKey: "test",
			}

			if tc.query != "" {
				ctx.Request = httptest.NewRequest("GET", "/v1/jobs"+tc.query, nil)
			}

			ctx.Set("JWT_PAYLOAD", claims)
			c.GetJobs(ctx)

			// if w.Code != 204 {
			// 	t.Errorf("expected status code 200, got %v", w.Code)
			// }

			if tc.expect != nil {
				t.Errorf("expected error to be nil, got %v", tc.expect)
			}

		})
	}
}

func TestGetJob(t *testing.T) {
	c := Controller{
		Client: &MockNomadClientJobs{},
	}

	tcs := []struct {
		name   string
		path   string
		expect error
	}{
		{
			name:   "valid job",
			path:   "test-valid-endpoint",
			expect: nil,
		},
		{
			name:   "invalid job",
			path:   "test-invalid-endpoint",
			expect: nil,
		},
		{
			name:   "invalid response",
			path:   "test-invalid-response",
			expect: nil,
		},
		{
			name:   "no job",
			path:   "test-no-job",
			expect: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			claims := jwt.MapClaims{
				c.IdentityKey: "test",
			}
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: tc.path})
			ctx.Set("JWT_PAYLOAD", claims)

			c.GetJob(ctx)

			if w.Code != 200 {
				t.Errorf("expected status code 200, got %v", w.Code)
			}

			if tc.expect != nil {
				t.Errorf("expected error to be nil, got %v", tc.expect)
			}
		})
	}
}

func TestGetJobDefinition(t *testing.T) {
	c := Controller{
		Client: &MockNomadClientJobs{},
	}

	tcs := []struct {
		name   string
		path   string
		expect error
	}{
		{
			name:   "valid job",
			path:   "test-valid-endpoint-prospector",
			expect: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: tc.path})

			claims := jwt.MapClaims{
				c.IdentityKey: "test",
			}
			ctx.Set("JWT_PAYLOAD", claims)
			c.GetJobDefinition(ctx)

			if w.Code != 500 {
				t.Errorf("expected status code 500, got %v", w.Code)
			}

			if tc.expect != nil {
				t.Errorf("expected error to be nil, got %v", tc.expect)
			}
		})
	}
}

func TestCreateJob(t *testing.T) {
	c := Controller{
		Client: &MockNomadClientJobs{},
	}

	tcs := []struct {
		name    string
		path    string
		project Project
		expect  error
	}{
		{
			name: "valid docker job",
			path: "job-create-valid-docker",
			project: Project{
				Name: "create-valid-endpoint",
				Type: "docker",
				Components: []Component{
					{
						Name:  "component",
						Image: "component",
						Volumes: []string{
							"test",
						},
						Resources: Resources{
							Cpu:    100,
							Memory: 100,
						},
						Network: Network{
							Port:   80,
							Expose: true,
						},
					},
				},
			},
			expect: nil,
		},
		{
			name: "valid vm job",
			path: "job-create-valid-vm",
			project: Project{
				Name: "create-valid-endpoint",
				Type: "vm",
				Components: []Component{
					{
						Name:  "component",
						Image: "component",
						Volumes: []string{
							"test",
						},
						Resources: Resources{
							Cpu:    100,
							Memory: 100,
						},
						Network: Network{
							Port:   80,
							Expose: true,
						},
					},
				},
			},
			expect: nil,
		},
		{
			name: "invalid job type",
			path: "job-create-invalid-type",
			project: Project{
				Name: "create-invalid-type",
				Type: "invalid",
				Components: []Component{
					{
						Name:  "component",
						Image: "component",
						Volumes: []string{
							"test",
						},
						Resources: Resources{
							Cpu:    100,
							Memory: 100,
						},
						Network: Network{
							Port:   80,
							Expose: true,
						},
					},
				},
			},
			expect: fmt.Errorf("invalid job type"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: tc.path})

			body, err := json.Marshal(tc.project)
			if err != nil {
				t.Errorf("failed to marshal project: %v", err)
			}

			ctx.Request = httptest.NewRequest("POST", "/v1/jobs", bytes.NewReader(body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			claims := jwt.MapClaims{
				c.IdentityKey: "test",
			}
			ctx.Set("JWT_PAYLOAD", claims)
			c.CreateJob(ctx)
		})
	}
}

func TestDeleteJob(t *testing.T) {
	c := Controller{
		Client: &MockNomadClientJobs{},
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
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: tc.path})

			claims := jwt.MapClaims{
				c.IdentityKey: "test",
			}
			ctx.Set("JWT_PAYLOAD", claims)
			c.DeleteJob(ctx)
		})
	}
}

func TestRestartJob(t *testing.T) {
	c := Controller{
		Client: &MockNomadClientJobs{},
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
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: tc.path})

			claims := jwt.MapClaims{
				c.IdentityKey: "test",
			}
			ctx.Set("JWT_PAYLOAD", claims)
			c.RestartJob(ctx)
		})
	}
}

func TestStartJob(t *testing.T) {
	c := Controller{
		Client: &MockNomadClientJobs{},
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
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: tc.path})

			claims := jwt.MapClaims{
				c.IdentityKey: "test",
			}
			ctx.Set("JWT_PAYLOAD", claims)
			c.StartJob(ctx)
		})
	}
}
