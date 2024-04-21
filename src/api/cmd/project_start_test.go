package cmd

import (
	"bytes"
	"testing"
)

func TestProjectStartExecute(t *testing.T) {
	actual := new(bytes.Buffer)

	startCmd.SetOut(actual)
	startCmd.SetErr(actual)

	httpServer := createTestServer(200, "")
	startCmd.Flags().Set("address", httpServer.URL)
	defer httpServer.Close()
	startCmd.Run(startCmd, []string{"test"})

	expected := ``
	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}
