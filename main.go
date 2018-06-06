package main

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/dihedron/go-log"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
)

func usage() {
	fmt.Println(`usage: swift "my-container/my-object" [filename]`)
	os.Exit(1)
}

// swift get container/
func main() {
	var filename string

	switch {
	case len(os.Args) < 2:
		usage()
	case len(os.Args) > 2:
		filename = os.Args[2]
	}

	match := regexp.MustCompile(`^(?:/{0,1})([^/]+)(?:/{0,1})(.+)$`).FindStringSubmatch(os.Args[1])
	if len(match) != 3 {
		usage()
	}

	container := match[1]
	object := match[2]
	log.Infof("Retrieving %q from container %q", object, container)

	// obtain credentials from the environment
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		log.Fatalf("Unable to acquire credentials: %v", err)
	}

	client, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Fatalf("Unable to authenticate: %v", err)
	}

	if client.TokenID == "" {
		log.Fatalf("No token ID assigned to the client")
	}

	log.Infof("Client successfully acquired a token: %v", client.TokenID)

	// find the storage service in the service catalog
	storage, err := openstack.NewObjectStorageV1(client, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		log.Fatalf("Unable to locate a storage service: %v", err)
	}

	log.Infof("Located a storage service at endpoint: [%s]", storage.Endpoint)

	// download one of the objects that was created above
	response := objects.Download(storage, container, object, nil)
	if response.Err != nil {
		log.Fatalf("Unable to read object data: %v", err)
	}
	defer response.Body.Close()

	var output *os.File
	if filename != "" {
		output, err = os.Create(filename)
		if err != nil {
			log.Fatalf("Unable to open output file: %v", err)
		}
	} else {
		output = os.Stdout
	}

	count, err := io.Copy(output, response.Body)
	if err != nil {
		log.Fatalf("Unable to copy object data to file: %v", err)
	}
	log.Infof("Copied %d bytes to file", count)

	output.Close()
}
