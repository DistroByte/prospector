package cmd

import (
	"bytes"
	"testing"
)

func TestExecuteStopCommand(t *testing.T) {

	actual := new(bytes.Buffer)

	stopCmd.SetOut(actual)
	stopCmd.SetErr(actual)

	httpServer := createTestServer(200, "")

	stopCmd.Flags().Set("address", httpServer.URL)
	defer httpServer.Close()
	stopCmd.Run(stopCmd, []string{"my-job"})

	expected := ``

	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}

func TestExecuteStopCommandError(t *testing.T) {

	actual := new(bytes.Buffer)

	stopCmd.SetOut(actual)
	stopCmd.SetErr(actual)
	stopCmd.Run(stopCmd, []string{"my-job"})

	expected := ``

	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}

func TestExecuteStopCommandPurge(t *testing.T) {

	actual := new(bytes.Buffer)

	stopCmd.SetOut(actual)
	stopCmd.SetErr(actual)

	httpServer := createTestServer(200, "")

	stopCmd.Flags().Set("address", httpServer.URL)
	stopCmd.Flags().Set("purge", "true")
	defer httpServer.Close()
	stopCmd.Run(stopCmd, []string{"my-job"})

	expected := ``

	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}

func TestExecuteStopCommandDownstreamError(t *testing.T) {

	actual := new(bytes.Buffer)

	stopCmd.SetOut(actual)
	stopCmd.SetErr(actual)

	httpServer := createTestServer(500, "")

	stopCmd.Flags().Set("address", httpServer.URL)
	defer httpServer.Close()
	stopCmd.Run(stopCmd, []string{"my-job"})

	expected := ``

	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}
