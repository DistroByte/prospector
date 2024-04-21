package cmd

import (
	"bytes"
	"testing"
)

func TestProjectRestartExecute(t *testing.T) {
	actual := new(bytes.Buffer)

	restartCmd.SetOut(actual)
	restartCmd.SetErr(actual)

	httpServer := createTestServer(200, "")
	restartCmd.Flags().Set("address", httpServer.URL)
	defer httpServer.Close()
	restartCmd.Run(restartCmd, []string{"test"})

	expected := ``
	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}
