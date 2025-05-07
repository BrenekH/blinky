package commands

import (
	"fmt"
	"os"

	"github.com/BrenekH/blinky/cmd/blinky/util"
	"github.com/spf13/cobra"
)

func printServerInfo(n string, s util.Server) {
	if s.Username == "" {
		fmt.Printf("%s\n", n)
	} else {
		fmt.Printf("%s(%s)\n", n, s.Username)
	}
}

var serverCmd = &cobra.Command{
	Use:   "server [server_url...]",
	Short: "Manage Blinky servers",
	Long: `Manage Blinky servers. This command is used to list the servers
that are saved in the server database.
It can also be used to check the login information for a specific server.`,
	Args:   cobra.ArbitraryArgs,
	PreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		serverDB, err := util.ReadServerDB()
		if err != nil {
			fmt.Printf("Unexpected error while reading servers.json: %v", err)
			os.Exit(1)
		}

		if len(args) > 0 {
			for _, server := range args {
				if _, ok := serverDB.Servers[server]; !ok {
					fmt.Printf("Server %s not found.\n", server)
					os.Exit(1)
				} else {
					printServerInfo(server, serverDB.Servers[server])
				}
			}

		} else {
			for name, server := range serverDB.Servers {
				printServerInfo(name, server)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
