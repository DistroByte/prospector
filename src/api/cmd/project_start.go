package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a job",
	Long: `Starts a job with the given name. If the --purge flag is set, the job will be purged from the system.
	
For example:
	prospect job start my-job
	prospect job start my-job --purge`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		addr := cmd.Flag("address").Value.String()
		res, err := CmdPost(addr+"/api/v1/jobs/"+name+"/start", "")
		if err != nil {
			fmt.Println("Error: The server responded with an error. Please try again.")
			return
		}

		if res.StatusCode == http.StatusNotFound {
			fmt.Println("Error: Job not found")
			return
		}

		fmt.Println("Job started successfully")
	},
}

func init() {
	projectCmd.AddCommand(startCmd)

	startCmd.Args = cobra.ExactArgs(1)
	startCmd.Flags().StringP("address", "a", "https://prospector.ie", "The address of the Prospector server")
}
