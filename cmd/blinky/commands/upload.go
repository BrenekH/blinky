package commands

import (
	"fmt"
	"os"

	"github.com/BrenekH/blinky/clientlib"
	"github.com/BrenekH/blinky/cmd/blinky/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload repo_name package_files...",
	Short: "Upload packages to a Blinky server",
	Long: `Upload multiple packages to a pacman repository hosted
on a Blinky server. If the --server flag is not provided,
the default server will be used.

The user may override the saved username and password with the
--username and --password flags. upload will also prompt for a
password if --password is not used and --ask-pass is passed.

If a matching ".sig" is found alongside the package file, it will
be uploaded along with the package to the target server.`,
	Run: func(cmd *cobra.Command, args []string) {
		server := viper.GetString("server")
		username := viper.GetString("username")
		password := viper.GetString("password")
		promptForPasswd := viper.GetBool("ask-pass")

		if len(args) < 2 {
			fmt.Printf("Incorrect number of arguments for upload command. Expected >=2, got %v.\n\nUse blinky upload --help for more information.\n", len(args))
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

		client := clientlib.New(server, username, password)

		repoName := args[0]
		packageFiles := args[1:]

		if err := client.UploadPackageFiles(repoName, packageFiles...); err != nil {
			fmt.Printf("Error while removing packages: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringP("server", "s", "", "Server URL to remove package from")
	uploadCmd.Flags().String("username", "", "Set username for server login")
	uploadCmd.Flags().String("password", "", "Set password for server login. It is recommended to enter the password in the interactive prompt instead of using this flag")
	uploadCmd.Flags().BoolP("ask-pass", "K", false, "Ask for the password in an interactive prompt")

	viper.BindPFlag("server", uploadCmd.Flags().Lookup("server"))
	viper.BindPFlag("username", uploadCmd.Flags().Lookup("username"))
	viper.BindPFlag("password", uploadCmd.Flags().Lookup("password"))
	viper.BindPFlag("ask-pass", uploadCmd.Flags().Lookup("ask-pass"))
}
