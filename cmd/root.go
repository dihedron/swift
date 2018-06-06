package cmd

import (
	"fmt"
	"os"

	log "github.com/dihedron/go-log"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "swift",
	Short: "A minimalistic OpenStack Swift v1 client",
	Long: `This program provides a minimalistic OpenStack Swift v1 client with the 
ability to list all objects in a container, optionally filter the list, put a new 
file into an existing container, retrieve a file from a container, and delete it.`,
	Example:           "swift [command] [args...]",
	PersistentPreRun:  login,
	PersistentPostRun: logout,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing root command: %v", err)
	}
}

func init() {

	log.SetStream(os.Stderr, true)
	log.SetLevel(log.INF)

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

var storage *gophercloud.ServiceClient

func login(cmd *cobra.Command, args []string) {
	// obtain credentials from the environment
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		log.Fatalf("Unable to acquire credentials: %v", err)
	}

	// authenticate against keystone (v2 or v3)
	client, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Fatalf("Unable to authenticate: %v", err)
	}
	if client.TokenID == "" {
		log.Fatalf("No token ID assigned to the client")
	}

	log.Infof("Client successfully acquired a token: %v", client.TokenID)

	// find the storage service in the service catalog
	storage, err = openstack.NewObjectStorageV1(client, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		log.Fatalf("Unable to locate a storage service: %v", err)
	}

	log.Infof("Located a storage service at endpoint: [%s]", storage.Endpoint)
}

func logout(cmd *cobra.Command, args []string) {
	log.Infof("Logging out of storage service")
	if storage != nil {
		storage = nil
	}
}
