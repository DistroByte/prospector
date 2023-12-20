package commands

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	// https://brunoscheufler.com/blog/2019-04-26-choosing-the-right-go-web-framework (echo vs others)
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "prospector/docs"
)

//	@title			Prospector API
//	@version		0.0
//	@description	Prospector API

//	@host		prospector.ie
//	@BasePath	/api
//	@schemes	https

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func NewServerCommand() *ServerCommand {
	dc := &ServerCommand{
		flagset: flag.NewFlagSet("server", flag.ExitOnError),
	}

	dc.flagset.StringVar(&dc.port, "port", "8080", "server port")

	return dc
}

type ServerCommand struct {
	flagset *flag.FlagSet
	port    string
}

var listener *echo.Echo

func (d *ServerCommand) Name() string {
	return d.flagset.Name()
}

func (d *ServerCommand) Init(args []string) error {
	return d.flagset.Parse(args)
}

func (d *ServerCommand) Run() error {

	listener = echo.New()

	go func() {
		if err := listener.Start(":" + d.port); err != nil && err != http.ErrServerClosed {
			listener.Logger.Fatal("shutting down the server")
		}
	}()

	listener.GET("/api", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	listener.GET("/api/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	listener.GET("/api/users", func(c echo.Context) error {
		return c.String(http.StatusOK, "users")
	})

	listener.GET("/api/docs/*", echoSwagger.WrapHandler)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := listener.Shutdown(ctx); err != nil {
		listener.Logger.Fatal(err)
	}

	return nil
}

func (d *ServerCommand) Shutdown() error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := listener.Shutdown(ctx); err != nil {
		listener.Logger.Fatal(err)
	}
	return nil
}
