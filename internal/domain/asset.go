package domain

type Asset struct {
	ID          string `csv:"id"`
	ISC         string `csv:"isc"`
	Flaw        string `csv:"flaw"`
	Decision    string `csv:"decision"`
	Serial      string `csv:"serial"`
	Name        string `csv:"name"`
	InventoryID string `csv:"inventory_id"`
}
