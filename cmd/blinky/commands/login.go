package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Check and save the login information for a Blinky server.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		setAsDefault := viper.GetBool("default")
		flagUsername := viper.Get("username")
		flagPassword := viper.Get("password")

		fmt.Println(setAsDefault, flagUsername, flagPassword)

		// NOTE: If auth info is already in DB, ask if user would like to update it.
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
