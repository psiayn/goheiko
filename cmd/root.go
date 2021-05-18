package cmd

import (
	"fmt"
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

	// config file location
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $PWD/.heiko/config.yaml)")

	// the <name>
	rootCmd.PersistentFlags().StringP("name", "n", "", "Unique name to give (or given) to this heiko job")
	rootCmd.MarkPersistentFlagRequired("name")
	viper.BindPFlag("name", rootCmd.PersistentFlags().Lookup("name"))

	// flags for daemonizing
	startCmd.PersistentFlags().BoolP("daemon", "d", false, "Daemonizing heiko")
	viper.BindPFlag("daemon", startCmd.PersistentFlags().Lookup("daemon"))

	// add sub-commands here
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(initCmd)
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".heiko")
		viper.SetConfigName("config")
	}

	// this makes viper look for environment
	//  variables of the form HEIKO_SOMETHING
	viper.SetEnvPrefix("heiko")
	viper.AutomaticEnv()

	// data storage location - for logs, PID file, etc.
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)
	dataLocation := filepath.Join(homeDir, ".heiko")
	viper.SetDefault("dataLocation", dataLocation)

	// read config from ~/.heiko/config.{whatever-extension}
	err = viper.ReadInConfig()
	cobra.CheckErr(err)

	// unmarshalling is basically loading the data into our struct
	// look at internal/config/config.go for the structure
	err = viper.Unmarshal(&configuration)
	cobra.CheckErr(err)

	// set defaults using the "defaults" package
	// look for `default:"something"` in
	//    internal/config/config.go for the defaults
	err = defaults.Set(&configuration)
	cobra.CheckErr(err)

	// we need to ensure that all commands run on the nodes
	//   are run inside ~/.heiko/<name> directory
	//   so that the home directory isn't polluted because of heiko
	for i, job := range configuration.Jobs {
		// when initializing, we create this directory (for each job)
		//    and cd into it
		cdCommand := fmt.Sprintf("cd ~/.heiko/%s", job.Name)
		configuration.Jobs[i].Init = append([]string{
			fmt.Sprintf("mkdir -p ~/.heiko/%s", job.Name),
			cdCommand,
		}, job.Init...)
		// ^ this weird syntax above (and below) is for inserting elements
		//     to the beginning of the slice

		// when running, we only cd into it
		configuration.Jobs[i].Commands = append([]string{
			cdCommand,
		}, job.Commands...)
	}

	// ensure ~/.heiko/<name> exists
	os.MkdirAll(filepath.Join(dataLocation, viper.GetString("name")),
		os.ModePerm)
}
