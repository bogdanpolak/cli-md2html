package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

func main() {
	var inputFile = flag.String("input", "", "Input Markdown file (required)")
	var outputFile = flag.String("output", "", "Output HTML file (stdout if not specified)")
	var templateFile = flag.String("template", "", "HTML template file with %title% and %content% placeholders")
	var title = flag.String("title", "", "Title for the HTML document")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Usage: md2html -input <markdown-file> [-output <html-file>] [-template <template-file>] [-title <title>]")
		fmt.Println("  -input     Input Markdown file (required)")
		fmt.Println("  -output    Output HTML file (stdout if not specified)")
		fmt.Println("  -template  HTML template file with {{.Title}} and {{.Content}} placeholders")
		fmt.Println("  -title     Title for the HTML document")
		os.Exit(1)
	}

	convertMarkdown(*inputFile, *outputFile, *templateFile, *title)
}

func convertMarkdown(inputFile, outputFile, templateFile, title string) {
	// Read input file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Convert markdown to HTML
	var html string
	if templateFile != "" {
		html = markdownToHTMLWithTemplate(string(content), templateFile, title)
	} else {
		html = markdownToHTML(string(content))
	}

	// Write output
	if outputFile != "" {
		err = os.WriteFile(outputFile, []byte(html), 0644)
		if err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("HTML written to %s\n", outputFile)
	} else {
		fmt.Print(html)
	}
}

func markdownToHTML(markdown string) string {
	var result strings.Builder

	result.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"UTF-8\">\n<title>Converted Document</title>\n</head>\n<body>\n")

	lines := strings.Split(markdown, "\n")
	inCodeBlock := false
	codeBlockContent := strings.Builder{}
	inOrderedList := false
	inUnorderedList := false

	for i, line := range lines {
		// Handle code blocks (```)
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			// Close any open lists
			if inOrderedList {
				result.WriteString("</ol>\n")
				inOrderedList = false
			}
			if inUnorderedList {
				result.WriteString("</ul>\n")
				inUnorderedList = false
			}

			if inCodeBlock {
				// End of code block
				result.WriteString("<div class=\"source-code\">\n<pre><code>")
				result.WriteString(escapeHTML(codeBlockContent.String()))
				result.WriteString("</code></pre>\n</div>\n")
				codeBlockContent.Reset()
				inCodeBlock = false
			} else {
				// Start of code block
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			if codeBlockContent.Len() > 0 {
				codeBlockContent.WriteString("\n")
			}
			codeBlockContent.WriteString(line)
			continue
		}

		// Process regular lines
		lineType, processedLine := processLineWithType(line)

		// Handle list grouping
		if lineType == "ol" {
			if !inOrderedList {
				if inUnorderedList {
					result.WriteString("</ul>\n")
					inUnorderedList = false
				}
				result.WriteString("<ol>\n")
				inOrderedList = true
			}
			result.WriteString(processedLine)
			result.WriteString("\n")
		} else if lineType == "ul" {
			if !inUnorderedList {
				if inOrderedList {
					result.WriteString("</ol>\n")
					inOrderedList = false
				}
				result.WriteString("<ul>\n")
				inUnorderedList = true
			}
			result.WriteString(processedLine)
			result.WriteString("\n")
		} else {
			// Close any open lists
			if inOrderedList {
				result.WriteString("</ol>\n")
				inOrderedList = false
			}
			if inUnorderedList {
				result.WriteString("</ul>\n")
				inUnorderedList = false
			}

			if processedLine != "" {
				result.WriteString(processedLine)
				result.WriteString("\n")
			} else if i < len(lines)-1 && lines[i+1] != "" {
				// Add empty line only if next line is not empty (avoid double spacing)
				result.WriteString("\n")
			}
		}
	}

	// Close any remaining open lists
	if inOrderedList {
		result.WriteString("</ol>\n")
	}
	if inUnorderedList {
		result.WriteString("</ul>\n")
	}

	result.WriteString("</body>\n</html>")
	return result.String()
}

type TemplateData struct {
	Title   string
	Content string
}

func markdownToHTMLWithTemplate(markdown, templateFile, title string) string {
	// Read template file
	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		fmt.Printf("Error reading template file: %v\n", err)
		os.Exit(1)
	}

	// Parse template
	tmpl, err := template.New("document").Parse(string(templateContent))
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		os.Exit(1)
	}

	// Convert markdown to HTML content (without the full HTML structure)
	htmlContent := markdownToHTMLContent(markdown)

	// Use default title if no title provided
	if title == "" {
		title = "Converted Document"
	}

	// Prepare template data
	data := TemplateData{
		Title:   title,
		Content: htmlContent,
	}

	// Execute template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		os.Exit(1)
	}

	return buf.String()
}

func markdownToHTMLContent(markdown string) string {
	var result strings.Builder

	lines := strings.Split(markdown, "\n")
	inCodeBlock := false
	codeBlockContent := strings.Builder{}
	inOrderedList := false
	inUnorderedList := false

	for i, line := range lines {
		// Handle code blocks (```)
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			// Close any open lists
			if inOrderedList {
				result.WriteString("</ol>\n")
				inOrderedList = false
			}
			if inUnorderedList {
				result.WriteString("</ul>\n")
				inUnorderedList = false
			}

			if inCodeBlock {
				// End of code block
				result.WriteString("<div class=\"source-code\">\n<pre><code>")
				result.WriteString(escapeHTML(codeBlockContent.String()))
				result.WriteString("</code></pre>\n</div>\n")
				codeBlockContent.Reset()
				inCodeBlock = false
			} else {
				// Start of code block
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			if codeBlockContent.Len() > 0 {
				codeBlockContent.WriteString("\n")
			}
			codeBlockContent.WriteString(line)
			continue
		}

		// Process regular lines
		lineType, processedLine := processLineWithType(line)

		// Handle list grouping
		if lineType == "ol" {
			if !inOrderedList {
				if inUnorderedList {
					result.WriteString("</ul>\n")
					inUnorderedList = false
				}
				result.WriteString("<ol>\n")
				inOrderedList = true
			}
			result.WriteString(processedLine)
			result.WriteString("\n")
		} else if lineType == "ul" {
			if !inUnorderedList {
				if inOrderedList {
					result.WriteString("</ol>\n")
					inOrderedList = false
				}
				result.WriteString("<ul>\n")
				inUnorderedList = true
			}
			result.WriteString(processedLine)
			result.WriteString("\n")
		} else {
			// Close any open lists
			if inOrderedList {
				result.WriteString("</ol>\n")
				inOrderedList = false
			}
			if inUnorderedList {
				result.WriteString("</ul>\n")
				inUnorderedList = false
			}

			if processedLine != "" {
				result.WriteString(processedLine)
				result.WriteString("\n")
			} else if i < len(lines)-1 && lines[i+1] != "" {
				// Add empty line only if next line is not empty (avoid double spacing)
				result.WriteString("\n")
			}
		}
	}

	// Close any remaining open lists
	if inOrderedList {
		result.WriteString("</ol>\n")
	}
	if inUnorderedList {
		result.WriteString("</ul>\n")
	}

	return result.String()
}

func processLineWithType(line string) (string, string) {
	trimmed := strings.TrimSpace(line)

	// Empty lines
	if trimmed == "" {
		return "empty", ""
	}

	// Headers
	if strings.HasPrefix(trimmed, "## ") {
		content := strings.TrimPrefix(trimmed, "## ")
		return "h2", fmt.Sprintf("<h2>%s</h2>", processInlineElements(content))
	}
	if strings.HasPrefix(trimmed, "# ") {
		content := strings.TrimPrefix(trimmed, "# ")
		return "h1", fmt.Sprintf("<h1>%s</h1>", processInlineElements(content))
	}

	// Ordered lists
	if matched, _ := regexp.MatchString(`^\d+\.\s`, trimmed); matched {
		re := regexp.MustCompile(`^\d+\.\s(.*)`)
		matches := re.FindStringSubmatch(trimmed)
		if len(matches) > 1 {
			return "ol", fmt.Sprintf("    <li>%s</li>", processInlineElements(matches[1]))
		}
	}

	// Unordered lists (with proper indentation handling)
	if strings.HasPrefix(trimmed, "- ") {
		content := strings.TrimPrefix(trimmed, "- ")
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if indent > 4 { // Nested list - for now treat as regular list item
			return "ul", fmt.Sprintf("    <li>%s</li>", processInlineElements(content))
		}
		return "ul", fmt.Sprintf("    <li>%s</li>", processInlineElements(content))
	}

	// Regular paragraphs
	return "p", fmt.Sprintf("<p>%s</p>", processInlineElements(trimmed))
}

func processInlineElements(text string) string {
	// Process inline elements BEFORE escaping HTML

	// Inline code - extract and protect from escaping
	codeMap := make(map[string]string)
	codeCounter := 0
	re := regexp.MustCompile("`([^`]+)`")
	text = re.ReplaceAllStringFunc(text, func(match string) string {
		code := strings.Trim(match, "`")
		placeholder := fmt.Sprintf("__CODE_PLACEHOLDER_%d__", codeCounter)
		codeMap[placeholder] = fmt.Sprintf("<code>%s</code>", escapeHTML(code))
		codeCounter++
		return placeholder
	})

	// Links (before auto-links to avoid double processing)
	linkMap := make(map[string]string)
	linkCounter := 0
	re = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = re.ReplaceAllStringFunc(text, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) >= 3 {
			placeholder := fmt.Sprintf("__LINK_PLACEHOLDER_%d__", linkCounter)
			linkMap[placeholder] = fmt.Sprintf("<a href=\"%s\">%s</a>", escapeHTML(parts[2]), escapeHTML(parts[1]))
			linkCounter++
			return placeholder
		}
		return match
	})

	// Auto-links (standalone URLs)
	re = regexp.MustCompile(`(https?://[^\s\)<]+)`)
	text = re.ReplaceAllStringFunc(text, func(match string) string {
		// Don't replace if it's already in a placeholder (already processed as a link)
		if strings.Contains(match, "PLACEHOLDER") {
			return match
		}
		placeholder := fmt.Sprintf("__AUTOLINK_PLACEHOLDER_%d__", linkCounter)
		linkMap[placeholder] = fmt.Sprintf("<a href=\"%s\">%s</a>", escapeHTML(match), escapeHTML(match))
		linkCounter++
		return placeholder
	})

	// Now escape the remaining HTML in regular text
	text = escapeHTML(text)

	// Restore code placeholders
	for placeholder, codeHTML := range codeMap {
		text = strings.ReplaceAll(text, placeholder, codeHTML)
	}

	// Restore link placeholders
	for placeholder, linkHTML := range linkMap {
		text = strings.ReplaceAll(text, placeholder, linkHTML)
	}

	return text
}

func escapeHTML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	text = strings.ReplaceAll(text, "'", "&#39;")
	return text
}
