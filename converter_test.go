package main

import (
	"strings"
	"testing"
)

// ============================================================================
// MARKDOWN TO HTML TESTS
// ============================================================================

func TestMarkdownToHTML(t *testing.T) {
	markdown := `# Test Header

This is a test paragraph.

- Item 1
- Item 2`

	result := MarkdownToHTML(markdown)

	// Check if the result contains expected HTML elements
	if !strings.Contains(result, "<h1>Test Header</h1>") {
		t.Error("Expected h1 tag with 'Test Header'")
	}

	if !strings.Contains(result, "<p>This is a test paragraph.</p>") {
		t.Error("Expected paragraph tag")
	}

	if !strings.Contains(result, "<ul>") || !strings.Contains(result, "</ul>") {
		t.Error("Expected unordered list tags")
	}

	if !strings.Contains(result, "<li>Item 1</li>") {
		t.Error("Expected list item 'Item 1'")
	}

	// Must contain full HTML document structure
	if !strings.Contains(result, "<!DOCTYPE html>") {
		t.Error("Expected <!DOCTYPE html>")
	}
	if !strings.Contains(result, "<html>") || !strings.Contains(result, "</html>") {
		t.Error("Expected <html> tags")
	}
	if !strings.Contains(result, "<body>") || !strings.Contains(result, "</body>") {
		t.Error("Expected <body> tags")
	}
}

func TestMarkdownToHTMLContent(t *testing.T) {
	markdown := `## Header

Some content.`

	result := MarkdownToHTMLContent(markdown)

	// Should not contain full HTML document structure
	if strings.Contains(result, "<!DOCTYPE html>") || strings.Contains(result, "<html>") {
		t.Error("MarkdownToHTMLContent should not contain full HTML document structure")
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
// TABLE-DRIVEN TESTS FOR HEADERS
// ============================================================================

func TestHeaderConversion(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "H1 header",
			markdown: "# Main Title",
			expected: "<h1>Main Title</h1>",
		},
		{
			name:     "H2 header",
			markdown: "## Subtitle",
			expected: "<h2>Subtitle</h2>",
		},
		{
			name:     "H1 with inline code",
			markdown: "# Title with `code`",
			expected: "<h1>Title with <code>code</code></h1>",
		},
		{
			name:     "H2 with link",
			markdown: "## See [docs](https://example.com)",
			expected: "<h2>See <a href=\"https://example.com\">docs</a></h2>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MarkdownToHTMLContent(tt.markdown)
			if !strings.Contains(result, tt.expected) {
				t.Errorf("expected %q to contain %q", result, tt.expected)
			}
		})
	}
}

// ============================================================================
// TABLE-DRIVEN TESTS FOR LISTS
// ============================================================================

func TestUnorderedListConversion(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "Single item list",
			markdown: "- Item 1",
			expected: "<ul>\n    <li>Item 1</li>\n</ul>",
		},
		{
			name:     "Multiple items",
			markdown: "- First\n- Second\n- Third",
			expected: "<li>First</li>",
		},
		{
			name:     "List with inline code",
			markdown: "- Run `npm install`",
			expected: "<li>Run <code>npm install</code></li>",
		},
		{
			name:     "List with link",
			markdown: "- Visit [GitHub](https://github.com)",
			expected: "<li>Visit <a href=\"https://github.com\">GitHub</a></li>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MarkdownToHTMLContent(tt.markdown)
			if !strings.Contains(result, tt.expected) {
				t.Errorf("expected %q to contain %q", result, tt.expected)
			}
		})
	}
}

func TestOrderedListConversion(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "Single ordered item",
			markdown: "1. First step",
			expected: "<ol>\n    <li>First step</li>\n</ol>",
		},
		{
			name:     "Multiple ordered items",
			markdown: "1. Step one\n2. Step two\n3. Step three",
			expected: "<li>Step one</li>",
		},
		{
			name:     "Ordered list with code",
			markdown: "1. Install `package`\n2. Run `build`",
			expected: "<code>package</code>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MarkdownToHTMLContent(tt.markdown)
			if !strings.Contains(result, tt.expected) {
				t.Errorf("expected %q to contain %q", result, tt.expected)
			}
		})
	}
}

// ============================================================================
// TABLE-DRIVEN TESTS FOR CODE BLOCKS
// ============================================================================

func TestCodeBlockConversion(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "Simple code block",
			markdown: "```\ncode here\n```",
			expected: "<div class=\"source-code\">",
		},
		{
			name:     "Code block with language hint",
			markdown: "```go\nfunc main() {}\n```",
			expected: "func main() {}",
		},
		{
			name:     "Code block with HTML characters",
			markdown: "```\n<div>test</div>\n```",
			expected: "&lt;div&gt;test&lt;/div&gt;",
		},
		{
			name:     "Multiple code blocks",
			markdown: "```\nfirst\n```\n\n```\nsecond\n```",
			expected: "first",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MarkdownToHTMLContent(tt.markdown)
			if !strings.Contains(result, tt.expected) {
				t.Errorf("expected %q to contain %q", result, tt.expected)
			}
		})
	}
}

// ============================================================================
// TABLE-DRIVEN TESTS FOR INLINE ELEMENTS
// ============================================================================

func TestInlineCodeConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple inline code",
			input:    "Use `function()`",
			expected: "<code>function()</code>",
		},
		{
			name:     "Multiple inline codes",
			input:    "Call `foo()` and `bar()`",
			expected: "<code>foo()</code>",
		},
		{
			name:     "Inline code with special chars",
			input:    "Use `<tag>`",
			expected: "<code>&lt;tag&gt;</code>",
		},
		{
			name:     "Inline code at start",
			input:    "`code` at start",
			expected: "<code>code</code>",
		},
		{
			name:     "Inline code at end",
			input:    "at end `code`",
			expected: "<code>code</code>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processInlineElements(tt.input)
			if !strings.Contains(result, tt.expected) {
				t.Errorf("expected %q to contain %q", result, tt.expected)
			}
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
			input:    "Visit https://example.com today",
			expected: `<a href="https://example.com">https://example.com</a>`,
		},
		{
			name:     "Auto-detected HTTP URL",
			input:    "Go to http://example.com",
			expected: `<a href="http://example.com">http://example.com</a>`,
		},
		{
			name:     "Multiple links",
			input:    "[Google](https://google.com) and [Bing](https://bing.com)",
			expected: `<a href="https://google.com">Google</a>`,
		},
		{
			name:     "Link with special chars in URL",
			input:    "[Docs](https://example.com/docs?id=123&format=html)",
			expected: `<a href="https://example.com/docs?id=123&amp;format=html">Docs</a>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processInlineElements(tt.input)
			if !strings.Contains(result, tt.expected) {
				t.Errorf("expected %q to contain %q", result, tt.expected)
			}
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
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
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
			result := MarkdownToHTMLContent(tt.markdown)
			if result == "" && tt.markdown != "" && tt.name != "Unclosed code block" {
				// Allow empty result for some edge cases
			}
		})
	}
}

// ============================================================================
// INTEGRATION TESTS
// ============================================================================

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

1. First step
2. Second step with ` + "`code`" + `

` + "```go\n" + `func main() {
    fmt.Println("Hello")
}
` + "```\n" + `
Visit https://github.com for more info.`

	result := MarkdownToHTMLContent(markdown)

	// Verify structure
	if !strings.Contains(result, "<h1>Main Title</h1>") {
		t.Error("Missing h1")
	}
	if !strings.Contains(result, "<h2>Section 1</h2>") {
		t.Error("Missing h2")
	}
	if !strings.Contains(result, "<ul>") {
		t.Error("Missing unordered list")
	}
	if !strings.Contains(result, "<ol>") {
		t.Error("Missing ordered list")
	}
	if !strings.Contains(result, "source-code") {
		t.Error("Missing code block")
	}
	if !strings.Contains(result, "<a href=") {
		t.Error("Missing links")
	}
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
