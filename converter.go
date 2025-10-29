package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

// ConvertMarkdownToHTML converts markdown to HTML using a template file
func ConvertMarkdownToHTML(markdown string, templateText string, title string) (string, error) {

	// Parse template
	template, err := template.New("document").Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	// Convert markdown to HTML content (without the full HTML structure)
	htmlContent := GenerateHtmlBodyInternalContent(markdown)

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
	err = template.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}

type TemplateData struct {
	Title   string
	Content string
}

// converts markdown to HTML content (main converter function)
func GenerateHtmlBodyInternalContent(markdown string) string {
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
		lineType, processedLine := processSingleLine(line)

		// Handle list grouping
		switch lineType {
		case "ol":
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
		case "ul":
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
		default:
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

func processSingleLine(line string) (string, string) {
	trimmed := strings.TrimSpace(line)

	// Empty lines
	if trimmed == "" {
		return "empty", ""
	}

	// Headers
	if strings.HasPrefix(trimmed, "#### ") {
		content := strings.TrimPrefix(trimmed, "#### ")
		return "h4", fmt.Sprintf("<h4>%s</h4>", processInlineElements(content))
	}
	if strings.HasPrefix(trimmed, "### ") {
		content := strings.TrimPrefix(trimmed, "### ")
		return "h3", fmt.Sprintf("<h3>%s</h3>", processInlineElements(content))
	}
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
