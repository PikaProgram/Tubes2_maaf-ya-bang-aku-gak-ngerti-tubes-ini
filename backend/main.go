package main

import (
	"backend/services"
	"backend/utils"
	"fmt"
	"log"
)

func main() {
	targetURL := "https://motherfuckingwebsite.com/"

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

	utils.PrintTree(root)
}
