package cmd

import (
	"github.com/dihedron/swift/swift"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "about bucket object",
		Short: "Retrieve information about an object in a bucket",
		Long: `
Retrieves information about an object.`,
		Example: "  swift about my-bucket my-object",
		Aliases: []string{"info"},
		Args:    cobra.ExactArgs(2),
		Run:     swift.AboutObject,
	})
}
