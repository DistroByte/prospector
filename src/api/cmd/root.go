package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "prospector",
	Short: "Infrastructure-as-a-service and user management tool for running containers and virtual machines",
	Long: `Prospector is a user management and infrastructure-as-a-service tool
enabling easy on-demand deployment of containers and virtual machines.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
