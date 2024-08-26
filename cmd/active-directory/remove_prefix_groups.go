package active_directory

import (
	"log"

	"github.com/spf13/cobra"
)

var removePrefixGroupsCmd = &cobra.Command{
	Use:   "remove-prefix-groups",
	Short: "Remove groups that have a given prefix",
	Run: func(cmd *cobra.Command, args []string) {
		err := activeDirectoryHandler.RemovePrefixGroupsFromUsers()
		if err != nil {
			log.Fatal(err)
		}
	},
}
