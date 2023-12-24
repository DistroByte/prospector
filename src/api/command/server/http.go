package server

import (
	"net"
	"net/http"
	"strconv"

	_ "prospector/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

type HTTPServer struct {
	listener   net.Listener
	listenerCh chan struct{}
	Addr       string

	mux *http.ServeMux
}

// @title			Prospector API
// @description	API for Prospector, a tool for deploying and managing Nomad jobs
// @version		0
// @host			prospector.ie
// @BasePath		/api
// @Schemes		http https
func NewHTTPServer(config *Config) (*HTTPServer, error) {
	lnAddr, err := net.ResolveTCPAddr("tcp", config.BindAddr+":"+strconv.Itoa(config.Port))
	if err != nil {
		return nil, err
	}

	ln, err := config.Listener("tcp", lnAddr.IP.String(), lnAddr.Port)
	if err != nil {
		return nil, err
	}

	srv := &HTTPServer{
		listener:   ln,
		listenerCh: make(chan struct{}),
		Addr:       ln.Addr().String(),
		mux:        http.NewServeMux(),
	}

	srv.registerHandlers()

	http.Handle("/api/", srv.mux)

	httpServer := &http.Server{
		Addr: srv.Addr,
	}

	go func() {
		defer close(srv.listenerCh)
		httpServer.Serve(ln)
	}()

	return srv, nil
}

func (s *HTTPServer) Shutdown() {
	if s != nil {
		s.listener.Close()
		<-s.listenerCh
	}
}

func (s *HTTPServer) registerHandlers() {
	s.mux.HandleFunc("/api/health", wrapHandler(s.ServerHealthcheck))

	s.mux.HandleFunc("/api/docs/", httpSwagger.WrapHandler)
}

func wrapHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
