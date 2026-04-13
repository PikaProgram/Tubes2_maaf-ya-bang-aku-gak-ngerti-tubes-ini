package utils

import (
	"backend/models"
	"fmt"
	"strings"
)

func PrintTree(node *models.DOMNode, prefix string, isLast bool) {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	label := "<" + node.Tag + ">"
	if id, ok := node.Attributes["id"]; ok && id != "" {
		label += " #" + id
	}
	if len(node.Classes) > 0 {
		label += " ." + strings.Join(node.Classes, ".")
	}

	fmt.Println(prefix + connector + label)

	childPrefix := prefix + "│   "
	if isLast {
		childPrefix = prefix + "    "
	}

	for i, child := range node.Children {
		PrintTree(child, childPrefix, i == len(node.Children)-1)
	}
}
