package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a project",
	Long: `Stops a project with the given name. If the --purge flag is set, the project will be purged from the system.
	
For example:
	prospect project stop my-project
	prospect project stop my-project --purge`,
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
			fmt.Println("Error: Project not found")
			return
		}

		fmt.Println("Project stopped successfully")
	},
}

func init() {
	projectCmd.AddCommand(stopCmd)

	stopCmd.Args = cobra.ExactArgs(1)
	stopCmd.Flags().BoolP("purge", "p", false, "Purge the project from the system")
	stopCmd.Flags().StringP("address", "a", "https://prospector.ie", "The address of the Prospector server")
}
