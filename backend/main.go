package main

import (
	"backend/services"
	"backend/services/parser"
	"backend/utils"
	"fmt"
	"log"
)

func main() {
	targetURL := "https://google.com/"

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

	selector, err := parser.ParseCSSSelector("body")
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

	for _, child := range root.Children {
		matches, err := child.MatchesSelector(&selector.Steps[0])
		if err != nil {
			log.Fatal("Validation error:", err)
		}
		fmt.Printf("Root child matches selector: %v\n", matches)
	}
}
