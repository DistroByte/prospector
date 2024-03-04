package cmd

import (
	"prospector/middleware"
	"prospector/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Prospector API server",
	Long:  `Starts the Prospector API server. The server will be bound to the port specified by the --port flag.`,

	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()

		identityKey := cmd.Flag("identity-key").Value.String()

		middleware.CreateStandardMiddlewares(r)
		middleware.CreateAuthMiddlewares(r, identityKey)
		routes.Route(r, identityKey)

		r.Run(":" + cmd.Flag("port").Value.String())
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringP("port", "p", "3434", "Port to bind the server to")
	serverCmd.Flags().StringP("identity-key", "i", "id", "Identity key for the JWT middleware")
}
