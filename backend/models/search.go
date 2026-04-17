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

func (res *SearchResult) Serialize() map[int]interface{} {
	serialized := make(map[int]interface{})
	for nodeID, result := range res.Results {
		serialized[nodeID] = map[string]interface{}{
			"node": map[string]interface{}{
				"NodeID":     result.Node.NodeID,
				"Tag":        result.Node.Tag,
				"ID":         result.Node.ID,
				"Classes":    result.Node.Classes,
				"Attributes": result.Node.Attributes,
				"Content":    result.Node.Content,
				"Depth":      result.Node.Depth,
			},
			"path": result.Path,
		}
	}
	return serialized
}

func (log *SearchLog) Serialize() map[string]interface{} {
	entries := make([]int, len(log.Entries))
	for i, entry := range log.Entries {
		entries[i] = entry.NodeID
	}
	return map[string]interface{}{
		"Selector":   log.Selector.String(),
		"SearchType": log.SearchType,
		"Entries":    entries,
	}
}
