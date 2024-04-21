package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A subcommand for creating a job in the Prospector system",
	Long: `The create subcommand is used to create a prospector job.
	
For example:
	
	prospector job create job.json`,
	Run: func(cmd *cobra.Command, args []string) {
		buffer, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Println("Error: The file could not be read. Please try again.")
			return
		}

		addr := cmd.Flag("address").Value.String()
		res, err := CmdPost(addr+"/api/v1/jobs", string(buffer))
		if err != nil {
			fmt.Println("Error: The server responded with an error. Please try again.")
			return
		}

		if res.StatusCode == 400 {
			fmt.Println("Error: The job could not be created. Please try again.")
			return
		}

		if res.StatusCode == 200 {
			fmt.Println("Job created successfully")
		} else {
			fmt.Println("Error: The server responded with an error. Please try again.")
		}
	},
}

func init() {
	projectCmd.AddCommand(createCmd)
}
