package main

import (
	"backend/services"
	"backend/services/parser"
	"backend/services/search"
	"backend/utils"
	"fmt"
	"log"
)

func main() {
	targetURL := "https://google.com"

	fmt.Println("Fetching:", targetURL)
	rawHTML, err := services.FetchHTMLPage(targetURL)
	if err != nil {
		log.Fatal("Scraper error:", err)
	}
	fmt.Printf("Fetched %d bytes\n\n", len(rawHTML))

	root, err := parser.ParseHTML(rawHTML)
	if err != nil {
		log.Fatal("Parser error:", err)
	}

	utils.PrintTree(root)

	selector, err := parser.ParseCSSSelector("div")
	if err != nil {
		log.Fatal("Selector parse error:", err)
	}
	fmt.Printf("Parsed selector with %d step(s)\n", len(selector.Steps))

	for i, step := range selector.Steps {
		fmt.Printf("Step %d: Combinator='%s', Tag=%q, ID=%q, Classes=%v",
			i+1, step.Combinator, step.Compound.Tag, step.Compound.ID, step.Compound.Classes)
		for _, attr := range step.Compound.Attributes {
			fmt.Printf(", Attr{Name=%q, Operator='%s', Value=%q}", attr.Name, attr.Operator, attr.Value)
		}
		fmt.Println()
	}

	searchResult, searchLog := search.SearchElementDFS(root, &selector)

	fmt.Printf("Found %d matching element(s) with DFS:\n", len(searchResult.NodeIDs))
	fmt.Printf("Matching NodeIDs: %v\n", searchResult.NodeIDs)
	for _, nodeID := range searchResult.NodeIDs {
		res := searchResult.Results[nodeID]
		fmt.Printf("- NodeID: %d, Tag: <%s>, Path: %v\n", res.Node.NodeID, res.Node.Tag, res.Path)
	}

	fmt.Printf("Search Log (DFS):\n")
	for _, entry := range searchLog.Entries {
		fmt.Printf("  - NodeID: %d, Depth: %d\n", entry.NodeID, entry.Depth)
	}

	searchResultBFS, searchLogBFS := search.SearchElementBFS(root, &selector)

	fmt.Printf("Found %d matching element(s) with BFS:\n", len(searchResultBFS.NodeIDs))
	fmt.Printf("Matching NodeIDs: %v\n", searchResultBFS.NodeIDs)
	for _, nodeID := range searchResultBFS.NodeIDs {
		res := searchResultBFS.Results[nodeID]
		fmt.Printf("- NodeID: %d, Tag: <%s>, Path: %v\n", res.Node.NodeID, res.Node.Tag, res.Path)
	}

	fmt.Printf("Search Log (BFS):\n")
	for _, entry := range searchLogBFS.Entries {
		fmt.Printf("  - NodeID: %d, Depth: %d\n", entry.NodeID, entry.Depth)
	}
}
