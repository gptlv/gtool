package ad

import (
	"log"

	"github.com/spf13/cobra"
)

var moveUsersCmd = &cobra.Command{
	Use:   "move-users",
	Short: "Move users to a new ou",
	Run: func(cmd *cobra.Command, args []string) {
		err := activeDirectoryHandler.MoveUsersToNewOU()
		if err != nil {
			log.Fatal(err)
		}
	},
}
