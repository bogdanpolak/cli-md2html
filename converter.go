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
	htmlContent := GenerateHtmlBody(markdown)

	// Use default title if no title provided
	if title == "" {
		title = "Converted Document"
	}

	// Execute template
	var buf bytes.Buffer
	err = template.Execute(&buf, TemplateData{
		Title:   title,
		Content: htmlContent,
	})
	if err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}

type TemplateData struct {
	Title   string
	Content string
}

// Level information for the two-pass list processing
type ListLevel struct {
	Depth    int    // Number of spaces in source
	ListType string // "ul" or "ol"
}

func isCodeLine(ln string) bool {
	trimmed := strings.TrimSpace(ln)
	return trimmed == "```"
}

func isListLine(ln string) bool {
	if strings.HasPrefix(ln, "- ") {
		return true
	}
	matched, _ := regexp.MatchString(`^\d+\.\s`, ln)
	return matched
}

func isInsideListBlock(ln string, insideCodeBlock *bool) bool {
	trimmed := strings.TrimSpace(ln)
	if isCodeLine(trimmed) {
		*insideCodeBlock = !*insideCodeBlock
		return true
	}
	if *insideCodeBlock {
		return true
	}
	if isListLine(trimmed) || trimmed == "" {
		return true
	}
	return false
}

// converts markdown to HTML content (main converter function)
func GenerateHtmlBody(markdown string) string {
	var result strings.Builder

	lines := strings.Split(markdown, "\n")

	lineIdx := 0
	for lineIdx < len(lines) {
		currentLine := lines[lineIdx]

		// Handle multiline code blocks (```)
		if isCodeLine(currentLine) {
			newIdx, codeBlock := processCodeBlock(lineIdx, lines, "")
			result.WriteString(codeBlock)
			lineIdx = newIdx
			continue
		}

		// Check if this line starts a list block
		if isListLine(strings.TrimSpace(currentLine)) {
			listBlock := []string{}
			isInsideCode := false
			for lineIdx < len(lines) && isInsideListBlock(lines[lineIdx], &isInsideCode) {
				listBlock = append(listBlock, lines[lineIdx])
				lineIdx++
			}
			listHTML := processListBlock(listBlock)
			result.WriteString(listHTML)
			continue
		}

		// Skip empty lines. Add single empty HTML line only
		if currentLine == "" {
			lineIdx++
			// Skip consecutive empty lines
			for lineIdx < len(lines) && strings.TrimSpace(lines[lineIdx]) == "" {
				lineIdx++
			}
			result.WriteString("\n")
			continue
		}

		// Process single non-list line
		processedLine := processSingleLine(currentLine)
		if processedLine != "" {
			result.WriteString(processedLine)
			result.WriteString("\n")
		}

		lineIdx++
	}

	return result.String()
}

func processCodeBlock(lineIdx int, lines []string, indentation string) (int, string) {
	var codeBlock strings.Builder
	ln := lines[lineIdx]
	depth := getLineDepth(ln)
	lineIdx++
	startIdx := lineIdx
	for lineIdx < len(lines) && !isCodeLine(lines[lineIdx]) {
		lineIdx++
	}
	if startIdx < lineIdx {
		codeBlock.WriteString(indentation + "<section class=\"code\">\n" + indentation + "<pre><code>")
		for idx := startIdx; idx < lineIdx; idx++ {
			codeBlock.WriteString(escapeHTML(lines[idx][depth:]))
			if idx < lineIdx-1 {
				codeBlock.WriteString("\n")
			}
		}
		codeBlock.WriteString("</code></pre>\n" + indentation + "</section>\n")
	}

	lineIdx++ // Skip closing ```
	return lineIdx, codeBlock.String()
}

func getLineDepth(ln string) int {
	depth := 0
	for _, char := range ln {
		switch char {
		case ' ':
			depth++
		case '\t':
			depth += 4
		default:
			return depth
		}
	}
	return depth
}

// Process a block of list lines using the two-pass algorithm
func processListBlock(lines []string) string {
	if len(lines) == 0 {
		return ""
	}

	var result strings.Builder
	openTags := []string{} // Track open list tags for proper closing

	lineIdx := 0
	level := 0
	for lineIdx < len(lines) {
		currentLine := lines[lineIdx]
		trimmed := strings.TrimSpace(currentLine)
		currentDepth := getLineDepth(currentLine)

		if trimmed == "" {
			lineIdx++
			continue
		}

		if lineIdx == 0 {
			listType := "ol"
			if strings.HasPrefix(trimmed, "- ") {
				listType = "ul"
			}

			// Calculate indentation for block elements (ul/ol): 0, 8, 16, ...
			blockIndent := strings.Repeat(" ", level*8)
			tag := fmt.Sprintf("<%s>", listType)
			result.WriteString(blockIndent + tag + "\n")
			openTags = append(openTags, listType)
		}

		// Extract and process list item content
		var content string
		if after, ok := strings.CutPrefix(trimmed, "- "); ok {
			content = after
		} else {
			// Ordered list - remove number and dot
			re := regexp.MustCompile(`^\d+\.\s(.*)`)
			matches := re.FindStringSubmatch(trimmed)
			if len(matches) > 1 {
				content = matches[1]
			}
		}

		// Generate <li> with proper indentation: 4, 12, 20, ...
		lineIndent := strings.Repeat(" ", level*8+4)
		result.WriteString(fmt.Sprintf("%s<li>%s", lineIndent, processInlineElements(content)))

		// Move to next line and skip empty lines
		lineIdx++
		for lineIdx < len(lines) && strings.TrimSpace(lines[lineIdx]) == "" {
			lineIdx++
		}

		// Look ahead to see if next item is deeper (for nested lists)
		nextDepth := -1
		nextIsCode := false
		if lineIdx < len(lines) {
			// skip empty lines
			line := lines[lineIdx]
			isNextListLine := isListLine(strings.TrimSpace(line))
			nextIsCode = isCodeLine(line)
			nextDepth = getLineDepth(line)
			if nextIsCode || !isNextListLine {
				nextDepth = currentDepth
			}
		}

		if nextIsCode {
			blockIndent := strings.Repeat(" ", level*8+8)
			newIdx, codeBlock := processCodeBlock(lineIdx, lines, blockIndent)
			result.WriteString("\n" + codeBlock + lineIndent)
			lineIdx = newIdx
		}

		isNextLineIsDeeperList := nextDepth >= 0 && nextDepth > currentDepth
		isNextLineIsShallowerList := nextDepth >= 0 && level >= 0 && nextDepth < currentDepth
		if isNextLineIsDeeperList {
			level++
			listType := "ol"
			if strings.HasPrefix(strings.TrimSpace(lines[lineIdx]), "- ") {
				listType = "ul"
			}

			// Calculate indentation for block elements (ul/ol): 0, 8, 16, ...
			blockIndent := strings.Repeat(" ", level*8)
			tag := fmt.Sprintf("<%s>", listType)
			result.WriteString("\n" + blockIndent + tag + "\n")
			openTags = append(openTags, listType)
		} else if isNextLineIsShallowerList {
			result.WriteString("</li>\n")

			// Close nested list
			if len(openTags) > 0 {
				listType := openTags[len(openTags)-1]
				blockIndent := strings.Repeat(" ", level*8)
				result.WriteString(blockIndent + fmt.Sprintf("</%s>", listType) + "\n")
				openTags = openTags[:len(openTags)-1]
			}
			// Close the parent <li>
			level--
			blockIndent := strings.Repeat(" ", level*8+4)
			result.WriteString(blockIndent + "</li>\n")
		} else if level >= 0 {
			// Same level - close previous <li>
			result.WriteString("</li>\n")
		}

	}

	// Close all remaining open lists and list items
	for i := len(openTags) - 1; i >= 0; i-- {
		listType := openTags[i]
		blockIndent := strings.Repeat(" ", i*8)
		result.WriteString(blockIndent + fmt.Sprintf("</%s>", listType) + "\n")

		// If this list is nested (not the outermost), close the <li> that contains it
		if i > 0 {
			blockIndent := strings.Repeat(" ", (i-1)*8+4)
			result.WriteString(blockIndent + "</li>\n")
		}
	}

	return result.String()
}

func processSingleLine(line string) string {
	trimmed := strings.TrimSpace(line)

	// Empty lines
	if trimmed == "" {
		return ""
	}

	// Headers
	if strings.HasPrefix(trimmed, "#### ") {
		content := strings.TrimPrefix(trimmed, "#### ")
		return fmt.Sprintf("<h4>%s</h4>", processInlineElements(content))
	}
	if strings.HasPrefix(trimmed, "### ") {
		content := strings.TrimPrefix(trimmed, "### ")
		return fmt.Sprintf("<h3>%s</h3>", processInlineElements(content))
	}
	if strings.HasPrefix(trimmed, "## ") {
		content := strings.TrimPrefix(trimmed, "## ")
		return fmt.Sprintf("<h2>%s</h2>", processInlineElements(content))
	}
	if strings.HasPrefix(trimmed, "# ") {
		content := strings.TrimPrefix(trimmed, "# ")
		return fmt.Sprintf("<h1>%s</h1>", processInlineElements(content))
	}

	// Regular paragraphs
	return fmt.Sprintf("<p>%s</p>", processInlineElements(trimmed))
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
