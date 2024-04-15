package controllers

import "testing"

func TestToJson(t *testing.T) {
	project := Project{
		Name: "Test Project",
		Type: "Test Type",
		Components: []Component{
			{
				Name:  "Test Component",
				Image: "Test Image",
				Resources: Resources{
					Cpu:    1,
					Memory: 1,
				},
				Network: Network{
					Port:   1,
					Expose: false,
				},
			},
		},
	}

	json := project.ToJson()
	if json == nil {
		t.Errorf("expected json to not be nil")
	}

	if json.String() != `{"name":"Test Project","type":"Test Type","components":[{"name":"Test Component","image":"Test Image","resources":{"cpu":1,"memory":1},"network":{"port":1,"expose":false},"user_config":{"user":"","ssh_key":""}}],"user":""}` {
		t.Errorf("expected json to be %s, got %s", `{"name":"Test Project","type":"Test Type","components":[{"name":"Test Component","image":"Test Image","resources":{"cpu":1,"memory":1},"network":{"port":1,"expose":false},"user_config":{"user":"","ssh_key":""}}],"user":""}`, json.String())
	}
}
