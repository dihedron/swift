package swift

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	log "github.com/dihedron/go-log"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/spf13/cobra"
)

var client *gophercloud.ProviderClient
var storage *gophercloud.ServiceClient

// Login performs a login to OpenStack Keystone using parameters in the current
// process environment; in order to set the relevant environment variables, use
// the OpenStack-generated openstackrc file.
func Login(cmd *cobra.Command, args []string) {
	// obtain credentials from the environment
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		log.Fatalf("Unable to acquire credentials: %v", err)
	}

	// authenticate against keystone (v2 or v3)
	client, err = openstack.AuthenticatedClient(opts)
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

// Logout inhibits further requests to the object store.
func Logout(cmd *cobra.Command, args []string) {
	log.Infof("Logging out of storage service")
	if storage != nil {
		storage = nil
	}
}

// GetObject downloads an object from the object store, writing it either to
// STDOUT or, if a file name is provided, to a local file.
func GetObject(cmd *cobra.Command, args []string) {
	bucket := args[0]
	object := args[1]

	log.Debugf("Downloading %q from %q", object, bucket)

	// if a filename argument is provided, write to file, otherwise it's STDOUT
	var output *os.File
	if len(args) > 2 && args[2] != "" {
		log.Debugf("Writing object to output file: %q", args[2])
		var err error
		output, err = os.Create(args[2])
		if err != nil {
			log.Fatalf("Unable to open output file: %v", err)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	// download the object from the remote storage
	response := objects.Download(storage, bucket, object, nil)
	if response.Err != nil {
		log.Fatalf("Unable to read object data: %v", response.Err)
	}
	defer response.Body.Close()

	log.Debugf("Received response")

	count, err := io.Copy(output, response.Body)
	if err != nil {
		log.Fatalf("Unable to copy object data to file: %v", err)
	}
	log.Debugf("Copied %d bytes to file", count)
}

// PutObject stores an object in the object store under the given name; the file
// data may be read from STDIN or from a local file, if a valid name is provided.
func PutObject(cmd *cobra.Command, args []string) {
	bucket := args[0]
	object := args[1]

	// if a filename argument is provided, read from file, otherwise it's STDIN
	var input *os.File
	if len(args) > 2 && args[2] != "" {
		log.Debugf("Reading object from input file: %q", args[2])
		var err error
		input, err = os.Open(args[2])
		if err != nil {
			log.Fatalf("Unable to open input file: %v", err)
		}
		defer input.Close()
	} else {
		input = os.Stdin
	}

	// upload one of the objects that was created above
	options := &objects.CreateOpts{
		ContentType: "application/octet-stream", // as per RFC2046
		Content:     input,
	}
	result := objects.Create(storage, bucket, object, options)
	if result.Err != nil {
		log.Fatalf("Unable to store object data: %v", result.Err)
	}

	log.Debugf("Uploaded file to store")
}

// DeleteObject removes an object from the object store given its name.
func DeleteObject(cmd *cobra.Command, args []string) {
	bucket := args[0]
	object := args[1]

	result := objects.Delete(storage, bucket, object, nil)
	if result.Err != nil {
		log.Fatalf("Unable to delete object: %v", result.Err)
	}

	log.Debugf("Removed object from store")
}

// AboutObject retrieves metadata about an object.
func AboutObject(cmd *cobra.Command, args []string) {
	bucket := args[0]
	object := args[1]

	log.Debugf("Getting info about %q from %q", object, bucket)

	// retrieve the object metadata from the remote storage
	response := objects.Get(storage, bucket, object, nil)
	if response.Err != nil {
		log.Fatalf("Unable to read object metadadata: %v", response.Err)
	}
	metadata, err := response.ExtractMetadata()
	if err != nil {
		log.Fatalf("Unable to read object metadata: %v", err)
	}

	for k, v := range metadata {
		fmt.Printf("%q: %q\n", k, v)
	}
}

// ListObjects retrieves the list of objects in a bucket; if a filter is provided,
// only objects stores an object in the object store under the given name; the file
// data may be read from STDIN or from a local file, if a valid name is provided.
func ListObjects(cmd *cobra.Command, args []string) {
	bucket := args[0]
	filter := ".*"
	if len(args) > 1 {
		filter = args[1]
	}

	strings.ToLower(filter)

	pager := objects.List(storage, bucket, nil)
	if pager.Err != nil {
		log.Fatalf("Unable to get list of objects: %v", pager.Err)
	}

	list := []string{}

	pager.EachPage(func(page pagination.Page) (bool, error) {
		names, err := objects.ExtractNames(page)
		if err != nil {
			log.Fatalf("Unable to extract names: %v", err)
		}
		list = append(list, names...)
		return true, nil
	})

	sort.Strings(list)

	var re *regexp.Regexp
	if len(args) > 1 && len(args[1]) > 0 {
		log.Debugf("Filtering with %q", args[1])
		re = regexp.MustCompile(args[1])
	}
	for _, item := range list {
		if re == nil || re.MatchString(item) {
			fmt.Printf("%s\n", item)
		}
	}

	log.Debugf("Retrieved list of objects in bucket")
}
