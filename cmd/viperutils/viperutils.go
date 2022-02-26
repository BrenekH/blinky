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

	return viper.ReadInConfig()
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

	// ConfigDir
	viper.SetDefault("ConfigDir", "/var/lib/blinky")

	// GPGDir
	viper.SetDefault("GPGDir", "/tmp/blinky/gnupg")

	// HTTPPort
	viper.SetDefault("HTTPPort", "9000")

	return nil
}

func SetupEnvVars() {
	viper.SetEnvPrefix("BLINKY")

	viper.BindEnv("RepoPath", "BLINKY_REPO_PATH")
	viper.BindEnv("RequireSignedPkgs", "BLINKY_SIGNED_PKGS")
	viper.BindEnv("SigningKeyFile", "BLINKY_SIGNING_KEY")
	viper.BindEnv("ConfigDir", "BLINKY_CONFIG_DIR")
	viper.BindEnv("GPGDir", "BLINKY_GPG_DIR")
	viper.BindEnv("HTTPPort", "BLINKY_PORT")
}

func SetupFlags() {
	pflag.StringP("repo-path", "r", "", "--repo-path, -r <paths seperated with colons>")

	pflag.Bool("no-signed-pkgs", true, "--no-signed-pkgs")
	pflag.Lookup("no-signed-pkgs").NoOptDefVal = "false"

	pflag.String("signing-key", "", "--signing-key <filepath>")

	pflag.String("config-dir", "", "--config-dir <dir>")

	pflag.String("gpg-dir", "", "--gpg-dir <dir>")

	pflag.StringP("http-port", "p", "", "--http-port, -p <port number>")

	pflag.Parse()

	viper.BindPFlag("RepoPath", pflag.Lookup("repo-path"))
	viper.BindPFlag("RequireSignedPkgs", pflag.Lookup("no-signed-pkgs"))
	viper.BindPFlag("SigningKeyFile", pflag.Lookup("signing-key"))
	viper.BindPFlag("ConfigDir", pflag.Lookup("config-dir"))
	viper.BindPFlag("GPGDir", pflag.Lookup("gpg-dir"))
	viper.BindPFlag("HTTPPort", pflag.Lookup("http-port"))
}