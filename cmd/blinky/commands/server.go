package commands

import (
	"fmt"
	"os"

	"github.com/BrenekH/blinky/cmd/blinky/util"
	"github.com/spf13/cobra"
)

func printServerInfo(n string, s util.Server, showPassword bool) {
	var serverInfo string
	if s.Username == "" {
		serverInfo = fmt.Sprintf("%s\n", n)
	} else if showPassword {
		serverInfo = fmt.Sprintf("%s(%s:%s)\n", n, s.Username, s.Password)
	} else {
		serverInfo = fmt.Sprintf("%s(%s)\n", n, s.Username)
	}
	fmt.Print(serverInfo)
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

		showPassword, err := cmd.Flags().GetBool("password")
		if err != nil {
			fmt.Printf("Unexpected error while reading flag: %v", err)
			os.Exit(1)
		}

		showOnlyDefault, err := cmd.Flags().GetBool("default")
		if err != nil {
			fmt.Printf("Unexpected error while reading flag: %v", err)
			os.Exit(1)
		}

		if len(args) > 0 {
			for _, server := range args {
				if _, ok := serverDB.Servers[server]; !ok {
					fmt.Printf("Server %s not found.\n", server)
					os.Exit(1)
				} else {
					printServerInfo(server, serverDB.Servers[server], showPassword)
				}
			}
		}

		if showOnlyDefault {
			printServerInfo(serverDB.DefaultServer, serverDB.Servers[serverDB.DefaultServer], showPassword)
		}

		if !(len(args) > 0) && !showOnlyDefault {
			for name, server := range serverDB.Servers {
				printServerInfo(name, server, showOnlyDefault)
			}
		}

	},
}

func init() {
	serverCmd.Flags().Bool("default", false, "Show only the default server")
	serverCmd.Flags().Bool("password", false, "Show the password for the server")

	rootCmd.AddCommand(serverCmd)

}
