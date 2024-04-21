package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecuteBaseCommand(t *testing.T) {
	actual := new(bytes.Buffer)

	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{""})

	Execute()

	expected := `Prospector is a user management and infrastructure-as-a-service tool
enabling easy on-demand deployment of containers and virtual machines.

Usage:
  prospector [command]

Available Commands:
  auth        Authenticate with the server
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  job         A subcommand for managing jobs in the Prospector system
  server      Start the Prospector API server

Flags:
  -a, --address string   The address of the Prospector server (default "https://prospector.ie")
  -h, --help             help for prospector
  -t, --toggle           Help message for toggle

Use "prospector [command] --help" for more information about a command.
`

	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}

func TestExecuteUnKnownCommand(t *testing.T) {
	actual := new(bytes.Buffer)

	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"unknown"})
	rootCmd.Execute()

	expected := `Error: unknown command "unknown" for "prospector"
Run 'prospector --help' for usage.
`

	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}

func createTestServer(statusCode int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}))
}
