package models

type DOMNode struct {
	NodeID     int
	Tag        string
	ID         string
	Classes    []string
	Attributes map[string]string
	Parent     *DOMNode
	Children   []*DOMNode
	Depth      int
}
