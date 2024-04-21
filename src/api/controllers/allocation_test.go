package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRestartAlloc(t *testing.T) {
	tcs := []struct {
		name        string
		projectId   string
		componentId string
		response    int
	}{
		{
			name:        "valid job and alloc ID",
			projectId:   "test-project-prospector",
			componentId: "test-alloc",
			response:    http.StatusOK,
		},
		{
			name:        "valid job id with no allocations",
			projectId:   "test-project-lost-prospector",
			componentId: "test-alloc",
			response:    http.StatusOK,
		},
		{
			name:        "invalid job ID",
			projectId:   "invalid",
			componentId: "test-alloc",
			response:    http.StatusForbidden,
		},
		{
			name:        "empty alloc ID",
			projectId:   "test-project-prospector",
			componentId: "",
			response:    http.StatusBadRequest,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			c := Controller{
				Client: &MockNomadClient{},
			}
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.PUT("/:id/component/:component/restart", c.RestartAlloc)

			req, err := http.NewRequest("PUT", "/"+tc.projectId+"/component/"+tc.componentId+"/restart", nil)

			if tc.response == http.StatusForbidden && err != nil {
				return
			}

			if tc.response == http.StatusBadRequest && err != nil {
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.response {
				t.Errorf("Expected status code %d, got %d", tc.response, w.Code)
			}
		})
	}
}

func TestGetComponents(t *testing.T) {
	c := Controller{
		Client: &MockNomadClient{},
	}

	tcs := []struct {
		name      string
		projectId string
		response  int
	}{
		{
			name:      "valid project ID",
			projectId: "test-project-prospector",
			response:  http.StatusOK,
		},
		{
			name:      "invalid project ID",
			projectId: "invalid",
			response:  http.StatusForbidden,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: tc.projectId})

			c.GetComponents(ctx)

			if w.Code != tc.response {
				t.Errorf("Expected status code %d, got %d", tc.response, w.Code)
			}
		})
	}
}
