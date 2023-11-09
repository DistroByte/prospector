package commands

import (
	"flag"
	"log"
)

func NewUserCommand() *UserCommand {
	uc := &UserCommand{
		fs: flag.NewFlagSet("user", flag.ExitOnError),
	}

	return uc
}

type UserCommand struct {
	fs *flag.FlagSet
}

func (u *UserCommand) Name() string {
	return u.fs.Name()
}

func (u *UserCommand) Init(args []string) error {
	return u.fs.Parse(args)
}

func (u *UserCommand) Run() error {
	log.Println("user")
	return nil
}
