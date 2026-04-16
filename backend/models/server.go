package models

type Request struct {
	URL      string `json:"url"`
	Type     string `json:"type"`
	Amount   int    `json:"amount"`
	Selector string `json:"selector"`
}
