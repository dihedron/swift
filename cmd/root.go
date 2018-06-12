package cmd

import (
	"fmt"
	"os"

	log "github.com/dihedron/go-log"
	"github.com/dihedron/swift/swift"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "swift",
	Short: "A minimalistic OpenStack Swift v1 client",
	Long: `
This program provides a minimalistic OpenStack Swift v1 client with the ability
to list all objects in a bucket, optionally filter the list, put a new file into 
an existing bucket (upload), retrieve a file from a bucket (download), and delete 
it.`,
	Example:           "swift [command] [args...]",
	PersistentPreRun:  swift.Login,
	PersistentPostRun: swift.Logout,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing root command: %v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "configuration file (default is $HOME/.swift.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".swift" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".swift")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
