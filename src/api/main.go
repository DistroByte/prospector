package main

import (
	"fmt"
	"os"
	"prospector/command"

	"github.com/mitchellh/cli"
)

// pass all cli args to the base of commands
func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {
	metaPtr := new(command.Meta)

	agentUi := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	commands := command.Commands(metaPtr, agentUi)

	cli := &cli.CLI{
		Name:     "prospector",
		Args:     args,
		Commands: commands,
	}

	exitCode, err := cli.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %s\n", err.Error())
		return 1
	}

	return exitCode
}
