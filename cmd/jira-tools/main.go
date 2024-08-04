package main

import (
	"fmt"
	"main/internal/handlers"
	"main/internal/services"
	"os"
	"strconv"

	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/log"
	"github.com/go-ldap/ldap/v3"
	"github.com/joho/godotenv"
)

func init() {
	if envErr := godotenv.Load(".env"); envErr != nil {
		fmt.Println(".env file missing")
	}
}

func main() {
	tp := jira.BearerAuthTransport{
		Token: os.Getenv("JIRA_TOKEN"),
	}

	client, err := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))
	if err != nil {
		log.Fatal(err)
	}

	ldapURL := os.Getenv("LDAP_URL")
	conn, err := ldap.DialURL(ldapURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// err = conn.Bind(os.Getenv("ADMIN_DN"), os.Getenv("ADMIN_PASS"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	issueService := services.NewIssueService(client)
	assetService := services.NewAssetService(client)
	activeDirectoryService := services.NewActiveDirectoryService(conn)
	writeOffService := services.NewWriteOffService()

	issueHandler := handlers.NewIssueHandler(issueService, activeDirectoryService, assetService)
	assetHandler := handlers.NewAssetHandler(assetService)
	writeOffHandler := handlers.NewWriteOffHandler(writeOffService, assetService)

	fmt.Print("\033[H\033[2J")
	fmt.Println("1) Process deactivate insight account issue")
	fmt.Println("2) Generate write-off records")
	fmt.Println("3) Get laptop description")
	fmt.Println("4) Assign all deactivate insight issues to me")
	fmt.Println("5) Show issues with empty component")
	fmt.Println("6) Update block trainee cc issue")
	// fmt.Println("7) Remove VPN groups from users")
	fmt.Println("8) Process grant access issue")
	fmt.Println("9) Process disable active directory account issue")
	// fmt.Println("10) Move users to new OU")
	fmt.Println("11) Add user to ad group from cli")
	fmt.Println("12) Process dismissal or hiring issue")
	fmt.Println("13) Process return cc equipment issue")

	var n int
	for {
		var input string
		fmt.Print("Your choice: ")
		fmt.Scanln(&input)

		n, err = strconv.Atoi(input)
		if err == nil && (1 <= n && n <= 13) {
			break
		}
		fmt.Println("Invalid choice")
	}

	if n == 1 {
		err := issueHandler.ProcessDeactivateInsightAccountIssue()
		if err != nil {
			log.Fatal(err)
		}
	}

	if n == 2 {
		err := writeOffHandler.GenerateWriteOffRecords()
		if err != nil {
			log.Fatal(err)
		}
	}

	if n == 3 {
		err := assetHandler.PrintLaptopDescription()
		if err != nil {
			log.Fatal(err)
		}
	}

	if n == 4 {
		err := issueHandler.AssignAllDeactivateInsightIssuesToMe()
		if err != nil {
			log.Fatal(err)
		}
	}

	if n == 5 {
		err := issueHandler.ShowIssuesWithEmptyComponent()
		if err != nil {
			log.Fatal(err)
		}
	}

	if n == 6 {
		err := issueHandler.UpdateBlockTraineeIssue()
		if err != nil {
			log.Fatal(err)
		}
	}

	// if n == 7 {
	// 	err := activeDirectoryService.RemoveVPNGroupsFromUsers()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	if n == 8 {
		err := issueHandler.ProcessGrantAccessIssue()
		if err != nil {
			log.Fatal(err)
		}
	}

	if n == 9 {
		err := issueHandler.ProcessDisableActiveDirectoryAccountIssue()
		if err != nil {
			log.Fatal(err)
		}
	}

	// if n == 10 {
	// 	err := th.MoveUsersToNewOU()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	if n == 11 {
		err := issueHandler.AddUserToGroupFromCLI()
		if err != nil {
			log.Fatal(err)
		}
	}

	if n == 12 {
		err := issueHandler.ProcessDismissalOrHiringIssue()
		if err != nil {
			log.Fatal(err)
		}
	}

	if n == 13 {
		err := issueHandler.ProcessReturnCCEquipmentIssue()
		if err != nil {
			log.Fatal(err)
		}
	}

}
