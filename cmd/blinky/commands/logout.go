package commands

import (
	"fmt"
	"os"

	"github.com/BrenekH/blinky/cmd/blinky/util"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout server_url...",
	Short: "Delete the credentials for multiple servers",
	Long: `Delete the credentials for multiple servers.

logout will clear the default server if it is set
to one of the servers being removed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Printf("Incorrect number of arguments for logout command. Expected >=1, got %v.\n\nUse blinky logout --help for more information.\n", len(args))
			os.Exit(1)
		}
		serverURLs := args

		serverDB, err := util.ReadServerDB()
		if err != nil {
			fmt.Printf("Unexpected error while reading servers.json: %v\n", err)
			os.Exit(1)
		}

		for _, v := range serverURLs {
			// Reset the default server if the credentials are being removed.
			if serverDB.DefaultServer == v {
				serverDB.DefaultServer = ""
			}

			delete(serverDB.Servers, v)
		}

		if err := util.SaveServerDB(serverDB); err != nil {
			fmt.Printf("Unexpected error while saving servers.json: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
