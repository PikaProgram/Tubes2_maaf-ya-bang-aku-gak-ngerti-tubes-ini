package services

import (
    "strings"
    "golang.org/x/net/html"
)

func ParseHTML(rawHTML string) (*html.Node, error) {
    doc, err := html.Parse(strings.NewReader(rawHTML))
    if err != nil {
        return nil, err
    }
    return doc, nil
}