/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A subcommand for creating a job in the Prospector system",
	Long: `The create subcommand is used to create a prospector job.
	
For example:
	
	prospector job create --name my-job --image my-image --port 8080 --cpu 100 --memory 300`,
	Run: func(cmd *cobra.Command, args []string) {
		// name, _ := cmd.Flags().GetString("name")
		// image, _ := cmd.Flags().GetString("image")
		// application_port, _ := cmd.Flags().GetInt("application-port")
		// cpu, _ := cmd.Flags().GetInt("cpu")
		// memory, _ := cmd.Flags().GetInt("memory")

		// healthURL := "http://localhost:3434/api/health"

		// // make the request to the server
		// res, err := http.Get(healthURL)
		// if err != nil {
		// 	fmt.Println("Error: Could not connect to prospector server, is it running?")
		// 	return
		// }

		// // check the status code of the response
		// if res.StatusCode != 200 {
		// 	fmt.Println("Error: There was an error connecting to the server. Please try again.")
		// 	return
		// }

		// jobCreateUrl := "http://localhost:3434/api/jobs"

		// var job controllers.Job = controllers.Job{
		// 	ProjectName: name,
		// 	Component: controllers.Component{
		// 		Name: name,
		// 		Type: "docker",
		// 	},
		// 	Image: image,
		// 	Resources: controllers.Resources{
		// 		Cpu:    cpu,
		// 		Memory: memory,
		// 	},
		// 	Network: controllers.Network{
		// 		Port:   application_port,
		// 		Expose: true,
		// 	},
		// }

		// // make the request to the server
		// res, err = http.Post(jobCreateUrl, "application/json", job.ToJson())
		// if err != nil {
		// 	fmt.Println("Error: The server responded with an error. Please try again.")
		// 	return
		// }

		// // check the status code of the response
		// if res.StatusCode != 200 {
		// 	fmt.Println("Error: The server responded with an error. Please try again.")
		// 	return
		// }

		// fmt.Printf("Creating job %s\n", name)
	},
}

func init() {
	jobCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("name", "n", "", "Name of the job")
	createCmd.Flags().StringP("image", "i", "", "Docker image to use for the job")
	createCmd.Flags().IntP("application-port", "p", 0, "Port the application is bound to inside the container")
	createCmd.Flags().IntP("cpu", "c", 100, "CPU to allocate to the job")
	createCmd.Flags().IntP("memory", "m", 100, "Memory to allocate to the job")

	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("image")
	createCmd.MarkFlagRequired("application-port")
}
