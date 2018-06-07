package cmd

import (
	"github.com/dihedron/swift/swift"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "put bucket object [filename]",
		Short: "Store an object into a bucket",
		Long: `
Stores a file into a bucket, reading it from STDIN or from a local file named
[filename], if provided.`,
		Example: "  swift put my-bucket my-object filename.in",
		Aliases: []string{"retrieve", "download"},
		Args:    cobra.RangeArgs(2, 3),
		Run:     swift.PutObject,
	})
}
