package cmd

import (
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a project",
	Long: `Restarts a project with the given name. If the --purge flag is set, the project will be purged from the system.

For example:
	prospect project restart my-project`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		addr := cmd.Flag("address").Value.String()
		res, err := CmdPut(addr+"/api/v1/jobs/"+name+"/restart", "")
		if err != nil {
			println(err.Error())
			return
		}

		if res.StatusCode == 404 {
			return
		}

		println("Project restarted successfully")
	},
}

func init() {
	projectCmd.AddCommand(restartCmd)

	restartCmd.Args = cobra.ExactArgs(1)
	restartCmd.Flags().StringP("address", "a", "https://prospector.ie", "The address of the Prospector server")
}
