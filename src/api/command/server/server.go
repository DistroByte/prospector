package server

import (
	"fmt"
	"sync"
)

type Server struct {
	config *Config

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex

	server *Server
}

func NewServer(config *Config) (*Server, error) {
	s := &Server{
		config:     config,
		shutdownCh: make(chan struct{}),
	}

	if err := s.setupServer(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) setupServer() error {
	conf, err := s.serverConfig()

	if err != nil {
		return fmt.Errorf("server config setup failed: %s", err)
	}

	server, err := NewProspectorServer(conf)

	if err != nil {
		return fmt.Errorf("server setup failed: %v", err)
	}

	s.server = server

	return nil
}

func (s *Server) Shutdown() error {
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		return nil
	}

	if s.server != nil {
		if err := s.server.ShutdownServer(); err != nil {
			fmt.Println("server shutdown failed")
		}
	}

	fmt.Println("Successfully shutdown server")
	s.shutdown = true

	if s.shutdownCh != nil {
		close(s.shutdownCh)
	}

	return nil
}

func (s *Server) ShutdownServer() error {
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		return nil
	}

	s.shutdown = true

	return nil
}

func (s *Server) serverConfig() (*Config, error) {
	conf := &Config{
		BindAddr: s.config.BindAddr,
		Port:     s.config.Port,
	}

	return conf, nil
}

func NewProspectorServer(config *Config) (*Server, error) {
	s := &Server{
		config: config,
	}

	return s, nil
}
