package command

import (
	"os"
	"prospector/command/server"

	"github.com/mitchellh/cli"
)

type NamedCommand interface {
	Name() string
}

func Commands(metaPtr *Meta, agentUi cli.Ui) map[string]cli.CommandFactory {
	if metaPtr == nil {
		metaPtr = new(Meta)
	}

	meta := *metaPtr

	if meta.Ui == nil {
		meta.Ui = &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		}
	}

	all := map[string]cli.CommandFactory{
		"server": func() (cli.Command, error) {
			return &server.Command{
				Ui:         agentUi,
				ShutdownCh: make(chan struct{}),
			}, nil
		},
	}

	return all
}
