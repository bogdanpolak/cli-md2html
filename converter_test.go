package main

import (
	"strings"
	"testing"

	td "github.com/maxatome/go-testdeep/td"
)

func TestGenerateHtmlBodyInternalContent(t *testing.T) {
	markdown := `## Header

Some content.`

	result := GenerateHtmlBodyInternalContent(markdown)

	// Should not contain full HTML document structure
	if strings.Contains(result, "<!DOCTYPE html>") || strings.Contains(result, "<html>") {
		t.Error("GenerateHtmlBodyInternalContent should not contain full HTML document structure")
	}

	// Should contain the converted content
	if !strings.Contains(result, "<h2>Header</h2>") {
		t.Error("Expected h2 tag with 'Header'")
	}

	if !strings.Contains(result, "<p>Some content.</p>") {
		t.Error("Expected paragraph tag")
	}
}

// ============================================================================
// HEADERS
// ============================================================================

func TestHeaderConversion(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "01 H1 header",
			markdown: "# Main Title",
			expected: "<h1>Main Title</h1>\n",
		},
		{
			name:     "02 H2 header",
			markdown: "## Subtitle",
			expected: "<h2>Subtitle</h2>\n",
		},
		{
			name:     "03 H3 header",
			markdown: "### Sub Subtitle",
			expected: "<h3>Sub Subtitle</h3>\n",
		},
		{
			name:     "04 H1 with inline code",
			markdown: "# Title with `code`",
			expected: "<h1>Title with <code>code</code></h1>\n",
		},
		{
			name:     "05 H2 with link",
			markdown: "## See [docs](https://example.com)",
			expected: "<h2>See <a href=\"https://example.com\">docs</a></h2>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td.Cmp(t, GenerateHtmlBodyInternalContent(tt.markdown), tt.expected)
		})
	}
}

// ============================================================================
// LISTS
// ============================================================================

func TestUnorderedListConversion(t *testing.T) {
	tests := []struct {
		name     string
		markdown []string
		expected []string
	}{
		{
			name:     "01 Single item list",
			markdown: []string{"- Item 1"},
			expected: []string{
				"<ul>",
				"∘<li>Item 1</li>",
				"</ul>",
				""},
		},
		{
			name: "02 Multiple items",
			markdown: []string{
				"- First",
				"- Second",
				"- Third"},
			expected: []string{
				"<ul>",
				"∘<li>First</li>",
				"∘<li>Second</li>",
				"∘<li>Third</li>",
				"</ul>",
				""},
		},
		{
			name: "03 Nested list items",
			markdown: []string{
				"- First",
				"  - First Child",
				"  - Second Child",
				"- Second"},
			expected: []string{
				"<ul>",
				"∘<li>First",
				"∘∘<ul>",
				"∘∘∘<li>First Child</li>",
				"∘∘∘<li>Second Child</li>",
				"∘∘</ul>",
				"∘</li>",
				"∘<li>Second</li>",
				"</ul>",
				""},
		},
		{
			name: "04 List with inline code",
			markdown: []string{
				"- Run `npm install`"},
			expected: []string{
				"<ul>",
				"∘<li>Run <code>npm install</code></li>",
				"</ul>",
				""},
		},
		{
			name: "05 List with link",
			markdown: []string{
				"- Visit [GitHub](https://github.com)"},
			expected: []string{
				"<ul>",
				"∘<li>Visit <a href=\"https://github.com\">GitHub</a></li>",
				"</ul>",
				""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := strings.ReplaceAll(strings.Join(tt.expected, "\n"), "∘", "    ")
			markdown := strings.Join(tt.markdown, "\n")
			td.Cmp(t, GenerateHtmlBodyInternalContent(markdown), expected)
		})
	}
}

func TestOrderedListConversion(t *testing.T) {
	tests := []struct {
		name     string
		markdown []string
		expected []string
	}{
		{
			name: "01 Single ordered item",
			markdown: []string{
				"1. First step"},
			expected: []string{
				"<ol>",
				"→<li>First step</li>",
				"</ol>", ""},
		},
		{
			name: "02 Multiple ordered items",
			markdown: []string{
				"1. Step one",
				"2. Step two",
				"3. Step three"},
			expected: []string{
				"<ol>",
				"→<li>Step one</li>",
				"→<li>Step two</li>",
				"→<li>Step three</li>",
				"</ol>", ""},
		},
		{
			name: "03 Ordered list with code",
			markdown: []string{
				"1. Install `package`",
				"2. Run `build`"},
			expected: []string{
				"<ol>",
				"→<li>Install <code>package</code></li>",
				"→<li>Run <code>build</code></li>",
				"</ol>", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := strings.ReplaceAll(strings.Join(tt.expected, "\n"), "→", "    ")
			markdown := strings.Join(tt.markdown, "\n")
			actual := GenerateHtmlBodyInternalContent(markdown)
			td.Cmp(t, actual, expected)
		})
	}
}

// ============================================================================
// CODE BLOCKS
// ============================================================================

func TestCodeBlockConversion(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "01 Simple",
			markdown: "```\ncode here\n```",
			expected: `<div class="source-code">
<pre><code>code here</code></pre>
</div>
`,
		},
		{
			name: "02 Go function",
			markdown: "```\n" + `func main() {
    fmt.Println("Hello")
}
` + "```\n",
			expected: `<div class="source-code">
<pre><code>func main() {
    fmt.Println(&quot;Hello&quot;)
}</code></pre>
</div>
`,
		},
		{
			name:     "03 Code block with HTML characters",
			markdown: "```\n<div>test</div>\n```",
			expected: `<div class="source-code">
<pre><code>&lt;div&gt;test&lt;/div&gt;</code></pre>
</div>
`,
		},
		{
			name:     "04 Multiple code blocks",
			markdown: "```\nfirst\n```\n\n```\nsecond\n```",
			expected: `<div class="source-code">
<pre><code>first</code></pre>
</div>

<div class="source-code">
<pre><code>second</code></pre>
</div>
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td.Cmp(t, GenerateHtmlBodyInternalContent(tt.markdown), tt.expected)
		})
	}
}

// ============================================================================
// INLINE CODE
// ============================================================================

func TestInlineCodeConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Inline_code_simple",
			input:    "`function()`",
			expected: "<code>function()</code>",
		},
		{
			name:     "Inline_code_Multiple_elements",
			input:    "Call function `foo()` and `bar()`",
			expected: "Call function <code>foo()</code> and <code>bar()</code>",
		},
		{
			name:     "Inline_code_with_special_chars",
			input:    "Use `<tag>`",
			expected: "Use <code>&lt;tag&gt;</code>",
		},
		{
			name:     "Inline_code_at_start",
			input:    "`code` at start",
			expected: "<code>code</code> at start",
		},
		{
			name:     "Inline_code_at_end",
			input:    "at end `code`",
			expected: "at end <code>code</code>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td.Cmp(t, processInlineElements(tt.input), tt.expected)
		})
	}
}

func TestLinkConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Markdown link",
			input:    "[GitHub](https://github.com)",
			expected: `<a href="https://github.com">GitHub</a>`,
		},
		{
			name:     "Auto-detected HTTPS URL",
			input:    "Visit https://example.com to get started",
			expected: `Visit <a href="https://example.com">https://example.com</a> to get started`,
		},
		{
			name:     "Auto-detected HTTP URL",
			input:    "Go to http://example.com",
			expected: `Go to <a href="http://example.com">http://example.com</a>`,
		},
		{
			name:     "Multiple links",
			input:    "[Google](https://google.com) and [Bing](https://bing.com)",
			expected: `<a href="https://google.com">Google</a> and <a href="https://bing.com">Bing</a>`,
		},
		{
			name:     "Link with special chars in URL",
			input:    "[Docs](https://example.com/docs?id=123&format=html)",
			expected: `<a href="https://example.com/docs?id=123&amp;format=html">Docs</a>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td.Cmp(t, processInlineElements(tt.input), tt.expected)
		})
	}
}

// ============================================================================
// HTML ESCAPING TESTS
// ============================================================================

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Ampersand escape",
			input:    "Tom & Jerry",
			expected: "Tom &amp; Jerry",
		},
		{
			name:     "Less than escape",
			input:    "5 < 10",
			expected: "5 &lt; 10",
		},
		{
			name:     "Greater than escape",
			input:    "10 > 5",
			expected: "10 &gt; 5",
		},
		{
			name:     "Double quote escape",
			input:    `Say "hello"`,
			expected: "Say &quot;hello&quot;",
		},
		{
			name:     "Single quote escape",
			input:    "It's great",
			expected: "It&#39;s great",
		},
		{
			name:     "Script tag injection",
			input:    `<script>alert("XSS")</script>`,
			expected: `&lt;script&gt;alert(&quot;XSS&quot;)&lt;/script&gt;`,
		},
		{
			name:     "Multiple special chars",
			input:    `<script>alert("test & run")</script>`,
			expected: `&lt;script&gt;alert(&quot;test &amp; run&quot;)&lt;/script&gt;`,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeHTML(tt.input)
			td.Cmp(t, result, tt.expected)
		})
	}
}

// ============================================================================
// EDGE CASES AND MALFORMED INPUT
// ============================================================================

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		markdown       string
		shouldNotCrash bool
	}{
		{
			name:           "Empty string",
			markdown:       "",
			shouldNotCrash: true,
		},
		{
			name:           "Only whitespace",
			markdown:       "   \n\n   ",
			shouldNotCrash: true,
		},
		{
			name:           "Only newlines",
			markdown:       "\n\n\n",
			shouldNotCrash: true,
		},
		{
			name:           "Unclosed code block",
			markdown:       "```\ncode without closing",
			shouldNotCrash: true,
		},
		{
			name:           "Malformed link",
			markdown:       "[link without URL]()",
			shouldNotCrash: true,
		},
		{
			name:           "Nested code blocks",
			markdown:       "```\n```\nnested\n```\n```",
			shouldNotCrash: true,
		},
		{
			name:           "Mixed lists",
			markdown:       "- Item 1\n1. Ordered\n- Item 2",
			shouldNotCrash: true,
		},
		{
			name:           "Very long line",
			markdown:       strings.Repeat("a", 10000),
			shouldNotCrash: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("unexpected panic: %v", r)
				}
			}()
			result := GenerateHtmlBodyInternalContent(tt.markdown)
			if result == "" && tt.markdown != "" && tt.name != "Unclosed code block" {
				// Allow empty result for some edge cases
			}
		})
	}
}

// ============================================================================
// INTEGRATION TESTS
// ============================================================================

func TestConvertWithTemplate(t *testing.T) {
	markdown := "# Hello World"
	template := "<html><head><title>{{ .Title }}</title></head><body>{{ .Content }}</body></html>"
	title := "TestABC"

	result, err := ConvertMarkdownToHTML(markdown, template, title)

	td.Cmp(t, err, nil)
	td.Cmp(t, result, "<html><head><title>TestABC</title></head><body><h1>Hello World</h1>\n</body></html>")
}

func TestComplexDocument(t *testing.T) {
	markdown := `# Main Title

This is an introduction with a [link](https://example.com) and some ` + "`inline code`" + `.

## Section 1

Some text with ` + "`npm install`" + ` command.

- Install dependencies
- Run tests
- Deploy

## Section 2

### Nested subsection

#### Level four - Nested subsection

1. First step
2. Second step with ` + "`code`" + `

` + "```\n" + `func main() {
    fmt.Println("Hello")
}
` + "```\n" + `
Visit https://github.com for more info.`

	result := GenerateHtmlBodyInternalContent(markdown)

	td.Cmp(t, result, td.All(
		td.Contains("<h1>Main Title</h1>"),
		td.Contains("<p>This is an introduction with a <a href=\"https://example.com\">link</a> and some <code>inline code</code>.</p>"),
		td.Contains("<h2>Section 1</h2>"),
		td.Contains("<p>Some text with <code>npm install</code> command.</p>"),
		td.Contains("<ul>"),
		td.Contains("<li>Install dependencies</li>"),
		td.Contains("<li>Run tests</li>"),
		td.Contains("<li>Deploy</li>"),
		td.Contains("</ul>"),
		td.Contains("<h2>Section 2</h2>"),
		td.Contains("<h3>Nested subsection</h3>"),
		td.Contains("<h4>Level four - Nested subsection</h4>"),
		td.Contains("<ol>"),
		td.Contains("<li>First step</li>"),
		td.Contains("<li>Second step with <code>code</code></li>"),
		td.Contains("</ol>"),
		td.Contains(`<div class="source-code">
<pre><code>func main() {
    fmt.Println(&quot;Hello&quot;)
}</code></pre>
</div>
`),
		td.Contains("<p>Visit <a href=\"https://github.com\">https://github.com</a> for more info.</p>"),
	))
}

func TestProcessInlineElements(t *testing.T) {
	// Test inline code
	input := "This has `inline code` in it."
	result := processInlineElements(input)

	if !strings.Contains(result, "<code>inline code</code>") {
		t.Error("Expected inline code to be converted to <code> tags")
	}

	// Test links
	input = "Check out [example](https://example.com) link."
	result = processInlineElements(input)

	if !strings.Contains(result, `<a href="https://example.com">example</a>`) {
		t.Error("Expected link to be converted to <a> tags")
	}
}
