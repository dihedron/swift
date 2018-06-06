package swift

import (
	"io"
	"os"

	log "github.com/dihedron/go-log"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/spf13/cobra"
)

var storage *gophercloud.ServiceClient

func Login(cmd *cobra.Command, args []string) {
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

func Logout(cmd *cobra.Command, args []string) {
	log.Infof("Logging out of storage service")
	if storage != nil {
		storage = nil
	}
}

func GetObject(cmd *cobra.Command, args []string) {
	container := args[0]
	object := args[1]

	// if a filename argument is provided, write to file, otherwise it's STDOUT
	var output *os.File
	if len(args) > 2 {
		if args[2] != "" {
			log.Debugf("Writing object to output file: %q", args[2])
			var err error
			output, err = os.Create(args[2])
			if err != nil {
				log.Fatalf("Unable to open output file: %v", err)
			}
		} else {
			output = os.Stdout
		}
	}

	// download one of the objects that was created above
	response := objects.Download(storage, container, object, nil)
	if response.Err != nil {
		log.Fatalf("Unable to read object data: %v", response.Err)
	}
	defer response.Body.Close()

	count, err := io.Copy(output, response.Body)
	if err != nil {
		log.Fatalf("Unable to copy object data to file: %v", err)
	}
	log.Infof("Copied %d bytes to file", count)

	output.Close()
}
