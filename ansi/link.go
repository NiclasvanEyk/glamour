package ansi

import (
	"io"
	"net/url"

	"github.com/muesli/termenv"
)

// A LinkElement is used to render hyperlinks.
type LinkElement struct {
	Text    string
	BaseURL string
	URL     string
	Child   ElementRenderer
}

func (e *LinkElement) Render(w io.Writer, ctx RenderContext) error {
	var textRendered bool
	if len(e.Text) > 0 && e.Text != e.URL {
		textRendered = true

		el := &BaseElement{
			Token: e.Text,
			Style: ctx.options.Styles.LinkText,
		}

		// TODO: There is some logic at the end of this function that resolves
		//       URLs and whether to render them at all. This should be re-used
		//       here.
		hasLinkThatShouldBeRendered := true
		if ctx.options.OmitLinkUrls && hasLinkThatShouldBeRendered {
			// See https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda#the-escape-sequence
			// TODO: Replace test link to google.com
			el.Prefix = termenv.OSC + "8;;" + "https://google.com" + termenv.ST
			el.Suffix = termenv.OSC + "8;;" + termenv.ST
		}

		err := el.Render(w, ctx)
		if err != nil {
			return err
		}
	}

	if ctx.options.OmitLinkUrls {
		return nil
	}

	/*
		if node.LastChild != nil {
			if node.LastChild.Type == bf.Image {
				el := tr.NewElement(node.LastChild)
				err := el.Renderer.Render(w, node.LastChild, tr)
				if err != nil {
					return err
				}
			}
			if len(node.LastChild.Literal) > 0 &&
				string(node.LastChild.Literal) != string(node.LinkData.Destination) {
				textRendered = true
				el := &BaseElement{
					Token: string(node.LastChild.Literal),
					Style: ctx.style[LinkText],
				}
				err := el.Render(w, node.LastChild, tr)
				if err != nil {
					return err
				}
			}
		}
	*/

	u, err := url.Parse(e.URL)
	if err == nil &&
		"#"+u.Fragment != e.URL { // if the URL only consists of an anchor, ignore it
		pre := " "
		style := ctx.options.Styles.Link
		if !textRendered {
			pre = ""
			style.BlockPrefix = ""
			style.BlockSuffix = ""
		}

		el := &BaseElement{
			Token:  resolveRelativeURL(e.BaseURL, e.URL),
			Prefix: pre,
			Style:  style,
		}
		err := el.Render(w, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
