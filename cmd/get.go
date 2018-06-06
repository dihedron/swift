package cmd

import (
	"io"
	"os"

	log "github.com/dihedron/go-log"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "get container object [filename]",
		Short: "Retrieve an object from a container",
		Long: `
Retrieves a file from a container, optonally storing it as [filename]; if no
filename is provided, the file contents are written to STDOUT.`,
		Example: "  swift get my-container my-object filename.bin",
		Aliases: []string{"retrieve", "download"},
		Args: func(cmd *cobra.Command, args []string) error {
			err := (cobra.RangeArgs(2, 3))(cmd, args)
			if err != nil {
				return err
			}
			// add extra validation here
			return nil
		},
		Run: getObject,
	})
}

func getObject(cmd *cobra.Command, args []string) {
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
