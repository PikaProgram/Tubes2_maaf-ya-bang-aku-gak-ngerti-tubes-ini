package models

type SelectorResult struct {
	Node *DOMNode
	Path []int
}

type SearchResult struct {
	NodeIDs []int
	Results map[int]SelectorResult
}

type SearchLog struct {
	Selector   Selector
	SearchType string
	Entries    []SearchLogEntry
}

type SearchLogEntry struct {
	NodeID int
	Depth  int
}
