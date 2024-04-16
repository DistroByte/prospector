package controllers

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockNomadClientJobs struct{}

func (m *MockNomadClientJobs) Get(path string) ([]byte, error) {
	switch path {
	case "/jobs?meta=true":
		return []byte(`{"meta": {"count": 1}}`), nil
	case "/jobs/test-valid-endpoint":
		return []byte(`{"ID": "test-valid-endpoint"}`), nil
	case "/jobs/test-invalid-endpoint":
		return nil, nil
	case "/jobs/test-invalid-response":
		return nil, nil
	case "/jobs/test-no-job":
		return nil, nil
	}

	return nil, fmt.Errorf("error")
}

func (m *MockNomadClientJobs) Post(path string, reqBody *bytes.Buffer) ([]byte, error) {
	return nil, nil
}

func (m *MockNomadClientJobs) Delete(path string) ([]byte, error) {
	return nil, nil
}

func TestGetJobs(t *testing.T) {
	c := Controller{
		Client: &MockNomadClientJobs{},
	}

	tcs := []struct {
		name   string
		path   string
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
			expect: nil,
		},
		{
			name:   "no jobs",
			path:   "test-no-jobs",
			expect: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: tc.path})

			c.GetJobs(ctx)

			if w.Code != 204 {
				t.Errorf("expected status code 200, got %v", w.Code)
			}

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
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: tc.path})

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
