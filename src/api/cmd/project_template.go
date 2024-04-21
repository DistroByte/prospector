package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"reflect"

	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "A subcommand for creating a project template",
	Long: `The template subcommand is used to create a project template.

For example:

	prospector template create --name my-template --components 2 --type docker`,

	Run: func(cmd *cobra.Command, args []string) {
		var dockerTemplateSource = `{
	"components": [
		{{ range $i, $component := .Components }}{
			"image": "string",
			"name": "string",
			"network": {
				"expose": false,
				"port": 65535
			},
			"resources": {
				"cpu": 100,
				"memory": 100
			},
			"user_config": {
				"ssh_key": "string",
				"user": "string"
			},
			"volumes": [
				"string"
			]
		}{{ if not (last $i $.Components) }},{{ end }}{{ end }}
	],
	"name": "{{ .Name}}",
	"type": "{{ .Type }}"
}`

		var vmTemplateSource = `{
	"components": [
		{{ range $i, $component := .Components }}{
			"image": "string",
			"resources": {
				"cpu": 1000,
				"memory": 500
			},
			"user_config": {
				"ssh_key": "string",
				"user": "string"
			},
			"volumes": [
				"string"
			]
		}{{ if not (last $i $.Components) }},{{ end }}
		{{ end }}
	],
	"name": "{{ .Name}}",
	"type": "{{ .Type }}"
}`

		var templateSource string

		templateName, _ := cmd.Flags().GetString("name")
		componentCount, _ := cmd.Flags().GetInt("components")
		templateType, _ := cmd.Flags().GetString("type")

		if templateName == "" {
			fmt.Println("Error: The name of the template is required. Please try again.")
			return
		}

		if componentCount < 1 {
			componentCount = 1
		}

		if templateType != "docker" && templateType != "vm" {
			fmt.Println("Error: The type of the template must be either 'docker' or 'vm'. Please try again.")
			return
		}

		data := struct {
			Components []int
			Name       string
			Type       string
		}{
			Components: make([]int, componentCount),
			Name:       templateName,
			Type:       templateType,
		}

		if templateType == "docker" {
			templateSource = dockerTemplateSource
		} else {
			templateSource = vmTemplateSource
		}

		t, err := template.New("template").Funcs(template.FuncMap{"last": last}).Parse(templateSource)
		if err != nil {
			fmt.Println("Error: The template could not be parsed. Please try again.")
			return
		}

		var tpl bytes.Buffer
		err = t.Execute(&tpl, data)
		if err != nil {
			fmt.Println("Error: The template could not be executed. Please try again.")
			return
		}

		err = os.WriteFile(templateName+".json", tpl.Bytes(), 0644)
		if err != nil {
			fmt.Println("Error: The template could not be written. Please try again.")
			return
		}

		fmt.Println("Template created successfully")
	},
}

func init() {
	projectCmd.AddCommand(templateCmd)

	templateCmd.Flags().StringP("name", "n", "", "The name of the template")
	templateCmd.Flags().IntP("components", "c", 1, "The number of components in the template")
	templateCmd.Flags().StringP("type", "t", "docker", "The type of the template")
}

func last(i int, slice interface{}) bool {
	v := reflect.ValueOf(slice)
	return i == v.Len()-1
}
