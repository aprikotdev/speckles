package config

import (
	"strings"
	"unicode"

	"github.com/aprikotdev/speckles/internal/caser"
)

type Namespace struct {
	Name        string
	Description string
	Prefix      string
	Elements    []*Element
	Attributes  []*Attribute
}

type Element struct {
	NoChildren  bool
	Name        string
	Tag         string
	Description string
	Attributes  []*Attribute
}

type Attribute struct {
	Key         string
	Name        string
	Description string
	Type        AttributeType
}

type AttributeType any

type (
	attributeTypeBool      struct{}
	attributeTypeRune      struct{}
	attributeTypeInt       struct{}
	attributeTypeNumber    struct{}
	attributeTypeString    struct{}
	attributeTypeDelimited struct {
		Delimiter string
	}
	attributeTypeKeyValue struct {
		KeyValueDelimiter string
		PairDelimiter     string
	}
	attributeTypeChoice struct {
		Name        string
		Description string
	}
	attributeTypeChoices struct {
		Choices []*attributeTypeChoice
	}
)

func AttributeTypeBool() *attributeTypeBool {
	return &attributeTypeBool{}
}

func IsAttributeTypeBool(attributeType AttributeType) bool {
	_, ok := attributeType.(*attributeTypeBool)
	return ok
}

func AttributeTypeRune() *attributeTypeRune {
	return &attributeTypeRune{}
}

func IsAttributeTypeRune(attributeType AttributeType) bool {
	_, ok := attributeType.(*attributeTypeRune)
	return ok
}

func AttributeTypeInt() *attributeTypeInt {
	return &attributeTypeInt{}
}

func IsAttributeTypeInt(attributeType AttributeType) bool {
	_, ok := attributeType.(*attributeTypeInt)
	return ok
}

func AttributeTypeNumber() *attributeTypeNumber {
	return &attributeTypeNumber{}
}

func IsAttributeTypeNumber(attributeType AttributeType) bool {
	_, ok := attributeType.(*attributeTypeNumber)
	return ok
}

func AttributeTypeString() *attributeTypeString {
	return &attributeTypeString{}
}

func IsAttributeTypeString(attributeType AttributeType) bool {
	_, ok := attributeType.(*attributeTypeString)
	return ok
}

func IsAttributeTypeDelimited(attributeType AttributeType) bool {
	_, ok := attributeType.(*attributeTypeDelimited)
	return ok
}

func AttributeTypeDelimited(delimiter string) *attributeTypeDelimited {
	return &attributeTypeDelimited{Delimiter: delimiter}
}

func AttributeTypeSpaceDelimited() *attributeTypeDelimited {
	return AttributeTypeDelimited(" ")
}

func AttributeTypeCommaDelimited() *attributeTypeDelimited {
	return AttributeTypeDelimited(",")
}

func AttributeTypeChoice(name, description string) *attributeTypeChoice {
	return &attributeTypeChoice{Name: name, Description: description}
}

func AttributeTypeChoices(choices ...*attributeTypeChoice) *attributeTypeChoices {
	return &attributeTypeChoices{Choices: choices}
}

func IsAttributeTypeChoices(attributeType AttributeType) bool {
	_, ok := attributeType.(*attributeTypeChoices)
	return ok
}

func AttributeTypeKeyValue(keyValueDelimiter, pairDelimiter string) *attributeTypeKeyValue {
	return &attributeTypeKeyValue{
		KeyValueDelimiter: keyValueDelimiter,
		PairDelimiter:     pairDelimiter,
	}
}

func IsAttributeTypeKeyValue(attributeType AttributeType) bool {
	_, ok := attributeType.(*attributeTypeKeyValue)
	return ok
}

func AttributeTypeKeyValueColonSemicolon() *attributeTypeKeyValue {
	return AttributeTypeKeyValue(":", ";")
}

func ChoiceSuffix(choiceName string, choices []*attributeTypeChoice) string {
	if choiceName == "" {
		return "Empty"
	}

	// Handle single-character choice names that may conflict in casing
	if len(choiceName) == 1 {
		for _, choice := range choices {
			if choiceName != choice.Name {
				if strings.EqualFold(choiceName, choice.Name) {
					char := rune(choiceName[0])
					if unicode.IsUpper(char) {
						choiceName = "_upper_" + choiceName
					} else {
						choiceName = "_lower_" + choiceName
					}
					break
				}
			}
		}
	}

	return caser.GoPascal(choiceName)
}

func Namespaces() []*Namespace {
	return []*Namespace{HTML, SVG, MathML}
}
