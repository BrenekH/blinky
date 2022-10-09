package viperutils

import (
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Setup() error {
	SetupDefaults()

	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	if configDir, err := os.UserConfigDir(); err == nil {
		viper.AddConfigPath(configDir + "/blinky") // Prioritize user config directory
	}

	viper.AddConfigPath("/etc/blinky")

	SetupEnvVars()
	SetupFlags()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return nil
		} else {
			// Config file was found but another error was produced
			return err
		}
	}

	return nil
}

func SetupDefaults() error {
	// RepoPath
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	viper.SetDefault("RepoPath", cwd+"/repo")

	// RequireSignedPkgs
	viper.SetDefault("RequireSignedPkgs", true)

	// SigningKeyFile
	viper.SetDefault("SigningKeyFile", "")

	// DataDir
	viper.SetDefault("DataDir", "/var/lib/blinky")

	// GPGDir
	viper.SetDefault("GPGDir", "/tmp/blinky/gnupg")

	// HTTPPort
	viper.SetDefault("HTTPPort", "9000")

	// APIUsername
	viper.SetDefault("APIUsername", "")

	// APIPassword
	viper.SetDefault("APIPassword", "")

	return nil
}

func SetupEnvVars() {
	viper.SetEnvPrefix("BLINKY")

	viper.BindEnv("RepoPath", "BLINKY_REPO_PATH")
	viper.BindEnv("RequireSignedPkgs", "BLINKY_SIGNED_PKGS")
	viper.BindEnv("SigningKeyFile", "BLINKY_SIGNING_KEY")
	viper.BindEnv("DataDir", "BLINKY_DATA_DIR")
	viper.BindEnv("GPGDir", "BLINKY_GPG_DIR")
	viper.BindEnv("HTTPPort", "BLINKY_PORT")
	viper.BindEnv("APIUsername", "BLINKY_API_UNAME")
	viper.BindEnv("APIPassword", "BLINKY_API_PASSWD")
}

func SetupFlags() {
	pflag.StringP("repo-path", "r", "", "Colon-separated paths to use as repositories")

	pflag.Bool("no-signed-pkgs", true, "Do not require that packages be uploaded with a signature")
	pflag.Lookup("no-signed-pkgs").NoOptDefVal = "false"

	pflag.String("signing-key", "", "Filepath of a GPG key to use to sign the Pacman database")

	pflag.String("data-dir", "", "Directory to store Blinky's runtime files")

	pflag.String("gpg-dir", "", "Specify a custom location to construct a GPG keyring")

	pflag.StringP("http-port", "p", "", "Select the port to host Blinky on")

	pflag.String("api-uname", "", "The username to use to protect the API")

	pflag.String("api-passwd", "", "The password to use to protect the API")

	pflag.Parse()

	viper.BindPFlag("RepoPath", pflag.Lookup("repo-path"))
	viper.BindPFlag("RequireSignedPkgs", pflag.Lookup("no-signed-pkgs"))
	viper.BindPFlag("SigningKeyFile", pflag.Lookup("signing-key"))
	viper.BindPFlag("DataDir", pflag.Lookup("data-dir"))
	viper.BindPFlag("GPGDir", pflag.Lookup("gpg-dir"))
	viper.BindPFlag("HTTPPort", pflag.Lookup("http-port"))
	viper.BindPFlag("APIUsername", pflag.Lookup("api-uname"))
	viper.BindPFlag("APIPassword", pflag.Lookup("api-passwd"))
}
