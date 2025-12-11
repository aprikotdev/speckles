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
	Tag              []byte
	IsSelfClosing    bool
	IntAttributes    *treemap.TreeMap[string, int]
	FloatAttributes  *treemap.TreeMap[string, float64]
	StringAttributes *treemap.TreeMap[string, string]
	DelimitedStrings *treemap.TreeMap[string, *DelimitedBuilder[string]]
	KVStrings        *treemap.TreeMap[string, *KVBuilder]
	BoolAttributes   *treemap.TreeMap[string, bool]
	Descendants      []ElementRenderer
}

func (e *Element) Attr(name string, value string) *Element {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set(name, value)
	return e
}

func (e *Element) Attrs(attrs ...string) *Element {
	if len(attrs)%2 != 0 {
		panic("attrs must be a multiple of 2")
	}
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	for i := 0; i < len(attrs); i += 2 {
		k := attrs[i]
		v := attrs[i+1]
		e.StringAttributes.Set(k, v)
	}
	return e
}

func (e *Element) AttrsMap(attrs map[string]string) *Element {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	for k, v := range attrs {
		e.StringAttributes.Set(k, v)
	}
	return e
}

func (e *Element) Render(w io.Writer) error {
	w.Write(openBracket)
	w.Write(e.Tag)

	finalKeys := treemap.New[string, string]()

	if e.IntAttributes != nil && e.IntAttributes.Len() > 0 {
		for it := e.IntAttributes.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			finalKeys.Set(k, fmt.Sprint(v))
		}
	}

	if e.FloatAttributes != nil && e.FloatAttributes.Len() > 0 {
		for it := e.FloatAttributes.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			finalKeys.Set(k, fmt.Sprint(v))
		}
	}

	if e.StringAttributes != nil && e.StringAttributes.Len() > 0 {
		for it := e.StringAttributes.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			finalKeys.Set(k, v)
		}
	}

	if e.DelimitedStrings != nil && e.DelimitedStrings.Len() > 0 {
		for it := e.DelimitedStrings.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			buf := bytebufferpool.Get()
			if err := v.Render(buf); err != nil {
				return err
			}
			finalKeys.Set(k, buf.String())
			bytebufferpool.Put(buf)
		}
	}

	if e.KVStrings != nil && e.KVStrings.Len() > 0 {
		for it := e.KVStrings.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			buf := bytebufferpool.Get()
			if err := v.Render(buf); err != nil {
				return err
			}
			finalKeys.Set(k, buf.String())
			bytebufferpool.Put(buf)
		}
	}

	if e.BoolAttributes != nil && e.BoolAttributes.Len() > 0 {
		for it := e.BoolAttributes.Iterator(); it.Valid(); it.Next() {
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

	if e.IsSelfClosing {
		w.Write(closeBracket)
		return nil
	}

	w.Write(closeBracket)

	for _, d := range e.Descendants {
		if d == nil {
			continue
		}
		if err := d.Render(w); err != nil {
			return err
		}
	}

	w.Write(openBracket)
	w.Write(slash)
	w.Write(e.Tag)
	w.Write(closeBracket)

	return nil
}

type DelimitedBuilder[T comparable] struct {
	Delimiter string
	Values    []T
}

func NewDelimitedBuilder[T comparable](delimiter string) *DelimitedBuilder[T] {
	return &DelimitedBuilder[T]{
		Delimiter: delimiter,
	}
}

func (d *DelimitedBuilder[T]) Add(values ...T) *DelimitedBuilder[T] {
	d.Values = append(d.Values, values...)
	return d
}

func (d *DelimitedBuilder[T]) Remove(values ...T) *DelimitedBuilder[T] {
	toRemove := make(map[T]struct{}, len(values))
	for _, v := range values {
		toRemove[v] = struct{}{}
	}

	n := 0
	for _, v := range d.Values {
		if _, found := toRemove[v]; !found {
			d.Values[n] = v
			n++
		}
	}
	d.Values = d.Values[:n]

	return d
}

func (d *DelimitedBuilder[T]) Render(w io.Writer) error {
	for i, v := range d.Values {
		b := fmt.Append(nil, v)
		if _, err := w.Write(b); err != nil {
			return err
		}

		if i < len(d.Values)-1 {
			if _, err := w.Write([]byte(d.Delimiter)); err != nil {
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

type KVBuilder struct {
	KeyPairDelimiter string
	EntryDelimiter   string
	Values           []keyValue
	keysIdx          map[string]int
}

func NewKVBuilder(keyPairDelimiter, entryDelimiter string) *KVBuilder {
	return &KVBuilder{
		KeyPairDelimiter: keyPairDelimiter,
		EntryDelimiter:   entryDelimiter,
		keysIdx:          make(map[string]int),
	}
}

func (d *KVBuilder) Add(key, value string) *KVBuilder {
	if i, found := d.keysIdx[key]; found {
		d.Values[i].Value = value
	} else {
		d.Values = append(d.Values, keyValue{key, value})
		d.keysIdx[key] = len(d.Values) - 1
	}
	return d
}

func (d *KVBuilder) Remove(keys ...string) *KVBuilder {
	for _, key := range keys {
		if idx, found := d.keysIdx[key]; found {
			d.Values = append(d.Values[:idx], d.Values[idx+1:]...)
			delete(d.keysIdx, key)

			// Rebuild keysIdx map
			for j := idx; j < len(d.Values); j++ {
				d.keysIdx[d.Values[j].Key] = j
			}
		}
	}
	return d
}

func (d *KVBuilder) Render(w io.Writer) error {
	for i, kv := range d.Values {
		if _, err := w.Write([]byte(kv.Key)); err != nil {
			return err
		}
		if _, err := w.Write([]byte(d.KeyPairDelimiter)); err != nil {
			return err
		}
		if _, err := w.Write([]byte(kv.Value)); err != nil {
			return err
		}

		if i < len(d.Values)-1 {
			if _, err := w.Write([]byte(d.EntryDelimiter)); err != nil {
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

func TextF(format string, args ...any) *TextContent {
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

func EscapedF(format string, args ...any) *EscapedContent {
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

func If(condition bool, children ...ElementRenderer) ElementRenderer {
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

func Range[T any](values []T, cb func(T) ElementRenderer) ElementRenderer {
	children := make([]ElementRenderer, 0, len(values))
	for _, value := range values {
		children = append(children, cb(value))
	}
	return Group(children...)
}

func RangeI[T any](values []T, cb func(int, T) ElementRenderer) ElementRenderer {
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

func DynIf(condition bool, childrenFuncs ...ElementRendererFunc) ElementRenderer {
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
		Tag:         []byte(tag),
		Descendants: children,
	}
}

func Error(err error) ElementRenderer {
	return Text(err.Error())
}
