package cmd

import (
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a job",
	Long: `Restarts a job with the given name. If the --purge flag is set, the job will be purged from the system.

For example:
	prospect job restart my-job`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		addr := cmd.Flag("address").Value.String()
		res, err := CmdPost(addr+"/api/v1/jobs/"+name+"/restart", "")
		if err != nil {
			return
		}

		if res.StatusCode == 404 {
			return
		}

	},
}

func init() {
	projectCmd.AddCommand(restartCmd)

	restartCmd.Args = cobra.ExactArgs(1)
	restartCmd.Flags().StringP("address", "a", "https://prospector.ie", "The address of the Prospector server")
}
