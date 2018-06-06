package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "list container [filter]",
		Short: "List objects in a container",
		Long: `Lists all objects in a container; it can optionally apply a user-provided 
	regular expression to filter objects.`,
		Example: "  swift list my-container '^my-obj.*'",
		Aliases: []string{"search", "find", "l", "s", "f"},
		//Args:    cobra.RangeArgs(1, 2),
		Args: func(cmd *cobra.Command, args []string) error {
			err := (cobra.RangeArgs(1, 2))(cmd, args)
			if err != nil {
				return err
			}
			// add extra validation here
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("list called")
			// match := regexp.MustCompile(`^(?:/{0,1})([^/]+)(?:/{0,1})(.+)$`).FindStringSubmatch(os.Args[1])
			// if len(match) != 3 {
			// }
		},
	})
}
