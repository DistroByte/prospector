package commands

import (
	"flag"

	"github.com/labstack/echo/v4"
	// https://brunoscheufler.com/blog/2019-04-26-choosing-the-right-go-web-framework (echo vs others)
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "prospector/docs"
)

//	@title		Prospector API
//	@version	0.0

//	@host		https://prospector.ie
//	@BasePath	/api

func NewServerCommand() *ServerCommand {
	dc := &ServerCommand{
		fs: flag.NewFlagSet("server", flag.ExitOnError),
	}

	dc.fs.StringVar(&dc.port, "port", "8080", "server port")

	return dc
}

type ServerCommand struct {
	fs   *flag.FlagSet
	port string
}

func (d *ServerCommand) Name() string {
	return d.fs.Name()
}

func (d *ServerCommand) Init(args []string) error {
	return d.fs.Parse(args)
}

func (d *ServerCommand) Run() error {
	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start("0.0.0.0:" + d.port))
	// https://github.com/swaggo/echo-swagger#canonical-example
	return nil
}
