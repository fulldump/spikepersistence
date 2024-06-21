package persistence

type Item struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Version     string     `json:"version"`
	Subitems    []*SubItem `json:"subitems"`
	Counter     int        `json:"counter"`
}

type SubItem struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}
