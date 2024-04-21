package cmd

import (
	"github.com/spf13/cobra"
)

// projectCmd represents the job command
var projectCmd = &cobra.Command{
	Use:   "job",
	Short: "A subcommand for managing jobs in the Prospector system",
	Long: `The job subcommand is used to manage prospector jobs.
	
	You can use this subcommand to create, delete, and list jobs.
	
For example:
	
	prospector job create --name my-job --image my-image --port 8080 --cpu 100 --memory 300`,

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
