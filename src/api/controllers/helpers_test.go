package controllers

import "testing"

func TestCreateJobFromTemplate(t *testing.T) {
	// test cases
	tcs := []struct {
		name      string
		job       Project
		jobSource string
		expect    int
	}{
		{
			name: "valid docker",
			job: Project{
				Name: "test",
				Type: "docker",
				Components: []Component{
					{
						Name:  "test",
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
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			status, err := CreateJobFromTemplate(tc.job, tc.jobSource)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if status != tc.expect {
				t.Errorf("expected %d, got %d", tc.expect, status)
			}
		})
	}
}
