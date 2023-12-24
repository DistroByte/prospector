package command

import (
	"flag"

	"github.com/mitchellh/cli"
)

type FlagSetFlags uint

const (
	FlagSetNone    FlagSetFlags = 0
	FlagSetClient  FlagSetFlags = 1 << iota
	FlagSetDefault              = FlagSetClient
)

type Meta struct {
	Ui       cli.Ui
	flagPort int
}

func (m *Meta) FlagSet(n string, fs FlagSetFlags) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)

	if fs&FlagSetClient != 0 {
		f.IntVar(&m.flagPort, "port", 0, "")
	}

	f.SetOutput(&uiErrorWriter{ui: m.Ui})

	return f
}
