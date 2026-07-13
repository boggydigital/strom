package strom

import (
	"errors"
	"io"
	"iter"
	"maps"
	"slices"
	"strings"
	"sync"
)

const (
	openingAngleBracket   = "<"
	closingAngleBracket   = ">"
	forwardSlash          = "/"
	equalSign             = "="
	singleQuote           = "'"
	colon                 = ":"
	semicolon             = ";"
	classAttributeName    = "class"
	styleAttributeName    = "style"
	singleSpace           = " "
	openingCommentTagName = "<!--"
	closingCommentTagName = "-->"
)

var restrictedAttributes = []string{classAttributeName, styleAttributeName}

type elementNode struct {
	tagName             string
	textContent         []byte
	classes             map[string]any
	attributes          map[string]string
	styles              map[string]string
	prefix              []byte
	suffix              []byte
	children            []Element
	getChildrenDelegate func() iter.Seq[Element]
	mtx                 *sync.Mutex
}

type Element interface {
	SetTextContent(textContent string) Element
	AddClass(classes ...string) Element
	SetAttribute(name string, value string) Element
	SetStyles(styles map[string]string) Element
	Append(nodes ...Element) Element
	Write(w io.Writer) error
}

func (e *elementNode) SetTextContent(textContent string) Element {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	e.textContent = []byte(textContent)
	return e
}

func (e *elementNode) AddClass(classes ...string) Element {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	if e.classes == nil {
		e.classes = make(map[string]any)
	}

	for _, class := range classes {
		e.classes[class] = nil
	}

	return e
}

func (e *elementNode) SetStyles(styles map[string]string) Element {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	if e.styles == nil {
		e.styles = make(map[string]string)
	}

	for p, v := range styles {
		e.styles[p] = v
	}

	return e
}

func (e *elementNode) classList() string {
	return strings.Join(slices.Collect(maps.Keys(e.classes)), singleSpace)
}

func (e *elementNode) SetAttribute(name string, value string) Element {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	if e.attributes == nil {
		e.attributes = make(map[string]string)
	}

	e.attributes[name] = value
	return e
}

func (e *elementNode) Append(nodes ...Element) Element {
	e.children = append(e.children, nodes...)
	return e
}

func writeStrings(w io.Writer, parts ...string) error {
	for _, p := range parts {
		if _, err := w.Write([]byte(p)); err != nil {
			return err
		}
	}
	return nil
}

func writeAttribute(w io.Writer, name, value string) error {
	return writeStrings(w, singleSpace, name, equalSign, singleQuote, value, singleQuote)
}

func writeStyles(w io.Writer, styles map[string]string) error {
	if err := writeStrings(w, singleSpace, styleAttributeName, equalSign, singleQuote); err != nil {
		return err
	}

	for p, v := range styles {
		if err := writeStrings(w, p, colon, v, semicolon); err != nil {
			return err
		}
	}

	return writeStrings(w, singleQuote)
}

func (e *elementNode) Write(w io.Writer) error {

	e.mtx.Lock()
	defer e.mtx.Unlock()

	switch e.prefix {
	case nil:
	// do nothing
	default:
		if _, err := w.Write(e.prefix); err != nil {
			return err
		}
	}

	switch e.tagName {
	case "":
		if len(e.classes) > 0 {
			return errors.New("transparent containers cannot have classes")
		}
		if len(e.attributes) > 0 {
			return errors.New("transparent containers cannot have attributes")
		}
	default:
		if err := writeStrings(w, openingAngleBracket, e.tagName); err != nil {
			return nil
		}

		if len(e.classes) > 0 {

			for _, ra := range restrictedAttributes {
				if _, ok := e.attributes[ra]; ok {
					return errors.New("restrictred attribute " + ra)
				}
			}

			if err := writeAttribute(w, classAttributeName, e.classList()); err != nil {
				return err
			}
		}

		if len(e.styles) > 0 {
			if err := writeStyles(w, e.styles); err != nil {
				return err
			}
		}

		for attributeName, attributeValue := range e.attributes {
			if err := writeAttribute(w, attributeName, attributeValue); err != nil {
				return err
			}
		}

		if err := writeStrings(w, closingAngleBracket); err != nil {
			return err
		}

		if _, err := w.Write(e.textContent); err != nil {
			return err
		}
	}

	var children iter.Seq[Element]

	switch e.getChildrenDelegate {
	case nil:
		children = slices.Values(e.children)
	default:
		children = e.getChildrenDelegate()
	}

	if children != nil {
		for child := range children {
			if err := child.Write(w); err != nil {
				return err
			}
		}
	}

	switch e.tagName {
	case "":
	// do nothing
	default:
		if err := writeStrings(w, openingAngleBracket, forwardSlash, e.tagName, closingAngleBracket); err != nil {
			return err
		}
	}

	switch e.suffix {
	case nil:
	// do nothing
	default:
		if _, err := w.Write(e.suffix); err != nil {
			return err
		}
	}

	return nil
}

func Create(tagName string) Element {
	return &elementNode{
		mtx:     new(sync.Mutex),
		tagName: tagName,
	}
}

func CreateText(tagName, textContent string) Element {
	return &elementNode{
		mtx:         new(sync.Mutex),
		tagName:     tagName,
		textContent: []byte(textContent),
	}
}

func Comment(tagName string) Element {
	return &elementNode{
		mtx:     new(sync.Mutex),
		prefix:  []byte(openingCommentTagName),
		tagName: tagName,
		suffix:  []byte(closingCommentTagName),
	}
}

func Defer(getChildredDelegate func() iter.Seq[Element]) Element {
	return &elementNode{
		getChildrenDelegate: getChildredDelegate,
		mtx:                 &sync.Mutex{},
	}
}
