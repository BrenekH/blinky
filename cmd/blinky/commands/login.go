package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Check and save the login information for a Blinky server.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")

		// NOTE: If auth info is already in DB, ask if user would like to update it.
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
