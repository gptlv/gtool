package handlers

import (
	"fmt"
	"io"
	"main/internal/interfaces"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-ldap/ldap/v3"
)

type activeDirectoryHandler struct {
	activeDirectoryService interfaces.ActiveDirectoryService
}

func NewActiveDirectoryHandler(activeDirectoryService interfaces.ActiveDirectoryService) *activeDirectoryHandler {
	return &activeDirectoryHandler{activeDirectoryService: activeDirectoryService}
}

func (activeDirectoryHandler *activeDirectoryHandler) RemoveVPNGroupsFromUsers() error {
	f, err := os.OpenFile("vpn.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	wrt := io.MultiWriter(f, os.Stdout)
	log.SetOutput(wrt)

	var failedUsers []string
	failedFile, err := os.Create("failed_users.txt")
	if err != nil {
		return fmt.Errorf("failed to create failed users file: %w", err)
	}
	defer failedFile.Close()

	prefix := "res.vpn."
	usersFilename := "users.txt"

	log.Info(fmt.Sprintf("Opening %s file", usersFilename))
	file, err := os.Open(usersFilename)
	if err != nil {
		log.Error(fmt.Errorf("failed to open file %s: %w", usersFilename, err))
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Error(fmt.Errorf("failed to read file %s: %w", usersFilename, err))
		return err
	}

	emails := strings.Split(string(content), "\n")

	for _, email := range emails {
		log.Info(fmt.Sprintf("Getting user by email %s", email))
		user, err := activeDirectoryHandler.activeDirectoryService.GetByEmail(email)
		if err != nil {
			log.Error(fmt.Errorf("failed to get user %s by email: %w", email, err))
			continue
		}

		userGroups := activeDirectoryHandler.activeDirectoryService.GetUserGroups(user)

		for _, groupDN := range userGroups {
			groupCN, err := activeDirectoryHandler.activeDirectoryService.ExtractCNFromDN(groupDN)
			if err != nil {
				log.Error(fmt.Errorf("failed to extract group CN from DN %s: %w", groupDN, err))
				continue
			}

			if strings.HasPrefix(groupCN, prefix) {
				group := &ldap.Entry{DN: groupDN}

				log.Info(fmt.Sprintf("Removing group %s from user %s", groupCN, user.GetAttributeValue("sAMAccountName")))

				_, err := activeDirectoryHandler.activeDirectoryService.RemoveUserFromGroup(user, group)
				if err != nil {
					log.Error(fmt.Errorf("failed to remove user %s from group %s: %w", email, groupCN, err))
					failedUsers = append(failedUsers, email)
					continue
				}
				time.Sleep(time.Second)
			}

		}
	}

	for _, failedUser := range failedUsers {
		_, err := failedFile.WriteString(failedUser + "\n")
		if err != nil {
			log.Error(fmt.Errorf("failed to write failed user %s to file: %v", failedUser, err))
		}
	}

	log.Info("Task completed successfully")
	return nil
}

func (activeDirectoryHandler *activeDirectoryHandler) MoveUsersToNewOU() error {
	f, err := os.OpenFile("ou.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	wrt := io.MultiWriter(f, os.Stdout)
	log.SetOutput(wrt)

	usersFilename := "users.txt"

	log.Info(fmt.Sprintf("Opening %s file", usersFilename))
	file, err := os.Open(usersFilename)
	if err != nil {
		log.Error(fmt.Errorf("failed to open file %s: %w", usersFilename, err))
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Error(fmt.Errorf("failed to read file %s: %w", usersFilename, err))
		return err
	}

	commonNames := strings.Split(string(content), "\n")

	for _, cn := range commonNames {
		user, err := activeDirectoryHandler.activeDirectoryService.GetByCN(cn)
		if err != nil {
			log.Error(fmt.Errorf("failed to get user %v by cn", cn))
			return err
		}

		newSup := "OU=External,OU=Users,OU=SBMT,DC=sbermarket,DC=ru"

		log.Info(fmt.Sprintf("Moving user %v to %v", cn, newSup))
		_, err = activeDirectoryHandler.activeDirectoryService.UpdateDN(user, newSup)
		if err != nil {
			log.Error(fmt.Errorf("failed to move user %v to %v", cn, newSup))
			return err
		}

		time.Sleep(time.Second)

	}

	return nil
}
