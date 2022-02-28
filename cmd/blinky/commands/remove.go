package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/BrenekH/blinky/cmd/blinky/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove repo_name packages...",
	Short: "Remove packages from a Blinky server.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		server := viper.GetString("server")
		// username := viper.GetString("username")
		password := viper.GetString("password")
		promptForPasswd := viper.GetBool("ask-pass")

		if len(args) < 2 {
			fmt.Printf("Incorrect number of arguments for remove command. Expected >=2, got %v.\n\nUse blinky remove --help for more information.\n", len(args))
			os.Exit(1)
		}

		serverDB, err := util.ReadServerDB()
		if err != nil {
			fmt.Printf("Unexpected error while reading servers.json: %v\n", err)
			os.Exit(1)
		}

		if server == "" {
			server = serverDB.DefaultServer
		}

		// if username == "" {
		// 	username = serverDB.Servers[server].Username
		// }

		if password == "" {
			if promptForPasswd {
				password = util.SecureInput("Password: ")
			} else {
				password = serverDB.Servers[server].Password
			}
		}

		repoName := args[0]
		packagesToRemove := args[1:]

		// Will have to loop through each package and call the API one by one.
		// Ideally the API would support multiple package removals at once.

		for _, pkg := range packagesToRemove {
			fmt.Printf("Removing %s\n", pkg)

			r, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/unstable/%s/package/%s", server, repoName, pkg), bytes.NewBufferString(""))
			if err != nil {
				panic(err) // TODO: Handle this better
			}

			r.Header.Add("Authorization", password)

			resp, err := http.DefaultClient.Do(r)
			if err != nil {
				panic(err) // TODO: Handle better
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				b, _ := io.ReadAll(resp.Body)

				fmt.Printf("Received a non-200 status code while removing %s/%s: %s - %s", repoName, pkg, resp.Status, string(b))
				os.Exit(1)
			}

			fmt.Printf("%s removed.\n", pkg)
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
