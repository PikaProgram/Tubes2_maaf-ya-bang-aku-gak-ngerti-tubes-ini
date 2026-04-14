package utils

import (
	"backend/models"
	"fmt"
	"strings"
)

func PrintDOMTree(node *models.DOMNode, prefix string, isLast bool) {
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	var label strings.Builder
	label.WriteString("<" + node.Tag + ">" + fmt.Sprintf(" (NodeID: %d)", node.NodeID))

	for _, class := range node.Classes {
		label.WriteString("." + class)
	}
	if id, ok := node.Attributes["id"]; ok {
		label.WriteString(fmt.Sprintf(" #%s", id))
	}
	for key, val := range node.Attributes {
		if key != "id" && key != "class" {
			label.WriteString(fmt.Sprintf(" [%s=%q]", key, val))
		}
	}

	fmt.Println(prefix + connector + label.String())

	childPrefix := prefix + "│   "
	if isLast {
		childPrefix = prefix + "    "
	}

	for i, child := range node.Children {
		PrintDOMTree(child, childPrefix, i == len(node.Children)-1)
	}
}

func PrintTree(root *models.DOMNode) {
	fmt.Printf("<%s> (NodeID: %d)", root.Tag, root.NodeID)
	for i, child := range root.Children {
		PrintDOMTree(child, "", i == len(root.Children)-1)
	}
}
