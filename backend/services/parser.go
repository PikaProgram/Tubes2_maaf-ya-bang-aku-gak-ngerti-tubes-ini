package services

import (
	"backend/models"
	"strings"
	"golang.org/x/net/html"
)

func ParseHTML(rawHTML string) (*html.Node, error) {
	return html.Parse(strings.NewReader(rawHTML))
}

func BuildTree(n *html.Node, parent *models.DOMNode, depth int) *models.DOMNode {
	if n.Type == html.DocumentNode {
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode {
				return BuildTree(child, parent, depth)
			}
		}
		return nil
	}

	if n.Type != html.ElementNode {
		return nil
	}

	node := &models.DOMNode{
		Tag: n.Data,
		Depth: depth,
		Parent: parent,
		Classes: []string{},
		Attributes: make(map[string]string),
	}

	for _, attr := range n.Attr {
        if attr.Key == "id" {
            node.ID = attr.Val
        } else if attr.Key == "class" {
            node.Classes = strings.Fields(attr.Val)
        } else {
            node.Attributes[attr.Key] = attr.Val
        }
    }

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		childNode := BuildTree(child, node, depth+1)
		if childNode != nil {
			node.Children = append(node.Children, childNode)
		}
	}

	return node
}
