package parser

import (
	"backend/models"
	"strings"
	"unicode"
)

var voidElements = map[string]bool{
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"link":   true,
	"meta":   true,
	"param":  true,
	"source": true,
	"track":  true,
	"wbr":    true,
}

var rawTextElements = map[string]bool{
	"script": true,
	"style":  true,
}

type tokenType int

const (
	tokenEOF tokenType = iota
	tokenDoctype
	tokenComment
	tokenOpenTag
	tokenCloseTag
	tokenText
)

type htmlAttr struct {
	key string
	val string
}

type htmlToken struct {
	typ       tokenType
	tag       string
	attrs     []htmlAttr
	selfClose bool
}

type htmlTokenizer struct {
	src string
	pos int
}

func newHTMLTokenizer(src string) *htmlTokenizer {
	return &htmlTokenizer{src: src}
}

func (t *htmlTokenizer) skipWhitespace() {
	for t.pos < len(t.src) && unicode.IsSpace(rune(t.src[t.pos])) {
		t.pos++
	}
}

func (t *htmlTokenizer) consumeUntilString(s string) {
	idx := strings.Index(t.src[t.pos:], s)
	if idx == -1 {
		t.pos = len(t.src)
		return
	}
	t.pos += idx + len(s)
}

func (t *htmlTokenizer) readTagName() string {
	start := t.pos
	for t.pos < len(t.src) {
		ch := t.src[t.pos]
		if ch == '>' || ch == '/' || unicode.IsSpace(rune(ch)) {
			break
		}
		t.pos++
	}
	return strings.ToLower(t.src[start:t.pos])
}

func (t *htmlTokenizer) readAttrName() string {
	start := t.pos
	for t.pos < len(t.src) {
		ch := t.src[t.pos]
		if ch == '=' || ch == '>' || ch == '/' || unicode.IsSpace(rune(ch)) {
			break
		}
		t.pos++
	}
	return strings.ToLower(t.src[start:t.pos])
}

func (t *htmlTokenizer) readAttrValue() string {
	if t.pos >= len(t.src) {
		return ""
	}
	ch := t.src[t.pos]
	if ch == '"' || ch == '\'' {
		quote := ch
		t.pos++
		start := t.pos
		for t.pos < len(t.src) && t.src[t.pos] != quote {
			t.pos++
		}
		val := t.src[start:t.pos]
		if t.pos < len(t.src) {
			t.pos++
		}
		return val
	}
	start := t.pos
	for t.pos < len(t.src) {
		c := t.src[t.pos]
		if c == '>' || unicode.IsSpace(rune(c)) {
			break
		}
		t.pos++
	}
	return t.src[start:t.pos]
}

func (t *htmlTokenizer) readAttributes() []htmlAttr {
	var attrs []htmlAttr
	for {
		t.skipWhitespace()
		if t.pos >= len(t.src) {
			break
		}
		ch := t.src[t.pos]
		if ch == '>' || ch == '/' {
			break
		}
		name := t.readAttrName()
		if name == "" {
			t.pos++
			continue
		}
		t.skipWhitespace()
		if t.pos < len(t.src) && t.src[t.pos] == '=' {
			t.pos++
			t.skipWhitespace()
			val := t.readAttrValue()
			attrs = append(attrs, htmlAttr{key: name, val: val})
		} else {
			attrs = append(attrs, htmlAttr{key: name, val: ""})
		}
	}
	return attrs
}

// <script> and <style>
func (t *htmlTokenizer) consumeRawText(tag string) {
	end := "</" + tag
	src := t.src
	for t.pos < len(src) {
		remaining := len(src) - t.pos
		if remaining >= len(end) && strings.EqualFold(src[t.pos:t.pos+len(end)], end) {
			t.pos += len(end)
			for t.pos < len(src) && src[t.pos] != '>' {
				t.pos++
			}
			if t.pos < len(src) {
				t.pos++
			}
			return
		}
		t.pos++
	}
}

func (t *htmlTokenizer) next() htmlToken {
	if t.pos >= len(t.src) {
		return htmlToken{typ: tokenEOF}
	}
	if t.src[t.pos] != '<' {
		start := t.pos
		for t.pos < len(t.src) && t.src[t.pos] != '<' {
			t.pos++
		}
		return htmlToken{typ: tokenText, tag: t.src[start:t.pos]}
	}

	t.pos++

	if t.pos >= len(t.src) {
		return htmlToken{typ: tokenEOF}
	}

	ch := t.src[t.pos]

	// <!-- comment --> or <!DOCTYPE ...>
	if ch == '!' {
		t.pos++
		if t.pos+1 < len(t.src) && t.src[t.pos] == '-' && t.src[t.pos+1] == '-' {
			t.pos += 2
			t.consumeUntilString("-->")
			return htmlToken{typ: tokenComment}
		}
		t.consumeUntilString(">")
		return htmlToken{typ: tokenDoctype}
	}

	// </closing>
	if ch == '/' {
		t.pos++
		name := t.readTagName()
		t.consumeUntilString(">")
		return htmlToken{typ: tokenCloseTag, tag: name}
	}

	// <opening ...>
	name := t.readTagName()
	if name == "" {
		t.consumeUntilString(">")
		return htmlToken{typ: tokenComment}
	}
	attrs := t.readAttributes()

	selfClose := false
	if t.pos < len(t.src) && t.src[t.pos] == '/' {
		selfClose = true
		t.pos++
	}
	if t.pos < len(t.src) && t.src[t.pos] == '>' {
		t.pos++
	}
	if voidElements[name] {
		selfClose = true
	}
	return htmlToken{typ: tokenOpenTag, tag: name, attrs: attrs, selfClose: selfClose}
}

func ParseHTML(rawHTML string) (*models.DOMNode, error) {
	currentNodeID := 0
	nextNodeID := func() int {
		id := currentNodeID
		currentNodeID++
		return id
	}

	tkn := newHTMLTokenizer(rawHTML)
	root := &models.DOMNode{
		NodeID:     nextNodeID(),
		Tag:        "#document",
		Classes:    []string{},
		Attributes: make(map[string]string),
		Depth:      0,
	}
	stack := []*models.DOMNode{root}
	for {
		token := tkn.next()
		if token.typ == tokenEOF {
			break
		}

		if token.typ == tokenOpenTag {
			parent := stack[len(stack)-1]
			node := &models.DOMNode{
				NodeID:     nextNodeID(),
				Tag:        token.tag,
				Classes:    []string{},
				Attributes: make(map[string]string),
				Parent:     parent,
				Depth:      parent.Depth + 1,
			}

			for _, attr := range token.attrs {
				switch attr.key {
				case "id":
					node.ID = attr.val
					node.Attributes["id"] = attr.val
				case "class":
					node.Classes = strings.Fields(attr.val)
					node.Attributes["class"] = attr.val
				default:
					node.Attributes[attr.key] = attr.val
				}
			}

			parent.Children = append(parent.Children, node)

			if !token.selfClose {
				stack = append(stack, node)
				if rawTextElements[token.tag] {
					tkn.consumeRawText(token.tag)
					stack = stack[:len(stack)-1]
				}
			}
		} else if token.typ == tokenCloseTag {
			for i := len(stack) - 1; i >= 1; i-- {
				if stack[i].Tag == token.tag {
					stack = stack[:i]
					break
				}
			}
		}
	}

	for _, child := range root.Children {
		if child.Tag == "html" {
			return child, nil
		}
	}
	if len(root.Children) > 0 {
		return root.Children[0], nil
	}
	return root, nil
}
