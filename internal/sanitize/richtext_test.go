package sanitize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRichTextPolicyStripsScript(t *testing.T) {
	p := RichTextPolicy()
	tests := map[string]string{
		`<p>ok</p><script>alert(1)</script>`:               `<p>ok</p>`,
		`<img src="x" onerror="alert(1)">`:                 `<img src="x">`,
		`<a href="javascript:alert(1)">x</a>`:              `x`,
		`<p class="fixed inset-0">x</p>`:                   `<p>x</p>`,
		`<span style="position: fixed">x</span>`:           `<span>x</span>`,
		`<span style="color: url(javascript:x)">x</span>`:  `<span>x</span>`,
		`<iframe src="https://example.com"></iframe>ok`:    `ok`,
		`<img src="data:text/html;base64,PHNjcmlwdD4=">ok`: `ok`,
	}

	for input, want := range tests {
		assert.Equal(t, want, p.Sanitize(input), input)
	}
}

func TestRichTextPolicyKeepsEditorFormatting(t *testing.T) {
	p := RichTextPolicy()
	tests := []string{
		`<p class="ql-align-center">x</p>`,
		`<p class="ql-indent-2 ql-align-right">x</p>`,
		`<pre class="ql-syntax">x</pre>`,
		`<ol><li class="ql-indent-1">x</li></ol>`,
		`<h2>x</h2><blockquote>x</blockquote>`,
		`<p><strong>b</strong><em>i</em><u>u</u><s>s</s><sub>1</sub><sup>2</sup></p>`,
		`<span style="color: rgb(230, 0, 0)">x</span>`,
		`<strong style="background-color: rgb(255, 255, 0)">x</strong>`,
		`<img src="data:image/png;base64,iVBORw0KGgo=">`,
	}

	for _, input := range tests {
		assert.Equal(t, input, p.Sanitize(input), input)
	}
}
