package models

type DismissalRecord struct {
	//comes from csv
	ID       string
	ISC      string
	Flaw     string
	Decision string
	//from insight
	Serial      string
	Name        string
	InventoryID string
	//common
	Date string
	Boss string
	Lead string
}
