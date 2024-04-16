package name

import "fmt"

type FullName struct {
	FirstName  string
	LastName   string
	Patronymic string
}

func (f FullName) GetFullFormat() string {
	return fmt.Sprintf("%s %s %s", f.LastName, f.FirstName, f.Patronymic)
}

func (f FullName) GetShortFormat() string {
	initials := string(f.FirstName[0]) + "." + string(f.Patronymic[0]) + "."
	return fmt.Sprintf("%s %s", f.LastName, initials)
}
