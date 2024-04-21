package cmd

import (
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "A subcommand for managing projects in the Prospector system",
	Long: `The project subcommand is used to manage prospector projects.
	
	You can use this subcommand to create, delete, and list projects.
	
For example:
	
	prospector project create --name my-project --image my-image --port 8080 --cpu 100 --memory 300`,

	Run:  func(cmd *cobra.Command, args []string) {},
	Args: cobra.MinimumNArgs(1),

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("address")
		if !CheckProspectorReachability(addr) {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.Args = cobra.MinimumNArgs(1)
}
