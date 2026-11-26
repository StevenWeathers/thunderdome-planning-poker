// Package sanitize holds the HTML sanitization policies for user supplied content
package sanitize

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

// quillClasses matches the classes the UI's rich text editor uses for block formatting
var quillClasses = regexp.MustCompile(`^ql-(align-(center|right|justify)|indent-[1-8]|direction-rtl|syntax)` +
	`( ql-(align-(center|right|justify)|indent-[1-8]|direction-rtl|syntax))*$`)

// RichTextPolicy is the user generated content policy, plus the formatting
// produced by the UI's rich text editor that the base policy would strip
func RichTextPolicy() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()

	p.AllowAttrs("class").Matching(quillClasses).
		OnElements("p", "h1", "h2", "h3", "h4", "h5", "h6", "li", "blockquote", "pre")
	p.AllowStyles("color", "background-color").
		OnElements("span", "strong", "em", "u", "s", "sub", "sup", "a")
	// images pasted into the editor are stored inline
	p.AllowDataURIImages()

	return p
}
