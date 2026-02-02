package elements

import (
	"fmt"
	"html"
	"io"

	"github.com/igrmk/treemap/v2"
	"github.com/valyala/bytebufferpool"
)

var (
	openBracket  = []byte("<")
	closeBracket = []byte(">")
	slash        = []byte("/")
	equal        = []byte("=")
	doubleQuotes = []byte("\"")
	space        = []byte(" ")
)

type ElementRenderer interface {
	Render(w io.Writer) error
}

type ElementRendererFunc func() ElementRenderer

type Element struct {
	tag              []byte
	isSelfClosing    bool
	intAttributes    *treemap.TreeMap[string, int]
	floatAttributes  *treemap.TreeMap[string, float64]
	stringAttributes *treemap.TreeMap[string, string]
	delimitedStrings *treemap.TreeMap[string, *delimitedBuilder[string]]
	keyValueStrings  *treemap.TreeMap[string, *keyValueBuilder]
	boolAttributes   *treemap.TreeMap[string, bool]
	descendants      []ElementRenderer
}

func (e *Element) Attr(name string, value string) *Element {
	if e.stringAttributes == nil {
		e.stringAttributes = treemap.New[string, string]()
	}
	e.stringAttributes.Set(name, value)
	return e
}

func (e *Element) Attrs(attrs ...string) *Element {
	if len(attrs)%2 != 0 {
		panic("attrs must be a multiple of 2")
	}
	if e.stringAttributes == nil {
		e.stringAttributes = treemap.New[string, string]()
	}
	for i := 0; i < len(attrs); i += 2 {
		k := attrs[i]
		v := attrs[i+1]
		e.stringAttributes.Set(k, v)
	}
	return e
}

func (e *Element) AttrsMap(attrs map[string]string) *Element {
	if e.stringAttributes == nil {
		e.stringAttributes = treemap.New[string, string]()
	}
	for k, v := range attrs {
		e.stringAttributes.Set(k, v)
	}
	return e
}

func (e *Element) Render(w io.Writer) error {
	w.Write(openBracket)
	w.Write(e.tag)

	finalKeys := treemap.New[string, string]()

	if e.intAttributes != nil && e.intAttributes.Len() > 0 {
		for it := e.intAttributes.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			finalKeys.Set(k, fmt.Sprint(v))
		}
	}

	if e.floatAttributes != nil && e.floatAttributes.Len() > 0 {
		for it := e.floatAttributes.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			finalKeys.Set(k, fmt.Sprint(v))
		}
	}

	if e.stringAttributes != nil && e.stringAttributes.Len() > 0 {
		for it := e.stringAttributes.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			finalKeys.Set(k, v)
		}
	}

	if e.delimitedStrings != nil && e.delimitedStrings.Len() > 0 {
		for it := e.delimitedStrings.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			buf := bytebufferpool.Get()
			if err := v.Render(buf); err != nil {
				return err
			}
			finalKeys.Set(k, buf.String())
			bytebufferpool.Put(buf)
		}
	}

	if e.keyValueStrings != nil && e.keyValueStrings.Len() > 0 {
		for it := e.keyValueStrings.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			buf := bytebufferpool.Get()
			if err := v.Render(buf); err != nil {
				return err
			}
			finalKeys.Set(k, buf.String())
			bytebufferpool.Put(buf)
		}
	}

	if e.boolAttributes != nil && e.boolAttributes.Len() > 0 {
		for it := e.boolAttributes.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			if v {
				finalKeys.Set(k, "")
			}
		}
	}

	if finalKeys.Len() > 0 {
		for it := finalKeys.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()

			w.Write(space)
			w.Write([]byte(k))

			if v != "" {
				w.Write(equal)
				w.Write(doubleQuotes)
				w.Write(fmt.Append(nil, v))
				w.Write(doubleQuotes)
			}
		}
	}

	if e.isSelfClosing {
		w.Write(closeBracket)
		return nil
	}

	w.Write(closeBracket)

	for _, d := range e.descendants {
		if d == nil {
			continue
		}
		if err := d.Render(w); err != nil {
			return err
		}
	}

	w.Write(openBracket)
	w.Write(slash)
	w.Write(e.tag)
	w.Write(closeBracket)

	return nil
}

type delimitedBuilder[T comparable] struct {
	delimiter string
	values    []T
}

func newDelimitedBuilder[T comparable](delimiter string) *delimitedBuilder[T] {
	return &delimitedBuilder[T]{
		delimiter: delimiter,
	}
}

func (d *delimitedBuilder[T]) Add(values ...T) *delimitedBuilder[T] {
	d.values = append(d.values, values...)
	return d
}

func (d *delimitedBuilder[T]) Remove(values ...T) *delimitedBuilder[T] {
	toRemove := make(map[T]struct{}, len(values))
	for _, v := range values {
		toRemove[v] = struct{}{}
	}

	n := 0
	for _, v := range d.values {
		if _, found := toRemove[v]; !found {
			d.values[n] = v
			n++
		}
	}
	d.values = d.values[:n]

	return d
}

func (d *delimitedBuilder[T]) Render(w io.Writer) error {
	for i, v := range d.values {
		b := fmt.Append(nil, v)
		if _, err := w.Write(b); err != nil {
			return err
		}

		if i < len(d.values)-1 {
			if _, err := w.Write([]byte(d.delimiter)); err != nil {
				return err
			}
		}
	}
	return nil
}

type keyValue struct {
	Key   string
	Value string
}

type keyValueBuilder struct {
	keyPairDelimiter string
	entryDelimiter   string
	values           []keyValue
	keysIdx          map[string]int
}

func newKVBuilder(keyPairDelimiter, entryDelimiter string) *keyValueBuilder {
	return &keyValueBuilder{
		keyPairDelimiter: keyPairDelimiter,
		entryDelimiter:   entryDelimiter,
		keysIdx:          make(map[string]int),
	}
}

func (d *keyValueBuilder) Add(key, value string) *keyValueBuilder {
	if i, found := d.keysIdx[key]; found {
		d.values[i].Value = value
	} else {
		d.values = append(d.values, keyValue{key, value})
		d.keysIdx[key] = len(d.values) - 1
	}
	return d
}

func (d *keyValueBuilder) Remove(keys ...string) *keyValueBuilder {
	for _, key := range keys {
		if idx, found := d.keysIdx[key]; found {
			d.values = append(d.values[:idx], d.values[idx+1:]...)
			delete(d.keysIdx, key)

			// Rebuild keysIdx map
			for j := idx; j < len(d.values); j++ {
				d.keysIdx[d.values[j].Key] = j
			}
		}
	}
	return d
}

func (d *keyValueBuilder) Render(w io.Writer) error {
	for i, kv := range d.values {
		if _, err := w.Write([]byte(kv.Key)); err != nil {
			return err
		}
		if _, err := w.Write([]byte(d.keyPairDelimiter)); err != nil {
			return err
		}
		if _, err := w.Write([]byte(kv.Value)); err != nil {
			return err
		}

		if i < len(d.values)-1 {
			if _, err := w.Write([]byte(d.entryDelimiter)); err != nil {
				return err
			}
		}
	}
	return nil
}

type TextContent string

func (tc *TextContent) Render(w io.Writer) error {
	_, err := w.Write([]byte(*tc))
	return err
}

func Text(text string) *TextContent {
	return (*TextContent)(&text)
}

func Textf(format string, args ...any) *TextContent {
	return Text(fmt.Sprintf(format, args...))
}

type EscapedContent string

func (ec *EscapedContent) Render(w io.Writer) error {
	_, err := w.Write([]byte(html.EscapeString(string(*ec))))
	return err
}

func Escaped(text string) *EscapedContent {
	return (*EscapedContent)(&text)
}

func Escapedf(format string, args ...any) *EscapedContent {
	return Escaped(fmt.Sprintf(format, args...))
}

type Grouper struct {
	Children []ElementRenderer
}

func (g *Grouper) Render(w io.Writer) error {
	for _, child := range g.Children {
		if err := child.Render(w); err != nil {
			return fmt.Errorf("failed to build element: %w", err)
		}
	}
	return nil
}

func Group(children ...ElementRenderer) *Grouper {
	return &Grouper{
		Children: children,
	}
}

func If(condition bool, children ...ElementRenderer) *Grouper {
	if condition {
		return Group(children...)
	}
	return nil
}

func Tern(condition bool, trueChildren, falseChildren ElementRenderer) ElementRenderer {
	if condition {
		return trueChildren
	}
	return falseChildren
}

func Range[T any](values []T, cb func(T) ElementRenderer) *Grouper {
	children := make([]ElementRenderer, 0, len(values))
	for _, value := range values {
		children = append(children, cb(value))
	}
	return Group(children...)
}

func RangeI[T any](values []T, cb func(int, T) ElementRenderer) *Grouper {
	children := make([]ElementRenderer, 0, len(values))
	for i, value := range values {
		children = append(children, cb(i, value))
	}
	return Group(children...)
}

func DynGroup(childrenFuncs ...ElementRendererFunc) *Grouper {
	children := make([]ElementRenderer, 0, len(childrenFuncs))
	for _, childFunc := range childrenFuncs {
		child := childFunc()
		if child != nil {
			children = append(children, child)
		}
	}
	return &Grouper{
		Children: children,
	}
}

func DynIf(condition bool, childrenFuncs ...ElementRendererFunc) *Grouper {
	if condition {
		children := make([]ElementRenderer, 0, len(childrenFuncs))
		for _, childFunc := range childrenFuncs {
			child := childFunc()
			if child != nil {
				children = append(children, child)
			}
		}
		return Group(children...)
	}
	return nil
}

func DynTern(condition bool, trueChildren, falseChildren ElementRendererFunc) ElementRenderer {
	if condition {
		return trueChildren()
	}
	return falseChildren()
}

func NewElement(tag string, children ...ElementRenderer) *Element {
	return &Element{
		tag:         []byte(tag),
		descendants: children,
	}
}

func Error(err error) ElementRenderer {
	return Text(err.Error())
}
