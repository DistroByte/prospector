/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all running jobs on the server",
	Long: `List all running jobs on the server. This command will make a request to the server and print out all the running jobs.
	
For example:
    prospector job list`,
	Run: func(cmd *cobra.Command, args []string) {
		healthURL := "http://localhost:3434/api/health"

		// make the request to the server
		res, err := http.Get(healthURL)
		if err != nil {
			fmt.Println("Error: Could not connect to prospector server, is it running?")
			return
		}

		// check the status code of the response
		if res.StatusCode != 200 {
			fmt.Println("Error: There was an error connecting to the server. Please try again.")
			return
		}

		jobListUrl := "http://localhost:3434/api/jobs"

		// make the request to the server
		res, err = http.Get(jobListUrl)
		if err != nil {
			fmt.Println("Error: The server responded with an error. Please try again.")
			return
		}

		// check the status code of the response
		if res.StatusCode == 204 {
			fmt.Println("No jobs found.")
			return
		}

		// check the status code of the response
		if res.StatusCode != 200 {
			fmt.Println("Error: The server responded with an error. Please try again.")
			return
		}

		// read the response body
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error: There was an error reading the response body. Please try again.")
			return
		}

		// print the response body
		fmt.Println(string(body))
	},
}

func init() {
	jobCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
