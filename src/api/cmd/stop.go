package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a job",
	Long: `Stops a job with the given name. If the --purge flag is set, the job will be purged from the system.
	
For example:
	prospect job stop my-job
	prospect job stop my-job --purge`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		purge, _ := cmd.Flags().GetBool("purge")

		purgeText := "false"
		if purge {
			purgeText = "true"
		}

		healthURL := cmd.Flag("address").Value.String() + "/api/health"

		// make the request to the server
		res, err := http.Get(healthURL)
		if err != nil {
			fmt.Println("Error: Could not connect to prospector server, is it running?")
			return
		}

		// check the status code of the response
		if res.StatusCode != 200 {
			fmt.Println("Error: The server returned an error. Please try again.")
			return
		}

		jobStopUrl := cmd.Flag("address").Value.String() + "/api/jobs/" + name + "?purge=" + purgeText

		// make the request to the server
		req, err := http.NewRequest(http.MethodDelete, jobStopUrl, nil)
		if err != nil {
			fmt.Println("Error: The server responded with an error. Please try again.")
			return
		}

		res, err = http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error: The server responded with an error. Please try again.")
			return
		}

		// check the status code of the response
		if res.StatusCode != 200 {
			fmt.Println("Error: The server returned an error. Please try again.")
			return
		}

		fmt.Println("Job stopped successfully")
	},
}

func init() {
	jobCmd.AddCommand(stopCmd)

	stopCmd.Args = cobra.ExactArgs(1)
	stopCmd.Flags().BoolP("purge", "p", false, "Purge the job from the system")
	stopCmd.Flags().StringP("address", "a", "http://localhost:3434", "The address of the Prospector server")
}
