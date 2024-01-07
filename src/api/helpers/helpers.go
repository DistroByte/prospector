package helpers

import (
	"net/http"
	"net/http/httptest"
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
