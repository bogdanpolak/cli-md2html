package main

import (
	"strings"
	"testing"
)

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

func TestEscapeHTML(t *testing.T) {
	input := `<script>alert("test")</script> & "quotes"`
	expected := `&lt;script&gt;alert(&quot;test&quot;)&lt;/script&gt; &amp; &quot;quotes&quot;`

	result := escapeHTML(input)

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
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
