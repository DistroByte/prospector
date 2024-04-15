package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"prospector/middleware"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetUnauthenticatedRoutes(t *testing.T) {
	tcs := []struct {
		name   string
		path   string
		status int
	}{
		{
			name:   "valid",
			path:   "/api/health",
			status: http.StatusOK,
		},
		{
			name:   "invalid",
			path:   "/api/unhealthy",
			status: http.StatusUnauthorized,
		},
		{
			name:   "docs",
			path:   "/api/vm-config/notfound",
			status: http.StatusNotFound,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			getHasStatus(t, tc.path, tc.status)
		})
	}
}

func TestPostUnauthenticatedRoutes(t *testing.T) {
	tcs := []struct {
		name    string
		path    string
		body    string
		headers map[string]string
		status  int
	}{
		{
			name:   "valid credentials",
			path:   "/api/login",
			body:   `{"username": "tesla", "password": "password"}`,
			status: http.StatusOK,
		},
		{
			name:   "invalid credentials",
			path:   "/api/login",
			body:   `{"username": "test", "password": "invalid"}`,
			status: http.StatusUnauthorized,
		},
		{
			name:   "no credentials",
			path:   "/api/login",
			body:   `{"username": "", "password": ""}`,
			status: http.StatusUnauthorized,
		},
		{
			name:   "unauthorized user",
			path:   "/api/login",
			body:   `{"username": "test", "password": "password"}`,
			status: http.StatusUnauthorized,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			body := bytes.NewBufferString(tc.body)
			postHasStatus(t, tc.path, body, tc.headers, tc.status)
		})
	}
}

func TestGetAuthenticatedRoutes(t *testing.T) {
	token := getToken(t, "tesla", "password")

	tcs := []struct {
		name    string
		path    string
		headers map[string]string
		body    string
		status  int
	}{
		{
			name: "invalid login",
			path: "/api/login",
			headers: map[string]string{
				"Authorization": "Bearer invalid",
			},
			status: http.StatusUnauthorized,
		},
		{
			name: "valid login",
			path: "/api/v1/user",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			status: http.StatusOK,
		},
		{
			name: "valid login with invalid path",
			path: "/api/v1/invalid",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			status: http.StatusNotFound,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			getAuthenticatedHasStatus(t, tc.path, tc.status, tc.headers)
		})
	}
}

func getHasStatus(t *testing.T, path string, status int) *httptest.ResponseRecorder {
	t.Helper()
	ctx, r, w := getRouter(t)

	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}

	r.ServeHTTP(w, req)

	responseHasCode(t, w, status)

	return w
}

func postHasStatus(t *testing.T, path string, body *bytes.Buffer, headers map[string]string, status int) *httptest.ResponseRecorder {
	t.Helper()
	ctx, r, w := getRouter(t)

	req, err := http.NewRequestWithContext(ctx, "POST", path, body)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	r.ServeHTTP(w, req)

	responseHasCode(t, w, status)

	return w
}

func getAuthenticatedHasStatus(t *testing.T, path string, status int, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	ctx, r, w := getRouter(t)

	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	r.ServeHTTP(w, req)

	responseHasCode(t, w, status)

	return w
}

func getRouter(t *testing.T) (*gin.Context, *gin.Engine, *httptest.ResponseRecorder) {
	t.Helper()
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, r := gin.CreateTestContext(w)
	middleware.CreateAuthMiddlewares(r, "id")
	Route(r, "id")
	return ctx, r, w
}

func responseHasCode(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if expected != w.Code {
		t.Errorf("expected response code %d, got %d", expected, w.Code)
	}
}

func getToken(t *testing.T, username string, password string) string {
	t.Helper()
	body := bytes.NewBufferString(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
	w := postHasStatus(t, "/api/login", body, nil, http.StatusOK)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}

	token := response["token"].(string)

	return token
}

func TestGinReleaseMode(t *testing.T) {
	t.Helper()
	w := httptest.NewRecorder()
	gin.SetMode(gin.ReleaseMode)
	ctx, r := gin.CreateTestContext(w)
	middleware.CreateAuthMiddlewares(r, "id")
	Route(r, "id")

	req, err := http.NewRequestWithContext(ctx, "GET", "/api/health", nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}

	r.ServeHTTP(w, req)

	responseHasCode(t, w, http.StatusOK)
}
