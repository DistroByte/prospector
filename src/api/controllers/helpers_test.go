package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
)

func TestCreateJobFromTemplate(t *testing.T) {
	// test cases
	tcs := []struct {
		name      string
		job       Project
		jobSource string
		expect    int
	}{
		{
			name: "valid container",
			job: Project{
				Name: "test-valid-container",
				Type: "docker",
				Components: []Component{
					{
						Name:  "test-valid-container",
						Image: "nginx",
						Resources: Resources{
							Cpu:    10,
							Memory: 30,
						},
						Network: Network{
							Port:   80,
							Expose: true,
						},
					},
				},
			},
			jobSource: DockerSource,
			expect:    200,
		},
		{
			name: "valid vm",
			job: Project{
				Name: "test-valid-vm",
				Type: "vm",
				Components: []Component{
					{
						Name:  "test-vm",
						Image: "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-genericcloud-amd64.qcow2",
						Resources: Resources{
							Cpu:    1000,
							Memory: 600,
						},
						UserConfig: UserConfig{
							User:   "distro",
							SSHKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIDg4A77LUlRC9xiijAdNWgZFElCXkyickyh6g/FpltK distro@poseidon",
						},
						Network: Network{
							Mac: fmt.Sprintf("00:16:3e:%02x:%02x:%02x", rand.Intn(255), rand.Intn(255), rand.Intn(255)),
						},
					},
				},
			},
			jobSource: VMSource,
			expect:    200,
		},
		{
			name: "invalid container",
			job: Project{
				Name: "test-invalid-container",
				Type: "docker",
				Components: []Component{
					{
						Name:  "test-invalid-container",
						Image: "nginx",
						Resources: Resources{
							Cpu:    10,
							Memory: 30,
						},
						Network: Network{
							Port:   80,
							Expose: true,
						},
					},
				},
			},
			jobSource: "invalid",
			expect:    500,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			status, err := CreateJobFromTemplate(tc.job, tc.jobSource)
			if status != tc.expect {
				t.Errorf("expected %v but got %v with error %v", tc.expect, status, err)
			}

			client := &http.Client{}
			req, _ := http.NewRequest("DELETE", "http://zeus.internal:4646/v1/job/-"+tc.job.Name+"-prospector?purge=true", nil)

			response, err := client.Do(req)
			if err != nil {
				t.Errorf("error deleting job: %v", err)
				return
			}

			if response.StatusCode != 200 {
				t.Errorf("expected 200 but got %v", response.StatusCode)
			}

		})
	}
}

func TestWriteTextFilesForVM(t *testing.T) {
	tcs := []struct {
		name   string
		job    Project
		expect error
	}{
		{
			name: "valid vm",
			job: Project{
				Name: "test-valid-vm",
				Type: "vm",
				Components: []Component{
					{
						Name:  "test-vm",
						Image: "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-genericcloud-amd64.qcow2",
						Resources: Resources{
							Cpu:    1000,
							Memory: 600,
						},
						UserConfig: UserConfig{
							User:   "distro",
							SSHKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIDg4A77LUlRC9xiijAdNWgZFElCXkyickyh6g/FpltK distro@poseidon",
						},
					},
				},
				User: "distro",
			},
			expect: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := WriteTextFilesForVM(tc.job)
			if err != tc.expect {
				t.Errorf("expected %v but got %v", tc.expect, err)
			}

			err = os.RemoveAll("vm-config")
			if err != nil {
				t.Errorf("error removing directory: %v", err)
			}
		})
	}
}

func TestMakeDirsAndFiles(t *testing.T) {
	tcs := []struct {
		name      string
		directory string
		file      string
		expect    bool
	}{
		{
			name:      "valid directory and file",
			directory: "test dir",
			file:      "test.txt",
			expect:    false,
		},
		{
			name:      "invalid directory",
			directory: "./test",
			file:      "test.txt",
			expect:    false,
		},
		{
			name:      "invalid file",
			directory: "test",
			file:      "/test txt",
			expect:    false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, _ := makesDirsAndFiles(tc.file, tc.directory)
			if result != tc.expect {
				t.Errorf("expected %v but got %v", tc.expect, result)
			}

			// clean up
			if !result {
				err := os.RemoveAll(tc.directory)
				if err != nil {
					t.Errorf("error removing directory: %v", err)
				}
			}
		})
	}
}
