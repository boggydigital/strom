package strom

import (
	_ "embed"
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
	colon               = ":"
	semicolon           = ";"
	singleSpace         = " "

	commentPrefix     = "<!--"
	commentSuffix     = "-->"
	docTypeHtmlPrefix = "<!doctype html>"

	classAttributeName = "class"
	styleAttributeName = "style"
	idAttributeName    = "id"

	htmlTagName  = "html"
	styleTagName = "style"
)

//go:embed "styles/colors.css"
var colorStylesheet []byte

//go:embed "styles/units.css"
var unitsStylesheet []byte

//go:embed "styles/page.css"
var pageStylesheet []byte

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
	HasClass(classes ...string) bool
	GetTagName() string
	GetAttribute(name string) string
	SetAttribute(name string, value string) Element
	SetStyle(styles map[string]string) Element
	Append(nodes ...Element) Element
	GetElementById(id string) Element
	GetElementsByTagName(tagName string) iter.Seq[Element]
	GetElementsByClassName(classes ...string) iter.Seq[Element]
	Write(w io.Writer) error
}

func (en *elementNode) SetTextContent(textContent string) Element {
	en.mtx.Lock()
	defer en.mtx.Unlock()

	en.textContent = []byte(textContent)
	return en
}

func (en *elementNode) AddClass(classes ...string) Element {
	en.mtx.Lock()
	defer en.mtx.Unlock()

	if en.classes == nil {
		en.classes = make(map[string]any)
	}

	for _, class := range classes {
		en.classes[class] = nil
	}

	return en
}

func (en *elementNode) HasClass(classes ...string) bool {
	for _, class := range classes {
		if _, ok := en.classes[class]; !ok {
			return false
		}
	}
	return true
}

func (en *elementNode) GetTagName() string {
	return en.tagName
}

func (en *elementNode) SetStyle(styles map[string]string) Element {
	en.mtx.Lock()
	defer en.mtx.Unlock()

	if en.styles == nil {
		en.styles = make(map[string]string)
	}

	for p, v := range styles {
		en.styles[p] = v
	}

	return en
}

func (en *elementNode) GetAttribute(name string) string {
	return en.attributes[name]
}

func (en *elementNode) GetElementById(id string) Element {
	for _, child := range en.children {
		if cid := child.GetAttribute(idAttributeName); cid == id {
			return child
		}
		if el := child.GetElementById(id); el != nil {
			return el
		}
	}
	return nil
}

func (en *elementNode) GetElementsByTagName(tagName string) iter.Seq[Element] {
	return func(yield func(Element) bool) {
		for _, child := range en.children {
			if child.GetTagName() == tagName {
				if !yield(child) {
					return
				}
			}
			for match := range child.GetElementsByTagName(tagName) {
				if !yield(match) {
					return
				}
			}
		}
	}
}

func (en *elementNode) GetElementsByClassName(classes ...string) iter.Seq[Element] {
	return func(yield func(element Element) bool) {
		for _, child := range en.children {
			if child.HasClass(classes...) {
				if !yield(child) {
					return
				}
			}
			for match := range child.GetElementsByClassName(classes...) {
				if !yield(match) {
					return
				}
			}
		}
	}
}

func (en *elementNode) classList() string {
	return strings.Join(slices.Collect(maps.Keys(en.classes)), singleSpace)
}

func (en *elementNode) SetAttribute(name string, value string) Element {
	en.mtx.Lock()
	defer en.mtx.Unlock()

	if en.attributes == nil {
		en.attributes = make(map[string]string)
	}

	en.attributes[name] = value
	return en
}

func (en *elementNode) Append(nodes ...Element) Element {
	en.children = append(en.children, nodes...)
	return en
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

func (en *elementNode) Write(w io.Writer) error {

	en.mtx.Lock()
	defer en.mtx.Unlock()

	// writing element node is a sequence of writing:
	// 1. prefix (e.g. `<!--` for comments)
	// 2. `<tagName`
	// 3. classes (` class='class1 class2`)
	// 4. inline styles (` style='color:purple'`)
	// 5. attributes (` attr='value')
	// 6. `>`
	// 7. text content
	// 8. children or deferred children
	// 9. `</tagName>`
	// 10. suffix (e.g. `-->` for comments)

	// 1
	switch en.prefix {
	case nil:
	// do nothing
	default:
		if _, err := w.Write(en.prefix); err != nil {
			return err
		}

	}

	switch en.tagName {
	case "":
		if len(en.classes) > 0 || len(en.styles) > 0 || len(en.attributes) > 0 || len(en.textContent) > 0 {
			return errors.New("transparent container doesn't support classes, styles, attributes of textContent")
		}
	// do nothing
	default:
		// 2
		if err := writeStrings(w, openingAngleBracket, en.tagName); err != nil {
			return nil
		}

		// 3
		if len(en.classes) > 0 {

			for _, ra := range restrictedAttributes {
				if _, ok := en.attributes[ra]; ok {
					return errors.New("restrictred attribute " + ra)
				}
			}

			if err := writeAttribute(w, classAttributeName, en.classList()); err != nil {
				return err
			}
		}

		// 4
		if len(en.styles) > 0 {
			if err := writeStyles(w, en.styles); err != nil {
				return err
			}
		}

		// 5
		for attributeName, attributeValue := range en.attributes {
			if err := writeAttribute(w, attributeName, attributeValue); err != nil {
				return err
			}
		}

		// 6
		if err := writeStrings(w, closingAngleBracket); err != nil {
			return err
		}

		// 7
		if _, err := w.Write(en.textContent); err != nil {
			return err
		}
	}

	// 8
	var children iter.Seq[Element]

	switch en.getChildrenDelegate {
	case nil:
		children = slices.Values(en.children)
	default:
		children = en.getChildrenDelegate()
	}

	if children != nil {
		for child := range children {
			if err := child.Write(w); err != nil {
				return err
			}
		}
	}

	// 9
	switch en.tagName {
	case "":
	// do nothing
	default:
		if err := writeStrings(w, openingAngleBracket, forwardSlash, en.tagName, closingAngleBracket); err != nil {
			return err
		}
	}

	// 10
	switch en.suffix {
	case nil:
	// do nothing
	default:
		if _, err := w.Write(en.suffix); err != nil {
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
		prefix:  []byte(commentPrefix),
		tagName: tagName,
		suffix:  []byte(commentSuffix),
	}
}

func Defer(getChildredDelegate func() iter.Seq[Element]) Element {
	return &elementNode{
		getChildrenDelegate: getChildredDelegate,
		mtx:                 &sync.Mutex{},
	}
}

func DoctypeHml() Element {
	return &elementNode{
		mtx:     new(sync.Mutex),
		tagName: htmlTagName,
		prefix:  []byte(docTypeHtmlPrefix),
	}
}

func Stylesheet(content []byte) Element {
	return &elementNode{
		mtx:         new(sync.Mutex),
		tagName:     styleTagName,
		textContent: content,
	}
}

func Page(title string) Element {
	root := DoctypeHml().
		SetAttribute("id", "_top").
		SetAttribute("lang", "en")

	head := Create("head")
	head.Append(CreateText("title", title))
	head.Append(Defer(headDeferrals))
	root.Append(head)

	body := Create("body")
	root.Append(body)

	return root
}

func headDeferrals() iter.Seq[Element] {
	return func(yield func(element Element) bool) {
		if !yield(Create("meta").
			SetAttribute("charset", "utf-8")) {
			return
		}
		if !yield(Create("meta").
			SetAttribute("name", "viewport").
			SetAttribute("content", "width=device-width,initial-scale=1.0")) {
			return
		}
		if !yield(Create("meta").
			SetAttribute("color-scheme", "light dark")) {
			return
		}
		if !yield(Create("meta").
			SetAttribute("name", "format-detection").
			SetAttribute("content", "telephone=no")) {
			return
		}

		if !yield(Stylesheet(colorStylesheet)) {
			return
		}
		if !yield(Stylesheet(unitsStylesheet)) {
			return
		}
		if !yield(Stylesheet(pageStylesheet)) {
			return
		}
	}
}
