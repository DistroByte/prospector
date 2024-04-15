package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMetrics(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	_, r := gin.CreateTestContext(w)
	CreateAuthMiddlewares(r, "id")
	MetricsMiddleware(r)
}
