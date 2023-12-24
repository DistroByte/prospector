package server

import (
	"net/http"
	"testing"
)

func makeHTTPServer(t testing.TB) *TestServer {
	return NewTestServer(t, "http", func(c *Config) {
		c.Port = 3434
	})
}

func TestHTTPServer(t *testing.T) {
	s := makeHTTPServer(t)
	defer s.Shutdown()

	if s.HTTPServer == nil {
		t.Fatalf("Expected HTTP server to be non-nil")
	}

	if s.HTTPServer.Addr != "[::]:3434" {
		t.Fatalf("Unexpected HTTP address: %s", s.HTTPServer.Addr)
	}

	// Make a request to the server
	resp, err := http.Get("http://0.0.0.0:3434/api/health")
	if err != nil {
		t.Fatalf("Error making request: %s", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

}
