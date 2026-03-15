package main

import (
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

type simpleTestCase struct {
	name     string
	markdown string
	expected string
}

type multilineTestCase struct {
	name     string
	markdown []string
	expected []string
}

const indentHtmlWith4Spaces = "    "

func (tt multilineTestCase) toString(indentation string) (string, string) {
	expected := strings.ReplaceAll(strings.Join(tt.expected, "\n"), "•", indentation)
	markdown := strings.Join(tt.markdown, "\n")
	return expected, markdown
}

// ---------------------------------------------------------------------------
// Simple conversions
// ---------------------------------------------------------------------------

func TestSimpleConversion(t *testing.T) {
	tests := []simpleTestCase{
		{
			name:     "01 Empty markup",
			markdown: "",
			expected: "\n",
		},
		{
			name:     "02 Random Single line Text",
			markdown: "Lorem impsum dolor sit amet.",
			expected: "<p>Lorem impsum dolor sit amet.</p>\n",
		},
		{
			name:     "01 Empty markup",
			markdown: "",
			expected: "\n",
		},
		{
			name:     "02 Random Single line Text",
			markdown: "Lorem impsum dolor sit amet.",
			expected: "<p>Lorem impsum dolor sit amet.</p>\n",
		},
		{
			name:     "03 Two paragraphs",
			markdown: "Paragraph One.\nParagraph Two.",
			expected: "<p>Paragraph One.</p>\n<p>Paragraph Two.</p>\n",
		},
		{
			name:     "04 Only whitespace",
			markdown: "   \n\n   ",
			expected: "\n",
		},
		{
			name:     "05 Only newlines",
			markdown: "\n\n\n",
			expected: "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td.Cmp(t, GenerateHtmlBody(tt.markdown), tt.expected)
		})
	}
}

// ---------------------------------------------------------------------------
// Headers
// ---------------------------------------------------------------------------

func TestHeaderConversion(t *testing.T) {
	tests := []simpleTestCase{
		{
			name:     "01 H1 header",
			markdown: "# Main Title",
			expected: "<h1>Main Title</h1>",
		},
		{
			name:     "02 H2 header",
			markdown: "## Subtitle",
			expected: "<h2>Subtitle</h2>",
		},
		{
			name:     "03 H3 header",
			markdown: "### Sub Subtitle",
			expected: "<h3>Sub Subtitle</h3>",
		},
		{
			name:     "04 H1 with inline code",
			markdown: "# Title with `code`",
			expected: "<h1>Title with <code>code</code></h1>",
		},
		{
			name:     "05 H2 with link",
			markdown: "## See [docs](https://example.com)",
			expected: "<h2>See <a href=\"https://example.com\">docs</a></h2>",
		},
		{
			name:     "06 Bold text in header",
			markdown: "#### Title with **bold** text",
			expected: "<h4>Title with <strong>bold</strong> text</h4>",
		},
		{
			name:     "07 Italic text in paragraph",
			markdown: "#### Title with *italic* text",
			expected: "<h4>Title with <em>italic</em> text</h4>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td.Cmp(t, processSingleLine(tt.markdown), tt.expected)
		})
	}
}

// ---------------------------------------------------------------------------
// Block Quote conversions
// ---------------------------------------------------------------------------

func TestBlockQuoteConversion(t *testing.T) {
	tests := []simpleTestCase{
		{
			name:     "01 Basic Block Quote",
			markdown: "> Lorem impsum dolor sit amet.",
			expected: "<blockquote>Lorem impsum dolor sit amet.</blockquote>\n",
		},
		{
			name:     "02 Block Quote with Callout",
			markdown: "> Rule of thumb: Lorem impsum dolor sit amet.",
			expected: "<blockquote><strong>Rule of thumb:</strong> Lorem impsum dolor sit amet.</blockquote>\n",
		},
		{
			name:     "03 Callout ignored after comma",
			markdown: "> Rule of thumb, actually: Lorem impsum dolor sit amet.",
			expected: "<blockquote>Rule of thumb, actually: Lorem impsum dolor sit amet.</blockquote>\n",
		},
		{
			name:     "04 Callout ignored after period",
			markdown: "> Rule of thumb. Actually: Lorem impsum dolor sit amet.",
			expected: "<blockquote>Rule of thumb. Actually: Lorem impsum dolor sit amet.</blockquote>\n",
		},
		{
			name:     "05 Callout ignored after semicolon",
			markdown: "> Rule of thumb; actually: Lorem impsum dolor sit amet.",
			expected: "<blockquote>Rule of thumb; actually: Lorem impsum dolor sit amet.</blockquote>\n",
		},
		{
			name:     "06 Callout ignored after hyphen",
			markdown: "> Rule of thumb - actually: Lorem impsum dolor sit amet.",
			expected: "<blockquote>Rule of thumb - actually: Lorem impsum dolor sit amet.</blockquote>\n",
		},
		{
			name:     "07 Block Quote preserves inline formatting",
			markdown: "> Use **bold**, *italic*, `code`, and [link](https://example.com)",
			expected: "<blockquote>Use <strong>bold</strong>, <em>italic</em>, <code>code</code>, and <a href=\"https://example.com\">link</a></blockquote>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td.Cmp(t, GenerateHtmlBody(tt.markdown), tt.expected)
		})
	}
}

// ---------------------------------------------------------------------------
// Inline code
// ---------------------------------------------------------------------------

func TestInlineCodeConversion(t *testing.T) {
	tests := []simpleTestCase{
		{
			name:     "Inline_code_simple",
			markdown: "`function()`",
			expected: "<p><code>function()</code></p>",
		},
		{
			name:     "Inline_code_Multiple_elements",
			markdown: "Call function `foo()` and `bar()`",
			expected: "<p>Call function <code>foo()</code> and <code>bar()</code></p>",
		},
		{
			name:     "Inline_code_with_special_chars",
			markdown: "Use `<tag>`",
			expected: "<p>Use <code>&lt;tag&gt;</code></p>",
		},
		{
			name:     "Inline_code_at_start",
			markdown: "`code` at start",
			expected: "<p><code>code</code> at start</p>",
		},
		{
			name:     "Inline_code_at_end",
			markdown: "at end `code`",
			expected: "<p>at end <code>code</code></p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td.Cmp(t, processSingleLine(tt.markdown), tt.expected)
		})
	}
}

// ---------------------------------------------------------------------------
// Hyperlinks
// ---------------------------------------------------------------------------

func TestLinkConversion(t *testing.T) {
	tests := []simpleTestCase{
		{
			name:     "Markdown link",
			markdown: "[GitHub](https://github.com)",
			expected: `<p><a href="https://github.com">GitHub</a></p>`,
		},
		{
			name:     "Auto-detected HTTPS URL",
			markdown: "Visit https://example.com to get started",
			expected: `<p>Visit <a href="https://example.com">https://example.com</a> to get started</p>`,
		},
		{
			name:     "Auto-detected HTTP URL",
			markdown: "Go to http://example.com",
			expected: `<p>Go to <a href="http://example.com">http://example.com</a></p>`,
		},
		{
			name:     "Multiple links",
			markdown: "[Google](https://google.com) and [Bing](https://bing.com)",
			expected: `<p><a href="https://google.com">Google</a> and <a href="https://bing.com">Bing</a></p>`,
		},
		{
			name:     "Link with special chars in URL",
			markdown: "[Docs](https://example.com/docs?id=123&format=html)",
			expected: `<p><a href="https://example.com/docs?id=123&amp;format=html">Docs</a></p>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td.Cmp(t, processSingleLine(tt.markdown), tt.expected)
		})
	}
}

// ---------------------------------------------------------------------------
// Html Escaping
// ---------------------------------------------------------------------------

func TestEscapeHTML(t *testing.T) {
	tests := []simpleTestCase{
		{
			name:     "Ampersand escape",
			markdown: "Tom & Jerry",
			expected: "<p>Tom &amp; Jerry</p>",
		},
		{
			name:     "Less than escape",
			markdown: "5 < 10",
			expected: "<p>5 &lt; 10</p>",
		},
		{
			name:     "Greater than escape",
			markdown: "10 > 5",
			expected: "<p>10 &gt; 5</p>",
		},
		{
			name:     "Double quote escape",
			markdown: `Say "hello"`,
			expected: "<p>Say &quot;hello&quot;</p>",
		},
		{
			name:     "Single quote escape",
			markdown: "It's great",
			expected: "<p>It&#39;s great</p>",
		},
		{
			name:     "Script tag injection",
			markdown: `<script>alert("XSS")</script>`,
			expected: `<p>&lt;script&gt;alert(&quot;XSS&quot;)&lt;/script&gt;</p>`,
		},
		{
			name:     "Multiple special chars",
			markdown: `<script>alert("test & run")</script>`,
			expected: `<p>&lt;script&gt;alert(&quot;test &amp; run&quot;)&lt;/script&gt;</p>`,
		},
		{
			name:     "Empty markup",
			markdown: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td.Cmp(t, processSingleLine(tt.markdown), tt.expected)
		})
	}
}

// ---------------------------------------------------------------------------
// Block - Ordered and Point Lists
// ---------------------------------------------------------------------------

func TestUnorderedListConversion(t *testing.T) {
	tests := []multilineTestCase{
		{
			name: "01 Single item list",
			markdown: []string{
				"- Item 1"},
			expected: []string{
				"<ul>",
				"•<li>Item 1</li>",
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
				"•<li>First</li>",
				"•<li>Second</li>",
				"•<li>Third</li>",
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
				"•<li>First",
				"••<ul>",
				"•••<li>First Child</li>",
				"•••<li>Second Child</li>",
				"••</ul>",
				"•</li>",
				"•<li>Second</li>",
				"</ul>",
				""},
		},
		{
			name: "04 Different spaces for nesting",
			markdown: []string{
				"   - First level",
				"       - Second level 1",
				"       - Second level 2",
				"                 - Third level"},
			expected: []string{
				"<ul>",
				"•<li>First level",
				"••<ul>",
				"•••<li>Second level 1</li>",
				"•••<li>Second level 2",
				"••••<ul>",
				"•••••<li>Third level</li>",
				"••••</ul>",
				"•••</li>",
				"••</ul>",
				"•</li>",
				"</ul>",
				""},
		},
		{
			name: "05 List with inline code",
			markdown: []string{
				"- Run `npm install`"},
			expected: []string{
				"<ul>",
				"•<li>Run <code>npm install</code></li>",
				"</ul>",
				""},
		},
		{
			name: "06 List with link",
			markdown: []string{
				"- Visit [GitHub](https://github.com)"},
			expected: []string{
				"<ul>",
				"•<li>Visit <a href=\"https://github.com\">GitHub</a></li>",
				"</ul>",
				""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected, markdown := tt.toString(indentHtmlWith4Spaces)
			td.Cmp(t, GenerateHtmlBody(markdown), expected)
		})
	}
}

func TestOrderedListConversion(t *testing.T) {
	tests := []multilineTestCase{
		{
			name: "01 Single ordered item",
			markdown: []string{
				"1. First step"},
			expected: []string{
				"<ol>",
				"•<li>First step</li>",
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
				"•<li>Step one</li>",
				"•<li>Step two</li>",
				"•<li>Step three</li>",
				"</ol>", ""},
		},
		{
			name: "03 Ordered list with non-sequential numbers",
			markdown: []string{
				"3. Originally nr 3",
				"6. Originally nr 6",
				"7. Originally nr 7"},
			expected: []string{
				"<ol>",
				"•<li>Originally nr 3</li>",
				"•<li>Originally nr 6</li>",
				"•<li>Originally nr 7</li>",
				"</ol>", ""},
		},
		{
			name: "04 Ordered list with code",
			markdown: []string{
				"1. Install `package`",
				"2. Run `build`"},
			expected: []string{
				"<ol>",
				"•<li>Install <code>package</code></li>",
				"•<li>Run <code>build</code></li>",
				"</ol>", ""},
		},
		{
			name: "05 List with paragraph bellow",
			markdown: []string{
				"1. First",
				"Hello World"},
			expected: []string{
				"<ol>",
				"•<li>First</li>",
				"</ol>",
				"<p>Hello World</p>",
				""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected, markdown := tt.toString(indentHtmlWith4Spaces)
			td.Cmp(t, GenerateHtmlBody(markdown), expected)
		})
	}
}

func TestMixedListConversion(t *testing.T) {
	tests := []multilineTestCase{
		{
			name: "01. Ordered list with point subitems",
			markdown: []string{
				"1. Ordered",
				"2. Ordered Two",
				"   - Point Subitem 2.1"},
			expected: []string{
				"<ol>",
				"•<li>Ordered</li>",
				"•<li>Ordered Two",
				"••<ul>",
				"•••<li>Point Subitem 2.1</li>",
				"••</ul>",
				"•</li>",
				"</ol>",
				"",
			},
		},
		{
			name: "02. Point lists with ordered subitems",
			markdown: []string{
				"- Main",
				"   1. Ordered One",
				"   2. Ordered Two",
				"- Another Main"},
			expected: []string{
				"<ul>",
				"•<li>Main",
				"••<ol>",
				"•••<li>Ordered One</li>",
				"•••<li>Ordered Two</li>",
				"••</ol>",
				"•</li>",
				"•<li>Another Main</li>",
				"</ul>",
				"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected, markdown := tt.toString(indentHtmlWith4Spaces)
			td.Cmp(t, GenerateHtmlBody(markdown), expected)
		})
	}
}

// ---------------------------------------------------------------------------
// Block - Code
// ---------------------------------------------------------------------------

func TestCodeBlockConversion(t *testing.T) {
	tests := []multilineTestCase{
		{
			name: "01 Simple",
			markdown: []string{
				"```",
				"code here",
				"```"},
			expected: []string{
				"<section class=\"code\">",
				"<pre><code>code here</code></pre>",
				"</section>",
				""},
		},
		{
			name: "02 Go function",
			markdown: []string{
				"```",
				"func main() {",
				"    fmt.Println(\"Hello\")",
				"}",
				"```"},
			expected: []string{
				"<section class=\"code\">",
				"<pre><code>func main() {",
				"    fmt.Println(&quot;Hello&quot;)",
				"}</code></pre>",
				"</section>",
				""},
		},
		{
			name: "03 Code block with HTML characters",
			markdown: []string{
				"```",
				"<div>test</div>",
				"```"},
			expected: []string{
				"<section class=\"code\">",
				"<pre><code>&lt;div&gt;test&lt;/div&gt;</code></pre>",
				"</section>",
				""},
		},
		{
			name: "04 Multiple code blocks",
			markdown: []string{
				"```",
				"first",
				"```",
				"",
				"```",
				"second",
				"```"},
			expected: []string{
				"<section class=\"code\">",
				"<pre><code>first</code></pre>",
				"</section>",
				"",
				"<section class=\"code\">",
				"<pre><code>second</code></pre>",
				"</section>",
				""},
		},
		{
			name: "05 Unclosed code block",
			markdown: []string{
				"```",
				"const x = 5",
				"/* without closing hyphens */"},
			expected: []string{
				"<section class=\"code\">",
				"<pre><code>const x = 5",
				"/* without closing hyphens */</code></pre>",
				"</section>\n"},
		},
		{
			name: "06 HTML code block",
			markdown: []string{
				"```",
				"<header class='intro'>",
				"    <a href='abc.com?a=1&b=2'>Link</a>",
				"</header>",
				"```"},
			expected: []string{
				"<section class=\"code\">",
				"<pre><code>&lt;header class=&#39;intro&#39;&gt;",
				"    &lt;a href=&#39;abc.com?a=1&amp;b=2&#39;&gt;Link&lt;/a&gt;",
				"&lt;/header&gt;</code></pre>",
				"</section>",
				""},
		},
		{
			name: "07 Language-qualified code block keeps following prose outside",
			markdown: []string{
				"```pascal",
				"function Test: Integer;",
				"begin",
				"  Result := 42;",
				"end;",
				"```",
				"Following paragraph.",
			},
			expected: []string{
				"<section class=\"code\">",
				"<pre><code>function Test: Integer;",
				"begin",
				"  Result := 42;",
				"end;</code></pre>",
				"</section>",
				"<p>Following paragraph.</p>",
				""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected, markdown := tt.toString(indentHtmlWith4Spaces)
			td.Cmp(t, GenerateHtmlBody(markdown), expected)
		})
	}
}

// ---------------------------------------------------------------------------
// Block - Lists with Code
// ---------------------------------------------------------------------------

func TestListWithCodeBlockConversion(t *testing.T) {
	tests := []multilineTestCase{
		{
			name: "01 One line code in single point list",
			markdown: []string{
				"2. Create `001` migration",
				"   ```",
				"   migrate create -ext sql -dir ./migrations -seq \"Create Sets table\"",
				"   ```"},
			expected: []string{
				"<ol>",
				"•<li>Create <code>001</code> migration",
				"••<section class=\"code\">",
				"••<pre><code>migrate create -ext sql -dir ./migrations -seq &quot;Create Sets table&quot;</code></pre>",
				"••</section>",
				"•</li>",
				"</ol>",
				""},
		},
		{
			name: "02 Code block (3lines) separated with space in ordered list",
			markdown: []string{
				"2. Second step with `code`",
				"",
				"```",
				"func main() {",
				"      fmt.Println(\"Hello\")",
				"}",
				"```"},
			expected: []string{
				"<ol>",
				"•<li>Second step with <code>code</code>",
				"••<section class=\"code\">",
				"••<pre><code>func main() {",
				"      fmt.Println(&quot;Hello&quot;)",
				"}</code></pre>",
				"••</section>",
				"•</li>",
				"</ol>",
				""},
		},
		{
			name: "03 Nested list with code block",
			markdown: []string{
				"- One",
				"   - Level Two 1",
				"   ```",
				"     var x int = 10",
				"   ```",
				"   - Level Two 2",
				"- Two",
			},
			expected: []string{
				"<ul>",
				"•<li>One",
				"••<ul>",
				"•••<li>Level Two 1",
				"••••<section class=\"code\">",
				"••••<pre><code>  var x int = 10</code></pre>",
				"••••</section>",
				"•••</li>",
				"•••<li>Level Two 2</li>",
				"••</ul>",
				"•</li>",
				"•<li>Two</li>",
				"</ul>",
				""},
		},
		{
			name: "04 Ordered list with language-qualified code block",
			markdown: []string{
				"1. First step",
				"   ```pascal",
				"   WriteLn('Hello');",
				"   ```",
				"2. Second step",
			},
			expected: []string{
				"<ol>",
				"•<li>First step",
				"••<section class=\"code\">",
				"••<pre><code>WriteLn(&#39;Hello&#39;);</code></pre>",
				"••</section>",
				"•</li>",
				"•<li>Second step</li>",
				"</ol>",
				""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected, markdown := tt.toString(indentHtmlWith4Spaces)
			td.Cmp(t, GenerateHtmlBody(markdown), expected)
		})
	}
}

// ---------------------------------------------------------------------------
// Block Quote conversions
// ---------------------------------------------------------------------------

func TestYamlFrontmatterConversion(t *testing.T) {
	tests := []multilineTestCase{
		{
			name: "01 Empty Yaml Frontmatter",
			markdown: []string{
				"---",
				"---",
				"# Lorem impsum",
			},
			expected: []string{
				"<h1>Lorem impsum</h1>",
				""},
		},
		{
			name: "02 Sample Yaml Frontmatter",
			markdown: []string{
				"---",
				"postId: \"class-helpers-intro\"",
				"language: pl",
				"date: 2026-03-13",
				"---",
				"Impsum dolor"},
			expected: []string{
				"<p>Impsum dolor</p>",
				""},
		},
		{
			name: "03 Mid-document delimiter is preserved",
			markdown: []string{
				"Lorem impsum",
				"---",
				"Dolor sit amet",
			},
			expected: []string{
				"<p>Lorem impsum</p>",
				"<p>---</p>",
				"<p>Dolor sit amet</p>",
				""},
		},
		{
			name: "04 Unterminated leading delimiter is preserved",
			markdown: []string{
				"---",
				"postId: \"class-helpers-intro\"",
				"Impsum dolor",
			},
			expected: []string{
				"<p>---</p>",
				"<p>postId: &quot;class-helpers-intro&quot;</p>",
				"<p>Impsum dolor</p>",
				""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected, markdown := tt.toString(indentHtmlWith4Spaces)
			td.Cmp(t, GenerateHtmlBody(markdown), expected)
		})
	}
}

// ---------------------------------------------------------------------------
// Edge cases and malformed input
// ---------------------------------------------------------------------------

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		markdown       string
		shouldNotCrash bool
	}{
		{
			name:           "01 Malformed link",
			markdown:       "[link without URL]()",
			shouldNotCrash: true,
		},
		{
			name:           "02 Nested code blocks",
			markdown:       "```\n```\nnested\n```\n```",
			shouldNotCrash: true,
		},
		{
			name:           "03 Very long line",
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
			GenerateHtmlBody(tt.markdown)
		})
	}
}

// ---------------------------------------------------------------------------
// INTEGRATION TESTS
// ---------------------------------------------------------------------------

func TestConvertWithTemplate(t *testing.T) {
	markdown := "# Hello World\n\nGenerate HTML page"
	template := "<html><head><title>{{ .Title }}</title></head><body>{{ .Content }}</body></html>"
	title := "TestABC"

	result, err := ConvertMarkdownToHTML(markdown, template, title)

	expected := "<html><head><title>TestABC</title></head><body><h1>Hello World</h1>\n\n<p>Generate HTML page</p>\n</body></html>"
	td.Cmp(t, err, nil)
	td.Cmp(t, result, expected)
}

func TestConvertWithTemplateUsesFrontMatterMetadata(t *testing.T) {
	markdown := `---
title: "Front Matter Title"
description: "A document description"
date: 2026-03-13
author: "Bogdan Polak"
language: pl
coverImage: "cover.png"
coverImageCaption: "Cover caption"
---

# Hello World

Generate HTML page`
	template := `<html><head><title>{{ .Title }}</title><meta name="description" content="{{ .Description }}"></head><body><span class="date">{{ .Date }}</span><span class="author">{{ .Author }}</span><span class="language">{{ .Language }}</span><img src="{{ .CoverImage }}" alt="{{ .CoverImageCaption }}"><footer>{{ .PageFooter }}</footer><article>{{ .Content }}</article></body></html>`

	result, err := ConvertMarkdownToHTML(markdown, template, "")

	td.Cmp(t, err, nil)
	td.Cmp(t, result, `<html><head><title>Front Matter Title</title><meta name="description" content="A document description"></head><body><span class="date">2026-03-13</span><span class="author">Bogdan Polak</span><span class="language">pl</span><img src="cover.png" alt="Cover caption"><footer></footer><article><h1>Hello World</h1>

<p>Generate HTML page</p>
</article></body></html>`)
}

func TestConvertWithTemplateTitlePrecedence(t *testing.T) {
	tests := []struct {
		name          string
		title         string
		markdown      string
		expectedTitle string
	}{
		{
			name:          "cli title overrides front matter title",
			title:         "CLI Title",
			markdown:      "---\ntitle: \"Front Matter Title\"\n---\n\n# Hello",
			expectedTitle: "CLI Title",
		},
		{
			name:          "front matter title is used when cli title is empty",
			title:         "",
			markdown:      "---\ntitle: \"Front Matter Title\"\n---\n\n# Hello",
			expectedTitle: "Front Matter Title",
		},
		{
			name:          "default title is used when cli title and front matter title are empty",
			title:         "",
			markdown:      "# Hello",
			expectedTitle: "Converted Document",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertMarkdownToHTML(tt.markdown, "<title>{{ .Title }}</title><article>{{ .Content }}</article>", tt.title)

			td.Cmp(t, err, nil)
			td.Cmp(t, result, td.All(
				td.Contains("<title>"+tt.expectedTitle+"</title>"),
					td.Contains("<article>"),
					td.Contains("<h1>Hello</h1>"),
				td.Not(td.Contains("---")),
			))
		})
	}
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

	result := GenerateHtmlBody(markdown)

	result = strings.ReplaceAll(result, "\n", "│")
	td.Cmp(t, result, td.All(
		td.Contains("<h1>Main Title</h1>"),
		td.Contains("<p>This is an introduction with a <a href=\"https://example.com\">link</a> and some <code>inline code</code>.</p>"),
		td.Contains("<h2>Section 1</h2>"),
		td.Contains("<p>Some text with <code>npm install</code> command.</p>"),
		td.Contains("<ul>│    <li>Install dependencies</li>│    <li>Run tests</li>│    <li>Deploy</li>│</ul>"),
		td.Contains("<h2>Section 2</h2>"),
		td.Contains("<h3>Nested subsection</h3>"),
		td.Contains("<h4>Level four - Nested subsection</h4>"),
		td.Contains(
			strings.Join([]string{
				"<ol>",
				"    <li>First step</li>",
				"    <li>Second step with <code>code</code>",
				"        <section class=\"code\">",
				"        <pre><code>func main() {",
				"    fmt.Println(&quot;Hello&quot;)",
				"}</code></pre>",
				"        </section>",
				"    </li>",
				"</ol>",
			}, "│")),
		td.Contains("<p>Visit <a href=\"https://github.com\">https://github.com</a> for more info.</p>"),
	))
}
