package cmd

import "testing"

func TestCmdPut(t *testing.T) {
	httpServer := createTestServer(200, "")

	res, err := CmdPut(httpServer.URL, "test")
	if err != nil {
		t.Errorf("Expected nil but got %s", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Expected 200 but got %d", res.StatusCode)
	}

	defer httpServer.Close()
}
