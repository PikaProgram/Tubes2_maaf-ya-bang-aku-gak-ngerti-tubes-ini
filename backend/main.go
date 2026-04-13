package main

import (
	"backend/services"
	"backend/utils"
	"fmt"
	"log"
)

func main() {
	targetURL := "https://www.instagram.com/rasyad_2771/"

	fmt.Println("Fetching:", targetURL)
	rawHTML, err := services.FetchHTMLPage(targetURL)
	if err != nil {
		log.Fatal("Scraper error:", err)
	}
	fmt.Printf("Fetched %d bytes\n\n", len(rawHTML))

	root, err := services.ParseHTML(rawHTML)
	if err != nil {
		log.Fatal("Parser error:", err)
	}

	fmt.Println("Tree:")
	fmt.Println("<" + root.Tag + ">")
	for i, child := range root.Children {
		utils.PrintTree(child, "", i == len(root.Children)-1)
	}
}
