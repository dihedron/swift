package cmd

import (
	"github.com/dihedron/swift/swift"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "remove bucket object",
		Short: "Remove an object from a bucket",
		Long: `
Removes an object from a bucket by name.`,
		Example: "  swift remove my-bucket my-object",
		Aliases: []string{"drop", "delete", "rem", "del"},
		Args:    cobra.ExactArgs(2),
		Run:     swift.DeleteObject,
	})
}
