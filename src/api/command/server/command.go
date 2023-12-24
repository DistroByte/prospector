package server

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/mitchellh/cli"
)

type Command struct {
	Ui         cli.Ui
	ShutdownCh <-chan struct{}

	args   []string
	server *Server

	httpServer *HTTPServer
}

func (c *Command) readConfig() *Config {
	cmdConfig := &Config{
		BindAddr: "",
		Port:     0,
	}

	flags := flag.NewFlagSet("server", flag.ContinueOnError)

	flags.StringVar(&cmdConfig.BindAddr, "bind", "", "")
	flags.IntVar(&cmdConfig.Port, "port", 0, "")

	if err := flags.Parse(c.args); err != nil {
		return nil
	}

	var config *Config

	config = DefaultConfig()
	config = config.Merge(cmdConfig)

	return config
}

func (c *Command) Run(args []string) int {
	c.args = args

	config := c.readConfig()

	if config == nil {
		return 1
	}

	if err := c.setupServer(config); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	defer func() {
		c.server.Shutdown()
	}()

	return c.handleSignals()
}

func (c *Command) handleSignals() int {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Handle signals
	var sig os.Signal
	select {
	case s := <-signalCh:
		sig = s
	case <-c.ShutdownCh:
		sig = os.Interrupt
	}

	c.Ui.Output(fmt.Sprintf("==> Caught signal: %s", sig))

	return 1
}

func (c *Command) setupServer(config *Config) error {
	server, err := NewServer(config)

	if err != nil {
		return fmt.Errorf("server setup failed: %v", err)
	}

	c.server = server

	// print server address

	httpServer, err := NewHTTPServer(config)
	if err != nil {
		server.Shutdown()
		return fmt.Errorf("http server setup failed: %v", err)
	}

	c.Ui.Output(fmt.Sprintf("Server listening on http://%s:%s", c.server.config.BindAddr, strconv.Itoa(c.server.config.Port)))
	c.httpServer = httpServer

	return nil
}

func (c *Command) Synopsis() string {
	return "Run the server"
}

func (c *Command) Help() string {
	return `Usage: prospector server [options]

Starts the prospector server and runs until an interrupt is received.

Configuration for the server is listed below. 

Flags can be set via environment variables by prefixing the flag name 
with PROSPECTOR_ and uppercasing it. For example: PROSPECTOR_PORT=8080

Options:
  -port=<port>  Port to listen on
	
Examples:
  $ prospector server -port=8080`
}
