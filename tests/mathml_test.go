package tests

import (
	"testing"

	. "github.com/aprikotdev/speckles/elements"
)

func TestMathMLElement(t *testing.T) {
	run(t, []result{
		{
			Expected: `<math><mfrac><mn>1</mn><mn>3</mn></mfrac></math>`,
			Actual:   MathMLMath().Children(MathMLMfrac().Children(MathMLMn().Text("1"), MathMLMn().Text("3"))),
		},
	})
}

func TestMathMLAnnotationElement(t *testing.T) {
	run(t, []result{
		{
			Expected: `<annotation encoding="application/x-tex"></annotation>`,
			Actual:   MathMLAnnotation().Encoding("application/x-tex"),
		},
	})
}

func TestMathMLAnnotationXMLElement(t *testing.T) {
	run(t, []result{
		{
			Expected: `<annotation-xml encoding="application/mathml-presentation+xml"></annotation-xml>`,
			Actual:   MathMLAnnotationXML().Encoding("application/mathml-presentation+xml"),
		},
	})
}
