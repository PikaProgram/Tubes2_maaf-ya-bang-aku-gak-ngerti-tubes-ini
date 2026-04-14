package models

type Combinator string

const (
	CombinatorNone            Combinator = ""
	CombinatorDescendant      Combinator = " "
	CombinatorChild           Combinator = ">"
	CombinatorAdjacentSibling Combinator = "+"
	CombinatorGeneralSibling  Combinator = "~"
)

type AttrOperator string

const (
	AttrOperatorExists         AttrOperator = "exists"
	AttrOperatorEquals         AttrOperator = "="
	AttrOperatorIncludes       AttrOperator = "~="
	AttrOperatorDashMatch      AttrOperator = "|="
	AttrOperatorPrefixMatch    AttrOperator = "^="
	AttrOperatorSuffixMatch    AttrOperator = "$="
	AttrOperatorSubstringMatch AttrOperator = "*="
)

type AttributeSelector struct {
	Name     string
	Operator AttrOperator
	Value    string
}

type CompoundSelector struct {
	Tag        string
	ID         string
	Classes    []string
	Attributes []AttributeSelector
}

type SelectorStep struct {
	Combinator Combinator
	Compound   CompoundSelector
}

type Selector struct {
	Steps []SelectorStep
}
