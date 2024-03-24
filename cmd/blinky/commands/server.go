package commands

import (
	"fmt"
	"os"

	"github.com/BrenekH/blinky/cmd/blinky/util"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage Blinky servers",
	Run: func(cmd *cobra.Command, args []string) {
		serverDB, err := util.ReadServerDB()
		if err != nil {
			fmt.Printf("Unexpected error while reading servers.json: %v", err)
			os.Exit(1)
		}

		for name, server := range serverDB.Servers {
			fmt.Println(name + "(" + server.Username + ")")
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
