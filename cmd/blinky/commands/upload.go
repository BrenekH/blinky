package commands

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BrenekH/blinky/cmd/blinky/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload repo_name package_file...",
	Short: "Upload packages to a Blinky server.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		server := viper.GetString("server")
		// username := viper.GetString("username")
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
		packageFiles := args[1:]

		for _, pkg := range packageFiles {
			fmt.Printf("Uploading %s\n", pkg)

			r, w := io.Pipe()
			writer := multipart.NewWriter(w)

			go func() {
				defer w.Close()
				defer writer.Close()

				file, err := os.Open(pkg)
				if err != nil {
					panic(err)
				}
				defer file.Close()

				part, err := writer.CreateFormFile("package", filepath.Base(file.Name()))
				if err != nil {
					panic(err)
				}

				_, err = io.Copy(part, file)
				if err != nil {
					panic(err)
				}

				if sigFile, err := os.Open(pkg + ".sig"); err == nil {
					part, err := writer.CreateFormFile("signature", filepath.Base(sigFile.Name()))
					if err != nil {
						panic(err)
					}

					_, err = io.Copy(part, sigFile)
					if err != nil {
						panic(err)
					}
				}
			}()

			// TODO: Decompress passed package file to identify the package name to be sent through the API
			request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/unstable/%s/package/%s", server, repoName, "replace_me"), r)
			if err != nil {
				panic(err) // TODO: Handle this better
			}

			request.Header.Add("Authorization", password)
			request.Header.Add("Content-Type", writer.FormDataContentType())

			resp, err := http.DefaultClient.Do(request)
			if err != nil {
				panic(err) // TODO: Handle better
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				b, _ := io.ReadAll(resp.Body)

				fmt.Printf("Received a non-200 status code while uploading %s/%s: %s - %s", repoName, pkg, resp.Status, string(b))
				os.Exit(1)
			}

			fmt.Printf("%s uploaded.\n", pkg)
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
