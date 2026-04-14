package models

import (
	"fmt"
	"strings"
)

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

func (node *DOMNode) MatchesSelector(step *SelectorStep) (bool, error) {
	if node == nil {
		return false, fmt.Errorf("DOMNode is nil")
	} else if step == nil {
		return false, fmt.Errorf("SelectorStep is nil")
	}

	if step.Compound.Tag != "" && node.Tag != step.Compound.Tag {
		return false, nil
	}

	if step.Compound.ID != "" && node.ID != step.Compound.ID {
		return false, nil
	}

	for _, class := range step.Compound.Classes {
		if !strings.Contains(strings.Join(node.Classes, " "), class) {
			return false, nil
		}
	}

	for _, attr := range step.Compound.Attributes {
		nodeAttrValue, exists := node.Attributes[attr.Name]
		switch attr.Operator {
		case AttrOperatorExists:
			if !exists {
				return false, nil
			}
		case AttrOperatorEquals:
			if !exists || nodeAttrValue != attr.Value {
				return false, nil
			}
		case AttrOperatorIncludes:
			if !exists || !strings.Contains(strings.Join(strings.Fields(nodeAttrValue), " "), attr.Value) {
				return false, nil
			}
		case AttrOperatorDashMatch:
			if !exists || (nodeAttrValue != attr.Value && !strings.HasPrefix(nodeAttrValue, attr.Value+"-")) {
				return false, nil
			}
		case AttrOperatorPrefixMatch:
			if !exists || !strings.HasPrefix(nodeAttrValue, attr.Value) {
				return false, nil
			}
		case AttrOperatorSuffixMatch:
			if !exists || !strings.HasSuffix(nodeAttrValue, attr.Value) {
				return false, nil
			}
		case AttrOperatorSubstringMatch:
			if !exists || !strings.Contains(nodeAttrValue, attr.Value) {
				return false, nil
			}
		default:
			return false, fmt.Errorf("unknown attribute operator: %s", attr.Operator)
		}
	}

	return true, nil
}
