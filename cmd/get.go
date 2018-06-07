package cmd

import (
	"github.com/dihedron/swift/swift"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "get bucket object [filename]",
		Short: "Retrieve an object from a bucket",
		Long: `
Retrieves a file from a bucket, writing it to STDOUT or, if [filename] is provided,
writing it to a local file by the same name.`,
		Example: "  swift get my-bucket my-object filename.out",
		Aliases: []string{"retrieve", "download"},
		Args:    cobra.RangeArgs(2, 3),
		Run:     swift.GetObject,
	})
}
