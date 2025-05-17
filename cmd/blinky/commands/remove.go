package commands

import (
	"fmt"
	"os"

	"github.com/BrenekH/blinky/clientlib"
	"github.com/BrenekH/blinky/cmd/blinky/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove repo_name packages...",
	Short: "Remove packages from a Blinky server",
	Long: `Remove multiple packages from a pacman repository hosted
on a Blinky server. If the --server flag is not provided,
the default server will be used.

The user may override the saved username and password with the
--username and --password flags. remove will also prompt for a
password if --password is not used and --ask-pass is passed.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("incorrect number of arguments for remove command. Expected >=2, got %v", len(args))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server := viper.GetString("server")
		username := viper.GetString("username")
		password := viper.GetString("password")
		promptForPasswd := viper.GetBool("ask-pass")

		serverDB, err := util.ReadServerDB()
		if err != nil {
			fmt.Printf("Unexpected error while reading servers.json: %v\n", err)
			os.Exit(1)
		}

		if server == "" {
			server = serverDB.DefaultServer
		}

		if username == "" {
			username = serverDB.Servers[server].Username
		}

		if password == "" {
			if promptForPasswd {
				password = util.SecureInput("Password: ")
			} else {
				password = serverDB.Servers[server].Password
			}
		}

		client, err := clientlib.New(server, username, password)
		if err != nil {
			fmt.Printf("Error while creating client: %v", err)
			os.Exit(1)
		}

		repoName := args[0]
		packagesToRemove := args[1:]

		if err := client.RemovePackages(repoName, packagesToRemove...); err != nil {
			fmt.Printf("Error while removing packages: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringP("server", "s", "", "Server URL to remove package from")
	removeCmd.Flags().String("username", "", "Set username for server login")
	removeCmd.Flags().String("password", "", "Set password for server login. It is recommended to enter the password in the interactive prompt instead of using this flag")
	removeCmd.Flags().BoolP("ask-pass", "K", false, "Ask for the password in an interactive prompt")

	viper.BindPFlag("server", removeCmd.Flags().Lookup("server"))
	viper.BindPFlag("username", removeCmd.Flags().Lookup("username"))
	viper.BindPFlag("password", removeCmd.Flags().Lookup("password"))
	viper.BindPFlag("ask-pass", removeCmd.Flags().Lookup("ask-pass"))
}
