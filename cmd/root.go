package cmd

import (
	"os"
	"path/filepath"

	"github.com/creasty/defaults"
	"github.com/psiayn/heiko/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// stores path of config file
var configFile string

// stores entire configuration
var configuration config.Config

var rootCmd = &cobra.Command{
	Use:   "heiko",
	Short: "heiko is a lightweight not-a-load-balancer-but-something-like-that",
	Long: `heiko uses SSH to manage servers running on low power hardware such as
                Raspberry Pis or mobile phones.
                Made and maintained by PES Open Source.
                More details available at https://github.com/pesos/heiko`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $PWD/.heiko/config.yaml)")
	rootCmd.PersistentFlags().StringP("name", "n", "", "Unique name to give (or given) to this heiko job")
	rootCmd.MarkPersistentFlagRequired("name")

	viper.BindPFlag("name", rootCmd.PersistentFlags().Lookup("name"))

	rootCmd.AddCommand(startCmd)
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".heiko")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("heiko")
	viper.AutomaticEnv()

	// data storage location - for logs, PID file, etc.
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)
	dataLocation := filepath.Join(homeDir, ".heiko")
	viper.SetDefault("dataLocation", dataLocation)

	err = viper.ReadInConfig()
	cobra.CheckErr(err)

	err = viper.Unmarshal(&configuration)
	cobra.CheckErr(err)

	err = defaults.Set(&configuration)
	cobra.CheckErr(err)

	// ensure ~/.heiko/<name> exists
	os.MkdirAll(filepath.Join(dataLocation, viper.GetString("name")),
		os.ModePerm)
}
