package commands

import (
	"flag"
	"log"
	"net/http"
)

func NewUserCommand() *UserCommand {
	userCmd := &UserCommand{
		flagset: flag.NewFlagSet("user", flag.ExitOnError),
	}

	return userCmd
}

type UserCommand struct {
	flagset *flag.FlagSet
}

func (u *UserCommand) Name() string {
	return u.flagset.Name()
}

func (u *UserCommand) Init(args []string) error {
	return u.flagset.Parse(args)
}

func (u *UserCommand) Run() error {
	log.Println("user")

	// make a http request to the api

	req, err := http.NewRequest("GET", "http://localhost:8080/api/users", nil)

	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	log.Println(resp)

	return nil
}
