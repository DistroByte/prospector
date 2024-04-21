package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with the server",
	Long: `Authenticate with the server. This command will authenticate with the server and return a token.
	
For example:
	prospector auth --username my-username --password my-password`,

	Run: func(cmd *cobra.Command, args []string) {
		addr := cmd.Flag("address").Value.String()

		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		token, err := ProspectorAuth(addr, username, password)
		if err != nil {
			fmt.Println("Error: Could not authenticate with the server. Are your credentials correct?")
			return
		}

		// store token in ~/.prospector_token for future use
		storeToken(token)
		fmt.Println("Authenticated successfully")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringP("username", "u", "", "Username")
	authCmd.Flags().StringP("password", "p", "", "Password")
	authCmd.Flags().Lookup("username").NoOptDefVal = os.Getenv("PROSPECTOR_username")
	authCmd.Flags().Lookup("password").NoOptDefVal = os.Getenv("PROSPECTOR_password")
}

func storeToken(token string) {
	// create a file in ~/.prospector_token and store the token
	err := os.WriteFile(os.Getenv("HOME")+"/.prospector_token", []byte(token), 0644)
	if err != nil {
		fmt.Println("Error: Could not store token. Please try again.")
	}
}
