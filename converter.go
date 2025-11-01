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
	return ln == "```"
}

func isListLine(ln string) bool {
	s := strings.Trim(string(ln), " ")
	if strings.HasPrefix(s, "- ") {
		return true
	}
	matched, _ := regexp.MatchString(`^\d+\.\s`, s)
	return matched
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
			codeBlock := []string{}
			lineIdx++
			for lineIdx < len(lines) && !isCodeLine(lines[lineIdx]) {
				codeBlock = append(codeBlock, lines[lineIdx])
				lineIdx++
			}
			result.WriteString("<section class=\"code\">\n<pre><code>")
			result.WriteString(escapeHTML(strings.Join(codeBlock, "\n")))
			result.WriteString("</code></pre>\n</section>\n")
			lineIdx++ // Skip closing ``` if it exists
			continue
		}

		// Check if this line starts a list block
		if isListLine(currentLine) {
			listBlock := []string{currentLine}
			lineIdx++
			for lineIdx < len(lines) && (isListLine(lines[lineIdx])) {
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

// Process a block of list lines using the two-pass algorithm
func processListBlock(lines []string) string {
	if len(lines) == 0 {
		return ""
	}

	// First pass: detect indentation depths and list types
	levels := []ListLevel{}
	depthToIndex := make(map[int]int)

	for _, line := range lines {
		ln := strings.TrimSpace(line)
		if !isListLine(ln) {
			continue // Skip empty lines
		}

		// Calculate depth (tabs count as 4 spaces)
		depth := 0
		for _, char := range line {
			if char == ' ' {
				depth++
			} else if char == '\t' {
				depth += 4
			} else {
				break
			}
		}

		// Check if we've seen this depth before
		if _, exists := depthToIndex[depth]; !exists {
			// New depth level - determine list type from first occurrence
			trimmed := strings.TrimSpace(line)
			listType := "ol" // default
			if strings.HasPrefix(trimmed, "- ") {
				listType = "ul"
			}

			// Add to levels and create mapping
			depthToIndex[depth] = len(levels)
			levels = append(levels, ListLevel{Depth: depth, ListType: listType})
		}
	}

	// Sort levels by depth
	for i := 0; i < len(levels)-1; i++ {
		for j := i + 1; j < len(levels); j++ {
			if levels[i].Depth > levels[j].Depth {
				levels[i], levels[j] = levels[j], levels[i]
			}
		}
	}

	// Rebuild depth to index mapping after sorting
	depthToIndex = make(map[int]int)
	for i, level := range levels {
		depthToIndex[level.Depth] = i
	}

	// Second pass: generate HTML
	var result strings.Builder
	currentLevel := -1
	openTags := []string{} // Track open list tags for proper closing

	// Collect only the list lines for processing
	listItems := []string{}
	for _, line := range lines {
		if isListLine(line) {
			listItems = append(listItems, line)
		}
	}

	for i, line := range listItems {
		// Calculate depth
		depth := 0
	depthLoop:
		for _, char := range line {
			switch char {
			case ' ':
				depth++
			case '\t':
				depth += 4
			default:
				break depthLoop
			}
		}

		levelIndex := depthToIndex[depth]

		// Look ahead to see if next item is deeper (for nested lists)
		nextIsDeeper := false
		if i+1 < len(listItems) {
			nextDepth := 0
		nextDepthLoop:
			for _, char := range listItems[i+1] {
				switch char {
				case ' ':
					nextDepth++
				case '\t':
					nextDepth += 4
				default:
					break nextDepthLoop
				}
			}
			nextLevelIndex := depthToIndex[nextDepth]
			nextIsDeeper = nextLevelIndex > levelIndex
		}

		// Handle level changes
		if levelIndex > currentLevel {
			// Going deeper - open new nested lists
			for currentLevel < levelIndex {
				currentLevel++
				level := levels[currentLevel]

				// Calculate indentation for block elements (ul/ol): 0, 8, 16, ...
				blockIndent := strings.Repeat(" ", currentLevel*8)

				tag := fmt.Sprintf("<%s>", level.ListType)
				result.WriteString(blockIndent + tag + "\n")
				openTags = append(openTags, level.ListType)
			}
		} else if levelIndex < currentLevel {
			// Going shallower - close current item first, then close nested lists
			result.WriteString("</li>\n")

			// Close nested lists
			for currentLevel > levelIndex {
				if len(openTags) > 0 {
					listType := openTags[len(openTags)-1]
					blockIndent := strings.Repeat(" ", currentLevel*8)
					result.WriteString(blockIndent + fmt.Sprintf("</%s>", listType) + "\n")
					openTags = openTags[:len(openTags)-1]
				}
				currentLevel--
			}

			// Close the parent <li>
			liIndent := strings.Repeat(" ", levelIndex*8+4)
			result.WriteString(liIndent + "</li>\n")
		} else if currentLevel >= 0 {
			// Same level - close previous <li>
			result.WriteString("</li>\n")
		}

		// Extract and process list item content
		trimmed := strings.TrimSpace(line)
		var content string
		if strings.HasPrefix(trimmed, "- ") {
			content = strings.TrimPrefix(trimmed, "- ")
		} else {
			// Ordered list - remove number and dot
			re := regexp.MustCompile(`^\d+\.\s(.*)`)
			matches := re.FindStringSubmatch(trimmed)
			if len(matches) > 1 {
				content = matches[1]
			}
		}

		// Generate <li> with proper indentation: 4, 12, 20, ...
		liIndent := strings.Repeat(" ", levelIndex*8+4)
		result.WriteString(fmt.Sprintf("%s<li>%s", liIndent, processInlineElements(content)))

		// Only add newline here - closing </li> will be handled by next iteration or at the end
		if nextIsDeeper {
			result.WriteString("\n") // Just a newline before nested content
		} else if i == len(listItems)-1 {
			// Last item - close it
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
			liIndent := strings.Repeat(" ", (i-1)*8+4)
			result.WriteString(liIndent + "</li>\n")
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
