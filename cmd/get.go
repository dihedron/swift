package cmd

import (
	"github.com/dihedron/swift/swift"
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
		Run: swift.GetObject,
	})
}
