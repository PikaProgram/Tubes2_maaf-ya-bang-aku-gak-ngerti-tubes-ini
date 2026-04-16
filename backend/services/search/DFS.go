package search

import "backend/models"

func searchDFS(node *models.DOMNode, selector *models.Selector, stepIndex int, amount int, results *models.SearchResult, log *models.SearchLog) {
	if len(selector.Steps) == 0 || stepIndex < 0 || stepIndex >= len(selector.Steps) {
		return
	}
	if len(results.NodeIDs) >= amount {
		return
	}

	if stepIndex == 0 {
		traverse(node, func(node *models.DOMNode) bool {
			if len(results.NodeIDs) >= amount {
				return false
			}

			matchElement(node, selector, 0, amount, results, log)
			return len(results.NodeIDs) < amount
		})
		return
	}

	matchElement(node, selector, stepIndex, amount, results, log)
}

func SearchElementDFS(root *models.DOMNode, selector *models.Selector, amount int) (*models.SearchResult, *models.SearchLog) {
	if root == nil || selector == nil {
		return nil, nil
	}

	results := models.SearchResult{
		NodeIDs: []int{},
		Results: make(map[int]models.SelectorResult),
	}

	log := models.SearchLog{
		Selector:   *selector,
		SearchType: "DFS",
		Entries:    []models.SearchLogEntry{},
	}

	searchDFS(root, selector, 0, amount, &results, &log)
	return &results, &log
}

func matchElement(node *models.DOMNode, selector *models.Selector, stepIndex int, amount int, results *models.SearchResult, log *models.SearchLog) {
	if node == nil || selector == nil || results == nil {
		return
	}
	if stepIndex < 0 || stepIndex >= len(selector.Steps) {
		return
	}
	if len(results.NodeIDs) >= amount {
		return
	}

	log.Entries = append(log.Entries, models.SearchLogEntry{
		NodeID: node.NodeID,
		Depth:  node.Depth,
	})

	step := selector.Steps[stepIndex]
	if !step.Compound.Matches(node) {
		return
	}

	if stepIndex == len(selector.Steps)-1 {
		path := make([]int, 0, node.Depth+1)
		for current := node; current != nil; current = current.Parent {
			path = append([]int{current.NodeID}, path...)
		}

		if _, exists := (*results).Results[node.NodeID]; !exists {
			(*results).Results[node.NodeID] = models.SelectorResult{
				Node: node,
				Path: path,
			}

			(*results).NodeIDs = append((*results).NodeIDs, node.NodeID)
		}

		return
	}

	nextStep := selector.Steps[stepIndex+1]
	for _, candidate := range node.GetRelatedNodes(nextStep.Combinator) {
		if len(results.NodeIDs) >= amount {
			return
		}

		matchElement(candidate, selector, stepIndex+1, amount, results, log)
	}
}

func traverse(node *models.DOMNode, visit func(*models.DOMNode) bool) bool {
	if node == nil || visit == nil {
		return true
	}

	if !visit(node) {
		return false
	}

	for _, child := range node.Children {
		if !traverse(child, visit) {
			return false
		}
	}

	return true
}
