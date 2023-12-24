package server

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/nomad/testutil"
)

type TestServer struct {
	T          testing.TB
	Name       string
	Config     *Config
	HTTPServer *HTTPServer

	*Server
}

func NewTestServer(t testing.TB, name string, cb func(c *Config)) *TestServer {
	s := &TestServer{
		T:    t,
		Name: name,
	}

	s.Start()

	return s
}

func (s *TestServer) Start() *TestServer {
	if s.Config == nil {
		s.Config = DefaultConfig()
	}

	i := 10

	advertiseAddrs := s.Config.BindAddr

RETRY:
	i--

	newAddrs := advertiseAddrs

	s.Config.BindAddr = newAddrs

	server, err := s.start()
	if err == nil {
		s.Server = server
	} else if i == 0 {
		s.T.Fatalf("Failed to start server: %s", err)
	} else {
		if server != nil {
			server.Shutdown()
		}

		wait := time.Duration(rand.Int31n(2000)) * time.Millisecond
		s.T.Logf("Retrying server start in %s", wait)
		time.Sleep(wait)

		goto RETRY
	}

	failed := false
	testutil.WaitForResult(func() (bool, error) {
		req, _ := http.NewRequest("GET", "http://"+s.HTTPServer.Addr+"/api/health", nil)
		resp := httptest.NewRecorder()
		err := s.HTTPServer.ServerHealthcheck(resp, req)
		return err == nil && resp.Code == 200, err
	}, func(err error) {
		failed = true
		s.T.Logf("Server healthcheck failed: %s", err)
	})

	if failed {
		s.T.Fatalf("Failed to start server")
	}

	return s
}

func (s *TestServer) start() (*Server, error) {
	server, err := NewServer(s.Config)
	if err != nil {
		return nil, err
	}

	httpServer, err := NewHTTPServer(s.Config)
	if err != nil {
		return nil, err
	}

	s.HTTPServer = httpServer

	return server, nil
}
