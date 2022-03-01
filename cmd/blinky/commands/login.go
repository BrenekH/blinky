package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/BrenekH/blinky/cmd/blinky/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login server_url",
	Short: "Check and save the login information for a Blinky server.",
	Long: `Login to a blinky server and optionally set it as the default server
used by upload and remove.

A user may choose to pass login details using
the --username and --password flags. Any fields that are provided to login
will be prompted for.`,
	Run: func(cmd *cobra.Command, args []string) {
		setAsDefault := viper.GetBool("default")
		username := viper.GetString("username")
		password := viper.GetString("password")

		if len(args) != 1 {
			fmt.Printf("Incorrect number of arguments for login command. Expected 1, got %v.\n\nUse blinky login --help for more information.\n", len(args))
			os.Exit(1)
		}
		serverURL := args[0]

		serverDB, err := util.ReadServerDB()
		if err != nil {
			fmt.Printf("Unexpected error while reading servers.json: %v\n", err)
			os.Exit(1)
		}

		serverEntry, ok := serverDB.Servers[serverURL]
		if ok {
			fmt.Printf("Login information is already available for %s\n", serverURL)
			userIn := util.Input("Would you like to override (y/N)? ")
			if strings.ToLower(userIn) == "y" {
				if username == "" {
					username = util.Input("Username: ")
				}
				if password == "" {
					password = util.SecureInput("Password: ")
				}

				serverEntry.Username = username
				serverEntry.Password = password
				serverDB.Servers[serverURL] = serverEntry
			}
		} else {
			if username == "" {
				username = util.Input("Username: ")
			}
			if password == "" {
				password = util.SecureInput("Password: ")
			}

			serverEntry.Username = username
			serverEntry.Password = password
			serverDB.Servers[serverURL] = serverEntry
		}

		if setAsDefault {
			serverDB.DefaultServer = serverURL
		}

		if err := util.SaveServerDB(serverDB); err != nil {
			fmt.Printf("Unexpected error while saving servers.json: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().BoolP("default", "d", false, "Set server as default for upload and remove")
	loginCmd.Flags().String("username", "", "Set username for server login")
	loginCmd.Flags().String("password", "", "Set password for server login. It is recommended to enter the password in the interactive prompt instead of using this flag")

	viper.BindPFlag("default", loginCmd.Flags().Lookup("default"))
	viper.BindPFlag("username", loginCmd.Flags().Lookup("username"))
	viper.BindPFlag("password", loginCmd.Flags().Lookup("password"))
}
