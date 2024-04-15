package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPostUnauthenticatedPaths(t *testing.T) {
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

func responseHasCode(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if expected != w.Code {
		t.Errorf("expected response code %d, got %d", expected, w.Code)
	}
}

func getRouter(t *testing.T) (*gin.Context, *gin.Engine, *httptest.ResponseRecorder) {
	t.Helper()
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, r := gin.CreateTestContext(w)
	CreateAuthMiddlewares(r, "id")

	r.POST("/api/login", AuthMiddleware("id").LoginHandler)
	return ctx, r, w
}
