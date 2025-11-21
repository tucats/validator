package validator

import (
	"fmt"
	"strings"
	"text/scanner"
)

type tokenizer struct {
	position int
	tokens   []token
}

type token struct {
	line     int
	column   int
	spelling string
}

func (t *tokenizer) next() string {
	if t.position >= len(t.tokens) {
		return ""
	}

	token := t.tokens[t.position]
	t.position++

	return token.spelling
}

// Move the tokenizer position using the given offset. A negative offset will move the
// tokenizer position backwards. An offset of zero does nothing.
func (t *tokenizer) move(offset int) {
	if t.position+offset < 0 || t.position+offset >= len(t.tokens) {
		return
	}

	t.position += offset
}

// Fetch a token from the token queue, using the provided offset. An
// of zero of zero will return the current token.
func (t *tokenizer) peek(offset int) string {
	if t.position+offset < 0 || t.position+offset >= len(t.tokens) {
		return ""
	}

	token := t.tokens[t.position+offset]

	return token.spelling
}

func (t *tokenizer) pos() string {
	i := t.position - 1

	return fmt.Sprintf("line %d, column %d", t.tokens[i].line, t.tokens[i].column)
}

func Compile(src string) (*Item, error) {
	var (
		s      scanner.Scanner
		errors []error
	)

	src = UpdateLineEndings(src)

	s.Init(strings.NewReader(src))

	// Redirect any lexical scanning errors to the tokenizer log, if enabled.
	s.Error = func(s *scanner.Scanner, msg string) {
		errors = append(errors, NewError(msg))
	}

	s.Filename = "Input"
	tokenizer := &tokenizer{
		position: 0,
		tokens:   make([]token, 0),
	}

	// Scan as long as there are tokens left.
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		tokenizer.tokens = append(tokenizer.tokens, token{
			line:     s.Line,
			column:   s.Position.Column,
			spelling: s.TokenText(),
		})
	}

	return compileItem(tokenizer)
}

func UpdateLineEndings(src string) string {
	// Add line endings where needed in input source.
	lines := strings.Split(src, "\n")
	newSource := []string{}

	for i := range lines {
		text := lines[i]
		if strings.HasPrefix(text, "#") || strings.HasPrefix(text, "//") {
			text = ""
		}

		if len(text) > 0 {
			for strings.HasSuffix(text, " ") {
				text = text[:len(text)-1]
			}

			ch := text[len(text)-1:]
			if ch == "{" || ch == "," || ch == ";" {
				// no action needed
			} else {
				text += ";"
			}
		}

		newSource = append(newSource, text)
	}

	return strings.Join(newSource, "\n")
}

func compileItem(t *tokenizer) (*Item, error) {
	item := &Item{}

	if err := compileType(t, item); err != nil {
		return nil, err
	}

	if t.peek(0) == ";" {
		return item, nil
	}

	next := t.next()
	if next == ":" {
		if err := compileAttributes(t, item); err != nil {
			return nil, err
		}

		return item, nil
	}

	if next == "{" {
		err := compileObject(t, item)

		return item, err
	}

	return nil, ErrSyntaxError.Context(t.pos()).Value(next).Expected(";", ":", "{")
}

func compileAttributes(t *tokenizer, item *Item) error {
	text := ""

	next := t.next()
	for next != ";" {
		text += next + " "
		next = t.next()
	}

	return item.ParseTag(text)
}

func compileObject(t *tokenizer, item *Item) error {
	item.ItemType = TypeStruct

	for {
		// Compile the field
		field, err := compileItem(t)
		if err != nil {
			return err
		}

		item.Fields = append(item.Fields, field)

		if t.peek(0) == "}" {
			t.next()

			return nil
		}
	}
}

func compileType(t *tokenizer, item *Item) error {
	array := false
	ptr := false
	name := ""
	next := t.next()

	if next == "*" {
		ptr = true

		next = t.next()
	}

	if next == "[" && t.peek(0) == "]" {
		array = true

		t.next() // consume "]"
	}

	// Is it a reserved type word? If not, assume it's the item name, and
	// fetch the type name that follows it.
	if _, ok := TypeNamesMap[next]; !ok {
		name = next

		next = t.next()

		if next == ":" || next == ";" {
			return ErrUnsupportedType.Context(t.pos()).Value(name)
		}
	}

	item.Name = name

	if next == "{" {
		item.ItemType = TypeStruct

		t.move(-1)

		return nil
	} else {
		if kind, ok := TypeNamesMap[next]; ok {
			if array {
				item.ItemType = TypeArray
				item.BaseType = NewType(kind)
			} else if ptr {
				item.ItemType = TypePointer
				item.BaseType = NewType(kind)
			} else {
				item.ItemType = kind
			}

			return nil
		}
	}

	return ErrUnsupportedType.Value(next).Context(t.pos())
}
