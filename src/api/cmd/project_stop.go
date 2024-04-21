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

		addr := cmd.Flag("address").Value.String()
		res, err := CmdDelete(addr+"/api/v1/jobs/"+name+"?purge="+purgeText, "")
		if err != nil {
			fmt.Println("Error: The server responded with an error. Please try again.")
			return
		}

		if res.StatusCode == http.StatusNotFound {
			fmt.Println("Error: Job not found")
			return
		}

		fmt.Println("Job stopped successfully")
	},
}

func init() {
	projectCmd.AddCommand(stopCmd)

	stopCmd.Args = cobra.ExactArgs(1)
	stopCmd.Flags().BoolP("purge", "p", false, "Purge the job from the system")
	stopCmd.Flags().StringP("address", "a", "https://prospector.ie", "The address of the Prospector server")
}
