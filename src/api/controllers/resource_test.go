package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func TestGetAllocatedResources(t *testing.T) {
	tcs := []struct {
		name     string
		expected string
		response int
	}{
		{
			name:     "successful response",
			expected: ``,
			response: http.StatusNoContent,
		},
		{
			name:     "empty response",
			expected: ``,
			response: http.StatusNoContent,
		},
		{
			name:     "error response",
			expected: ``,
			response: http.StatusNoContent,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			c := Controller{
				Client: &MockNomadClient{},
			}
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(w)

			claims := jwt.MapClaims{
				c.IdentityKey: "test",
			}

			ctx.Set("JWT_PAYLOAD", claims)

			c.GetAllocatedResources(ctx)
		})
	}
}

func TestGetAllUsedResources(t *testing.T) {
	tcs := []struct {
		name     string
		expected string
		response int
	}{
		{
			name:     "successful response",
			expected: ``,
			response: http.StatusNoContent,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			c := Controller{
				Client: &MockNomadClient{},
			}
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(w)

			claims := jwt.MapClaims{
				c.IdentityKey: "test",
			}

			ctx.Set("JWT_PAYLOAD", claims)

			c.GetAllUsedResources(ctx)
		})
	}
}

func TestGetJobUsedResources(t *testing.T) {
	tcs := []struct {
		name     string
		expected string
		path     string
		response int
	}{
		{
			name:     "successful response",
			expected: ``,
			path:     "/test-job-prospector",
			response: http.StatusNoContent,
		},
		{
			name:     "empty response",
			expected: ``,
			path:     "/test-job-prospector",
			response: http.StatusNoContent,
		},
		{
			name:     "error response",
			expected: ``,
			path:     "/test-job-prospector",
			response: http.StatusNoContent,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			c := Controller{
				Client: &MockNomadClient{},
			}
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: tc.path})

			claims := jwt.MapClaims{
				c.IdentityKey: "test",
			}

			ctx.Set("JWT_PAYLOAD", claims)

			c.GetJobUsedResources(ctx)
		})
	}
}

func TestGetComponentUsedResources(t *testing.T) {
	tcs := []struct {
		name        string
		expected    string
		projectId   string
		componentId string
		response    int
	}{
		{
			name:        "successful response",
			expected:    ``,
			projectId:   "test-project-prospector",
			componentId: "test-component",
			response:    http.StatusNoContent,
		},
		{
			name:        "empty response",
			expected:    ``,
			projectId:   "test-project-prospector",
			componentId: "test-component",
			response:    http.StatusNoContent,
		},
		{
			name:        "error response",
			expected:    ``,
			projectId:   "test-project-prospector",
			componentId: "test-component",
			response:    http.StatusNoContent,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			c := Controller{
				Client: &MockNomadClient{},
			}
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: tc.projectId})
			ctx.Params = append(ctx.Params, gin.Param{Key: "component", Value: tc.componentId})

			claims := jwt.MapClaims{
				c.IdentityKey: "test",
			}

			ctx.Set("JWT_PAYLOAD", claims)

			c.GetComponentUsedResources(ctx)
		})
	}
}
