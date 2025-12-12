package tests

import (
	"testing"

	. "github.com/aprikotdev/speckles/elements"
)

func TestDivElement(t *testing.T) {
	run(t, []result{
		{
			Expected: `<div></div>`,
			Actual:   Div(),
		},
		{
			Expected: `<div data-foo="bar"></div>`,
			Actual:   Div().Attr("data-foo", "bar"),
		},
		{
			Expected: `<div data-baz data-bind-foo="bar"></div>`,
			Actual:   Div().Attr("data-bind-foo", "bar").BoolAttr("data-baz"),
		},
	})
}

func TestNavElement(t *testing.T) {
	run(t, []result{
		{
			Expected: `<nav class="navbar"><ol><li><a href="/">Home</a></li><li><a href="/contact">Contact</a></li><li><a href="/about">About</a></li></ol></nav>`,
			Actual: Nav().Class("navbar").Children(
				Ol(
					Li(A().Href("/").Text("Home")),
					Li(A().Href("/contact").Text("Contact")),
					Li(A().Href("/about").Text("About")),
				),
			),
		},
		{
			Expected: `<nav><ul class="navigation"><li><a href="/home">Home</a></li><li><a href="/about">About</a></li></ul></nav>`,
			Actual: func() ElementRenderer {
				type Navigation struct {
					Link string
					Item string
				}
				nav := []*Navigation{
					{Link: "/home", Item: "Home"},
					{Link: "/about", Item: "About"},
				}

				return Nav().Children(
					Ul().Class("navigation").Children(
						Range(nav, func(n *Navigation) ElementRenderer {
							return Li().Children(
								A().Href(n.Link).Text(n.Item),
							)
						}),
					),
				)
			}(),
		},
	})
}

func TestHeaderElement(t *testing.T) {
	run(t, []result{
		{
			Expected: `<header><title>Alice's Home Page</title><div class="header">Page Header</div></header>`,
			Actual: Header().Children(
				Title().TextF("%s's Home Page", "Alice"),
				Div().Class("header").Text("Page Header"),
			),
		},
	})
}

func TestFooterElement(t *testing.T) {
	run(t, []result{
		{
			Expected: `<footer><div class="footer">copyright 2016</div></footer>`,
			Actual:   Footer(Div().Class("footer").Text("copyright 2016")),
		},
	})
}

func TestSectionElement(t *testing.T) {
	run(t, []result{
		{
			Expected: `<section><div class="content"><div class="welcome"><h4>Hello John</h4><div class="raw"><a href="http://john.com">This is some raw content</a></div><div class="enc">&amp;lt;a href=&amp;#34;http://john.com&amp;#34;&amp;gt;This is some encoded content&amp;lt;/a&amp;gt;</div></div><p>John has 1 message</p><p>John has 2 messages</p><p>John has 3 messages</p><p>John has 4 messages</p><p>John has 5 messages</p></div></section>`,
			Actual: Section(
				Div().Class("content").Children(
					Div().Class("welcome").Children(
						H4().TextF("Hello %s", "John"),
						Div().Class("raw").Text(`<a href="http://john.com">This is some raw content</a>`),
						Div().Class("enc").Escaped(`&lt;a href=&#34;http://john.com&#34;&gt;This is some encoded content&lt;/a&gt;`),
					),
					Range([]int{0, 1, 2, 3, 4}, func(i int) ElementRenderer {
						count := i + 1
						return Tern(
							count == 1,
							P().TextF("%s has %d message", "John", count),
							P().TextF("%s has %d messages", "John", count),
						)
					}),
				),
			),
		},
	})
}

func TestOlElement(t *testing.T) {
	run(t, []result{
		{
			Expected: `<ol reversed start="5" type="a"><li>Item 5</li><li>Item 6</li><li>Item 7</li></ol>`,
			Actual: Ol().Reversed().Start(5).Type(OlTypeLowerA).Children(
				Range([]int{5, 6, 7}, func(i int) ElementRenderer {
					return Li().TextF("Item %d", i)
				}),
			),
		},
		{
			Expected: `<ol type="a"><li value="5">Item 5</li></ol>`,
			Actual: Ol().Type(OlTypeLowerA).Children(
				Li().Value(5).Text("Item 5"),
			),
		},
		{
			Expected: `<ol type="1"></ol>`,
			Actual:   Ol().Type(OlType1),
		},
		{
			Expected: `<ol type="I"></ol>`,
			Actual:   Ol().Type(OlTypeUpperI),
		},
	})
}

func TestHTMLElement(t *testing.T) {
	run(t, []result{
		{
			Expected: `<html><body><div class="header">Page Header</div><div autocapitalize="off" class="block bg-red-200 text-red-600 italic" style="font-size:12px">bar</div></body></html>`,
			Actual: HTML(
				Body(
					Div().Class("header").Text("Page Header"),
					Div().
						StyleAdd("color", "rad").
						StyleAdd("font-size", "12px").
						StyleRemove("color").
						Class("block bg-red-200 hidden").
						Class("text-red-600 italic").
						ClassRemove("hidden").
						Autocapitalize(DivAutocapitalizeOff).
						Text("bar"),
				),
			),
		},
	})
}

func TestGrouper(t *testing.T) {
	type User struct {
		FirstName      string
		Email          string
		FavoriteColors []string
		RawContent     string
		EscapedContent string
	}

	type Navigation struct {
		Item string
		Link string
	}

	user := &User{
		FirstName:      "Bob",
		FavoriteColors: []string{"blue", "green", "mauve"},
		RawContent:     "<div><p>Raw Content to be displayed</p></div>",
		EscapedContent: "<div><div><div>Escaped</div></div></div>",
	}

	nav := []*Navigation{
		{
			Item: "Link 1",
			Link: "http://www.mytest.com/",
		}, {
			Item: "Link 2",
			Link: "http://www.mytest.com/",
		}, {
			Item: "Link 3",
			Link: "http://www.mytest.com/",
		},
	}

	header := func(title string) ElementRenderer {
		return Header(
			Title().TextF("%s's Home Page", title),
			Div().Class("header").Text("Page Header"),
		)
	}

	navigation := func(nav []*Navigation) ElementRenderer {
		return Nav(
			Ul(
				Range(nav, func(n *Navigation) ElementRenderer {
					return Li(
						A().Href(n.Link).Text(n.Item),
					)
				}),
			).Class("navigation"),
		)
	}

	footer := func() ElementRenderer {
		return Footer(Div().Class("footer").Text("copyright 2016"))
	}

	index := func(u *User, nav []*Navigation, title string) ElementRenderer {
		return Group(
			Text("<!Doctype html>"),
			HTML(
				Body(
					header(title),
					navigation(nav),
					Section(
						Div(
							Div(
								H4().TextF("Hello %s", u.FirstName),
								Div().Class("raw").Text(u.RawContent),
								Div().Class("enc").Escaped(u.EscapedContent),
							).Class("welcome"),
							Range([]int{0, 1, 2, 3, 4}, func(i int) ElementRenderer {
								count := i + 1
								return Tern(
									count == 1,
									P().TextF("%s has %d message", u.FirstName, count),
									P().TextF("%s has %d messages", u.FirstName, count),
								)
							}),
						).Class("content"),
					),
					footer(),
				),
			),
		)
	}
	run(t, []result{
		{
			Expected: `<!Doctype html><html><body><header><title>Bob's Home Page</title><div class="header">Page Header</div></header><nav><ul class="navigation"><li><a href="http://www.mytest.com/">Link 1</a></li><li><a href="http://www.mytest.com/">Link 2</a></li><li><a href="http://www.mytest.com/">Link 3</a></li></ul></nav><section><div class="content"><div class="welcome"><h4>Hello Bob</h4><div class="raw"><div><p>Raw Content to be displayed</p></div></div><div class="enc">&lt;div&gt;&lt;div&gt;&lt;div&gt;Escaped&lt;/div&gt;&lt;/div&gt;&lt;/div&gt;</div></div><p>Bob has 1 message</p><p>Bob has 2 messages</p><p>Bob has 3 messages</p><p>Bob has 4 messages</p><p>Bob has 5 messages</p></div></section><footer><div class="footer">copyright 2016</div></footer></body></html>`,
			Actual:   index(user, nav, user.FirstName),
		},
	})
}
