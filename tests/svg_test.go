package tests

import (
	"testing"

	. "github.com/aprikotdev/speckles/elements"
)

func TestSVGSVGElement(t *testing.T) {
	run(t, []result{
		{
			Expected: `<svg height="200" viewBox="0 0 200 200" width="200" xmlns="http://www.w3.org/2000/svg"><circle cx="100" cy="100" r="80"></circle></svg>`,
			Actual: func() *SVGSVGElement {
				return SVGSVG().
					Attr("xmlns", "http://www.w3.org/2000/svg").
					Width("200").
					Height("200").
					ViewBox("0 0 200 200").
					Children(
						SVGCircle().Cx(100).Cy(100).R(80),
					)
			}(),
		},
	})
}

func TestSVGLinearGradientElement(t *testing.T) {
	run(t, []result{
		{
			Expected: `<linearGradient gradientTransform="skewX(20) translate(185, 0)" gradientUnits="objectBoundingBox" id="linear-gradient" x1="0.048" x2="0.963" y1="0.5" y2="0.5"><stop offset="0" stop-color="#000000"></stop><stop offset="1" stop-color="#0E67B4"></stop></linearGradient>`,
			Actual: SVGLinearGradient(
				SVGStop().Offset(0).StopColor("#000000"),
				SVGStop().Offset(1).StopColor("#0E67B4"),
			).
				ID("linear-gradient").
				GradientUnits("objectBoundingBox").
				GradientTransform("skewX(20) translate(185, 0)").
				X1(0.048).Y1(0.5).
				X2(0.963).Y2(0.5),
		},
	})
}

func TestSVGClipPathElement(t *testing.T) {
	run(t, []result{
		{
			Expected: `<clipPath id="clip-path"><rect class="cls-1" height="300" id="Rectangle_73" width="300"></rect></clipPath>`,
			Actual: SVGClipPath().ID("clip-path").Children(
				SVGRect().Class("cls-1").ID("Rectangle_73").Width(300).Height(300),
			),
		},
	})
}
