package main

import (
    "backend/services"
    "fmt"
    "log"
)

func main() {
    rawData, err := services.FetchHTMLPage("x.com")
    if err != nil {
        log.Fatalf("Error fetching: %v", err)
    }
    fmt.Println(rawData)

    parsedData, err := services.ParseHTML(rawData)
    if err != nil {
        log.Fatalf("Error parsing: %v", err)
    }
    fmt.Printf("%+v\n", parsedData)
}