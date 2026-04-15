package search

import "backend/models"

type QueueItem struct {
	Node      *models.DOMNode
	StepIndex int
}

func SearchElementBFS(root *models.DOMNode, selector *models.Selector) (*models.SearchResult, *models.SearchLog) {
	if root == nil || selector == nil {
		return nil, nil
	}

	results := models.SearchResult{
		NodeIDs: []int{},
		Results: make(map[int]models.SelectorResult),
	}

	log := models.SearchLog{
		Selector:   *selector,
		SearchType: "BFS",
		Entries:    []models.SearchLogEntry{},
	}

	queue := []QueueItem{{Node: root, StepIndex: 0}}

	for len(queue) > 0 {
		currentItem := queue[0]
		queue = queue[1:]

		currentNode := currentItem.Node
		currentStepIndex := currentItem.StepIndex

		if currentNode == nil || currentStepIndex < 0 || currentStepIndex >= len(selector.Steps) {
			continue
		}

		log.Entries = append(log.Entries, models.SearchLogEntry{
			NodeID: currentNode.NodeID,
			Depth:  currentNode.Depth,
		})

		step := selector.Steps[currentStepIndex]
		if step.Compound.Matches(currentNode) {
			if currentStepIndex == len(selector.Steps)-1 {
				path := make([]int, 0, currentNode.Depth+1)
				for current := currentNode; current != nil; current = current.Parent {
					path = append([]int{current.NodeID}, path...)
				}

				results.Results[currentNode.NodeID] = models.SelectorResult{
					Node: currentNode,
					Path: path,
				}

				results.NodeIDs = append(results.NodeIDs, currentNode.NodeID)
			} else {
				relatedNodes := currentNode.GetRelatedNodes(selector.Steps[currentStepIndex+1].Combinator)
				for _, relatedNode := range relatedNodes {
					queue = append(queue, QueueItem{Node: relatedNode, StepIndex: currentStepIndex + 1})
				}
			}
		}

		for _, child := range currentNode.Children {
			queue = append(queue, QueueItem{Node: child, StepIndex: currentStepIndex})
		}
	}

	return &results, &log
}
