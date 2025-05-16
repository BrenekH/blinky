package commands

import (
	"os"

	"github.com/BrenekH/blinky/vars"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "blinky",
	Short: "Manage packages in a Blinky repository system",
	Long: `blinky is a CLI used for uploading and removing packages
from a Blinky Pacman repository.`,
	Version: vars.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
