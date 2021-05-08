package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/psiayn/heiko/internal/config"
	"github.com/creasty/defaults"
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
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// Do Stuff Here
	// },
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

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Error in reading configuration file:", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		fmt.Fprintln(os.Stderr, "Error in reading configuration file:", err)
		os.Exit(1)
	}

	if err := defaults.Set(&configuration); err != nil {
		panic(fmt.Errorf("panik: Could not set defaults %v", err))
	}
}
