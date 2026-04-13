package main

import (
    "backend/services"
    "fmt"
    "log"
)

func main() {
    rawHTML := `<html><head><title>Test</title></head><body><div class="a"><p id="b">Hello</p></div></body></html>`
    doc, err := services.ParseHTML(rawHTML)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Node awal: ", doc.Type)
}