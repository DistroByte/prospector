package main

import (
	"fmt"
	"log"
	"os"

	"prospector/commands"
)

// Runner is an interface for a command line subcommand
type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

// root is the entrypoint for the CLI
func root(args []string) error {
	// Check if we have a subcommand
	if len(args) < 1 {
		return fmt.Errorf("expected subcommand")
	}

	cmds := []Runner{
		commands.NewServerCommand(),
		commands.NewUserCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	return fmt.Errorf("unknown subcommand %s", subcommand)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := root(os.Args[1:]); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
