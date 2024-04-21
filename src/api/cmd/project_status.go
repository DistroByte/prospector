package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a job",
	Long: `Get the status of a job. Returns the status of a job with the given name.
	
For example:
	prospect job status my-job`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		addr := cmd.Flag("address").Value.String()
		res, err := CmdGet(addr + "/api/v1/jobs/" + name)
		if err != nil {
			fmt.Println("Error: The server responded with an error. Please try again.")
			return
		}

		if res.StatusCode == 204 {
			fmt.Println("No job found")
			return
		}

		// res.Body is a json string
		io.Copy(os.Stdout, res.Body)
	},
}

func init() {
	projectCmd.AddCommand(statusCmd)

	statusCmd.Args = cobra.ExactArgs(1)
	statusCmd.Flags().StringP("address", "a", "https://prospector.ie", "The address of the Prospector server")
}
