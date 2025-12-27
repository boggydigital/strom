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
	openingAngleBracket = "<"
	closingAngleBracket = ">"
	forwardSlash        = "/"
	equalSign           = "="
	singleQuote         = "'"
	classAttributeName  = "class"
	singleSpace         = " "
)

type element struct {
	tagName     string
	textContent string
	classes     map[string]any
	attributes  map[string]string
	children    []Element
	deferrals   []func() iter.Seq[Element]
	mtx         *sync.Mutex
}

type Element interface {
	SetTextContent(textContent string) Element

	AddClass(classes ...string) Element

	SetAttribute(name string, value string) Element

	Append(nodes ...Element) Element
	Defer(d func() iter.Seq[Element]) Element

	Write(w io.Writer) error
}

func (e *element) SetTextContent(textContent string) Element {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	e.textContent = textContent
	return e
}

func (e *element) AddClass(classes ...string) Element {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	for _, class := range classes {
		e.classes[class] = nil
	}

	return e
}

func (e *element) classList() string {
	return strings.Join(slices.Collect(maps.Keys(e.classes)), singleSpace)
}

func (e *element) SetAttribute(name string, value string) Element {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	e.attributes[name] = value
	return e
}

func (e *element) Append(nodes ...Element) Element {
	e.children = append(e.children, nodes...)
	return e
}

func (e *element) Defer(d func() iter.Seq[Element]) Element {
	e.deferrals = append(e.deferrals, d)
	return e
}

func writeAttribute(sb *strings.Builder, name, value string) error {
	var err error
	if _, err = sb.WriteString(name); err != nil {
		return err
	}
	if _, err = sb.WriteString(equalSign); err != nil {
		return err
	}
	if _, err = sb.WriteString(singleQuote); err != nil {
		return err
	}
	if _, err = sb.WriteString(value); err != nil {
		return err
	}
	if _, err = sb.WriteString(singleQuote); err != nil {
		return err
	}
	if _, err = sb.WriteString(singleSpace); err != nil {
		return err
	}
	return nil
}

func (e *element) Write(w io.Writer) error {

	e.mtx.Lock()
	defer e.mtx.Unlock()

	sb := new(strings.Builder)

	var err error
	if _, err = sb.WriteString(openingAngleBracket); err != nil {
		return err
	}
	if _, err = sb.WriteString(e.tagName); err != nil {
		return err
	}

	if len(e.attributes) > 0 {
		if _, err = sb.WriteString(singleSpace); err != nil {
			return err
		}
	}

	if len(e.classes) > 0 {
		if err = writeAttribute(sb, classAttributeName, e.classList()); err != nil {
			return err
		}

		if _, ok := e.attributes[classAttributeName]; ok {
			return errors.New("you must add classes with AddClass, not SetAttribute")
		}
	}

	for attributeName, attributeValue := range e.attributes {
		if err = writeAttribute(sb, attributeName, attributeValue); err != nil {
			return err
		}
	}

	if _, err = sb.WriteString(closingAngleBracket); err != nil {
		return err
	}
	if _, err = sb.WriteString(e.textContent); err != nil {
		return err
	}

	if len(e.children) > 0 {

		// flush currently accumulated element content before writing children content

		if _, err = io.Copy(w, strings.NewReader(sb.String())); err != nil {
			return err
		}
		sb.Reset()

		// intentionally not writing children to strings.Builder to prioritize streaming content
		// over accumulating the full element subtree

		for _, child := range e.children {
			if err = child.Write(w); err != nil {
				return err
			}
		}
	}

	if len(e.deferrals) > 0 {

		// flush currently accumularted element content before writing deferred content
		if _, err = io.Copy(w, strings.NewReader(sb.String())); err != nil {
			return err
		}
		sb.Reset()

		for _, deferral := range e.deferrals {
			for d := range deferral() {
				if err = d.Write(w); err != nil {
					return err
				}
			}
		}
	}

	if _, err = sb.WriteString(openingAngleBracket); err != nil {
		return err
	}
	if _, err = sb.WriteString(forwardSlash); err != nil {
		return err
	}
	if _, err = sb.WriteString(e.tagName); err != nil {
		return err
	}
	if _, err = sb.WriteString(closingAngleBracket); err != nil {
		return err
	}

	_, err = io.Copy(w, strings.NewReader(sb.String()))

	return err
}

func CreateElement(options ...string) Element {

	var tagName string
	var textContent string

	if len(options) > 0 {
		tagName = options[0]
	}
	if len(options) > 1 {
		textContent = options[1]
	}

	return &element{
		tagName:     tagName,
		textContent: textContent,
		classes:     make(map[string]any),
		attributes:  make(map[string]string),
		mtx:         &sync.Mutex{},
	}
}
