package helpers

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func PerformAuthRequest(r http.Handler, method, path, username, password string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	req.SetBasicAuth(username, password)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func CheckJobHasValidName(ctx *gin.Context, id string) bool {
	if strings.Contains(id, "-prospector") || id != "" {
		return true
	}

	ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid job ID"})
	return false
}
