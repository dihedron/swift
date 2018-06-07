package cmd

import (
	"github.com/dihedron/swift/swift"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "list bucket [filter]",
		Short: "List objects in a bucket",
		Long: `
Lists all objects in a bucket; it can optionally apply a user-provided regular
expression to filter objects.`,
		Example: "  swift list my-bucket '^my-obj.*'",
		Aliases: []string{"search", "find", "l", "s", "f"},
		Args:    cobra.RangeArgs(1, 2),
		Run:     swift.ListObjects,
	})
}
