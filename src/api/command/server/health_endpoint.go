package server

import (
	"fmt"
	"net/http"
)

// @Summary		Healthcheck
// @Description	Healthcheck endpoint for the server
// @Produce		json
// @Success		200	{string}	string	"OK"
// @Router			/health [get]
func (s *HTTPServer) ServerHealthcheck(w http.ResponseWriter, r *http.Request) (err error) {
	// print the request
	fmt.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte("OK"))

	return nil
}
