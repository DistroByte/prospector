package commands

import (
	"testing"

	"net/http"
	"time"
)

func TestUserCommand(t *testing.T) {

	listener := NewServerCommand()

	// start the server
	go func() {
		if err := listener.Run(); err != nil {
			t.Error("failed to start server")
		}
	}()

	// wait for the server to start
	time.Sleep(3 * time.Second)

	// make a request to the server
	resp, err := http.Get("http://localhost:8080/api/users")
	if err != nil {
		t.Error("failed to make request to server")
	}

	// check the response
	if resp.StatusCode != http.StatusOK {
		t.Error("server did not return 200 OK")
	}

	// stop the server
	listener.Shutdown()
}
