package helpers

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
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

func CheckJobHasValidName(id string) bool {
	if strings.Contains(id, "-prospector") && id != "" {
		return true
	}

	return false
}
