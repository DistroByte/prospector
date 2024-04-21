package cmd

import (
	"bytes"
	"testing"
)

func TestExecuteListCommand(t *testing.T) {

	actual := new(bytes.Buffer)

	listCmd.SetOut(actual)
	listCmd.SetErr(actual)

	httpServer := createTestServer(200, "")

	listCmd.Flags().Set("address", httpServer.URL)
	defer httpServer.Close()
	listCmd.Run(listCmd, []string{})

	expected := ``

	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}

func TestExecuteListCommandNoJobs(t *testing.T) {

	actual := new(bytes.Buffer)

	listCmd.SetOut(actual)
	listCmd.SetErr(actual)

	httpServer := createTestServer(200, `[]`)
	listCmd.Flags().Set("address", httpServer.URL)
	defer httpServer.Close()
	listCmd.Run(listCmd, []string{})

	expected := ``

	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}
