package cmd

import (
	"bytes"
	"testing"
)

// authCmd represents the auth command
func TestAuthCommand(t *testing.T) {
	actual := new(bytes.Buffer)

	authCmd.SetOut(actual)
	authCmd.SetErr(actual)

	httpServer := createTestServer(200, "")

	authCmd.Flags().Set("address", httpServer.URL)
	authCmd.Flags().Set("username", "tesla")
	authCmd.Flags().Set("password", "password")
	defer httpServer.Close()
	authCmd.Run(authCmd, []string{})

	expected := ``

	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}
