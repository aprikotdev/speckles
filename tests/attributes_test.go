package tests

import (
	"testing"

	. "github.com/aprikotdev/speckles/elements"
	"github.com/stretchr/testify/assert"
)

func TestBoolAttributes(t *testing.T) {
	run(t, []result{
		{
			Expected: "<input disabled>",
			Actual:   Input().Disabled(),
		},
		{
			Expected: "<button></button>",
			Actual:   Button().IfDisabled(false),
		},
		{
			Expected: "<select multiple></select>",
			Actual:   Select().Multiple(),
		},
		{
			Expected: "<option selected></option>",
			Actual:   Option().Selected(),
		},
		{
			Expected: "<textarea readonly></textarea>",
			Actual:   Textarea().Readonly(),
		},
		{
			Expected: "<form novalidate></form>",
			Actual:   Form().Novalidate(),
		},
		{
			Expected: "<iframe allowfullscreen></iframe>",
			Actual:   Iframe().Allowfullscreen(),
		},
		{
			Expected: "<fieldset disabled></fieldset>",
			Actual:   Fieldset().Disabled(),
		},
		{
			Expected: "<form disabled></form>",
			Actual:   Form().BoolAttr("disabled"),
		},
		{
			Expected: "<input>",
			Actual:   Input().IfBoolAttr(false, "disabled"),
		},
		{
			Expected: "<video autoplay muted></video>",
			Actual:   Video().Autoplay().Muted(),
		},
		{
			Expected: "<video autoplay muted></video>",
			Actual:   Video().Autoplay().MutedSet(true),
		},
		{
			Expected: "<video autoplay></video>",
			Actual:   Video().Autoplay().MutedSet(false),
		},
		{
			Expected: "<video autoplay></video>",
			Actual:   Video().Autoplay().Muted().MutedRemove(),
		},
	})
}

func TestStringAttributes(t *testing.T) {
	run(t, []result{
		{
			Expected: "<button popovertarget=\"my-popover\">Open Popover</button>",
			Actual:   Button().Popovertarget("my-popover").Text("Open Popover"),
		},
	})
}

func TestKVAttributes(t *testing.T) {
	run(t, []result{
		{
			Expected: "<div id=\"elt\" style=\"border-top:1px solid blue;color:red\">An example div</div>",
			Actual:   Div().ID("elt").Style("border-top: 1px solid blue; color: red;").Text("An example div"),
		},
		{
			Expected: "<div style=\"display:none\"></div>",
			Actual:   Div().StyleAdd("display", "none"),
		},
		{
			Expected: "<span style=\"color:red;display:block\"></span>",
			Actual: Span().
				StyleAdd("color", "red").
				StyleMap(map[string]string{"display": "block", "font-size": "12px", "font-weight": "bold"}).
				StyleRemove("font-size", "font-weight"),
		},
		{
			Expected: "<p style=\"display:block;margin:10px;padding:5px\"></p>",
			Actual:   P().StyleMap(map[string]string{"margin": "10px", "padding": "5px", "display": "block"}),
		},
	})

	assert.NotPanics(t, func() { P().StyleMap(map[string]string{"foo": ""}) })

	assert.Panics(t, func() { A().StylePairs("foo") })
	assert.Panics(t, func() { Div().StyleAdd("", "bar") })
	assert.Panics(t, func() { Span().Style(";;;;;;") })
	assert.Panics(t, func() { Div().Style("font-size; color: red;") })
}

func TestChoiceAttributes(t *testing.T) {
	run(t, []result{
		{
			Expected: "<div id=\"my-popover\" popover=\"auto\">Greetings, one and all!</div>",
			Actual:   Div().Popover(DivPopoverAuto).ID("my-popover").Text("Greetings, one and all!"),
		},
		{
			Expected: "<div popover></div>",
			Actual:   Div().Popover(DivPopoverEmpty),
		},
		{
			Expected: "<a hidden></a>",
			Actual:   A().Hidden(AHiddenEmpty),
		},
		{
			Expected: "<a hidden=\"until-found\"></a>",
			Actual:   A().Hidden(AHiddenUntilFound),
		},
	})
}

func TestSpaceDelimitedAttributes(t *testing.T) {
	run(t, []result{
		{
			Expected: "<div class=\"foo bar baz\"></div>",
			Actual:   Div().Class("foo bar baz hello").ClassRemove("hello"),
		},
		{
			Expected: "<div class=\"foo bar\"></div>",
			Actual:   Div().Class("foo").Class("bar"),
		},
	})
}

func TestCommaDelimitedAttributes(t *testing.T) {
	run(t, []result{
		{
			Expected: "<area alt=\"Html\" coords=\"260,96,209,249,130,138\" href=\"https://developer.mozilla.org/docs/Web/Html\" shape=\"poly\">",
			Actual:   Area().Shape(AreaShapePoly).Coords("260,96,209,249,130,138").Href("https://developer.mozilla.org/docs/Web/Html").Alt("Html"),
		},
		{
			Expected: "<area alt=\"Html\" coords=\"260,96,209,249,130\" href=\"https://developer.mozilla.org/docs/Web/Html\" shape=\"rect\">",
			Actual:   Area().Shape(AreaShapeRect).Coords("260,96,209,249,130,138").CoordsRemove("138").Href("https://developer.mozilla.org/docs/Web/Html").Alt("Html"),
		},
	})
}
