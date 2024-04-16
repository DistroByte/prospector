package cmd

import (
	"testing"
)

func TestServerCommand(t *testing.T) {
	// call `server` function and kill the server once it has started
	go func() {
		if serverCmd.Runnable() {
			serverCmd.Run(serverCmd, []string{})
		}
	}()
}
