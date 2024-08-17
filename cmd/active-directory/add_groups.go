package active_directory

import (
	"log"

	"github.com/spf13/cobra"
)

var addGroupsCmd = &cobra.Command{
	Use:   "add-group",
	Short: "Add user to group from cli",
	Run: func(cmd *cobra.Command, args []string) {
		err := activeDirectoryHandler.AddUsersToGroupsFromCLI()
		if err != nil {
			log.Fatal(err)
		}
	},
}
