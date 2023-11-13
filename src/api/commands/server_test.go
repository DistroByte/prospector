package commands

import (
	"net/http"
	"testing"
	"time"
)

func TestServerStarts(t *testing.T) {
	cmd := NewServerCommand()

	// start the listener
	go func() {
		if err := cmd.Run(); err != nil {
			t.Error("failed to start server")
		}
	}()

	// wait for the server to start
	time.Sleep(3 * time.Second)

	// make a request to the server
	resp, err := http.Get("http://localhost:8080/api/health")
	if err != nil {
		t.Error("failed to make request to server")
	}

	// check the response
	if resp.StatusCode != http.StatusOK {
		t.Error("server did not return 200 OK")
	}

	// stop the server
	cmd.Shutdown()
}
