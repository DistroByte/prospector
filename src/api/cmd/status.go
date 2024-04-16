package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a job",
	Long: `Get the status of a job. Returns the status of a job with the given name.
	
For example:
	prospect job status my-job`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		healthURL := cmd.Flag("address").Value.String() + "/api/health"

		// make the request to the server
		res, err := http.Get(healthURL)
		if err != nil {
			fmt.Println("Error: Could not connect to prospector server, is it running?")
			return
		}

		// check the status code of the response
		if res.StatusCode != 200 {
			fmt.Println("Error: There was an error connecting to the server. Please try again.")
			return
		}

		jobStatusUrl := cmd.Flag("address").Value.String() + "/api/jobs/" + name

		// make the request to the server
		res, err = http.Get(jobStatusUrl)
		if err != nil {
			fmt.Println("Error: The server responded with an error. Please try again.")
			return
		}

		// check the status code of the response
		if res.StatusCode != 200 {
			fmt.Println("Error: There was an error connecting to the server. Please try again.")
			return
		}

		// read the response body
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error: There was an error reading the response body. Please try again.")
			return
		}

		// pretty print the json response body
		fmt.Println(string(body))

		// close the response body
		res.Body.Close()
	},
}

func init() {
	jobCmd.AddCommand(statusCmd)

	statusCmd.Args = cobra.ExactArgs(1)
	statusCmd.Flags().StringP("address", "a", "http://localhost:3434/", "The address of the Prospector server")
}
