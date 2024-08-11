package ad

import (
	"fmt"

	"github.com/gptlv/gtools/common"
	"github.com/gptlv/gtools/internal/handlers"
	"github.com/gptlv/gtools/internal/interfaces"
	"github.com/gptlv/gtools/internal/services"
	"github.com/spf13/cobra"
)

func init() {
	initActiveDirectoryHandler()
	ActiveDirectoryCmd.AddCommand(addGroupsCmd)
	ActiveDirectoryCmd.AddCommand(moveUsersCmd)
	ActiveDirectoryCmd.AddCommand(removePrefixGroupsCmd)
}

var activeDirectoryHandler interfaces.ActiveDirectoryHandler

var ActiveDirectoryCmd = &cobra.Command{
	Use:   "active-directory",
	Short: "Modify active directory entries",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'myapp active-directory --help' for more information on modifying active directory entries")
	},
}

func initActiveDirectoryHandler() error {
	conn, err := common.GetLDAPConnection()
	if err != nil {
		return fmt.Errorf("failed to get ldap connection: %w", err)
	}

	activeDirectoryService := services.NewActiveDirectoryService(conn)
	activeDirectoryHandler = handlers.NewActiveDirectoryHandler(activeDirectoryService)

	return nil
}
