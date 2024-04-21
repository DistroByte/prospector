package cmd

import (
	"encoding/json"
	"fmt"
	"prospector/controllers"
	"time"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all running jobs on the server",
	Long: `List all running jobs on the server. This command will make a request to the server and print out all the running jobs.
	
For example:
    prospector job list`,
	Run: func(cmd *cobra.Command, args []string) {
		addr := cmd.Flag("address").Value.String()
		res, err := CmdGet(addr + "/api/v1/jobs")
		if err != nil {
			println(err.Error())
			fmt.Println("Error: The server responded with an error. Please try again.")
			return
		}

		// check the status code of the response
		if res.StatusCode == 204 {
			fmt.Println("No jobs found")
			return
		}

		var response []controllers.ShortJob
		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			fmt.Println("Error: Could not decode the response. Please try again.")
			return
		}

		// print out the response
		for _, job := range response {
			println("------------")
			fmt.Printf("Name: %s\n", job.ID)
			fmt.Printf("Status: %s\n", job.Status)
			fmt.Printf("Type: %s\n", job.Type)
			fmt.Printf("Created: %s\n", time.Unix(0, job.Created).Format(time.RFC1123Z))
		}
	},
}

func init() {
	projectCmd.AddCommand(listCmd)
}
