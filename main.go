package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var help = flag.Bool("help", false, "Show help information")
	var inputFile = flag.String("input", "", "Input Markdown file (stdin if not specified)")
	var outputFile = flag.String("output", "", "Output HTML file (stdout if not specified)")
	var templateFile = flag.String("template", "", "HTML template file with %title% and %content% placeholders (optional)")
	var title = flag.String("title", "", "Title for the HTML document (optional)")
	flag.Parse()

	if *help {
		fmt.Println("Usage: md2html -input <markdown-file> [-output <html-file>] [-template <template-file>] [-title <title>]")
		fmt.Println("  -input     Input Markdown file (stdin if not specified)")
		fmt.Println("  -output    Output HTML file (stdout if not specified)")
		fmt.Println("  -template  HTML template file with {{.Title}} and {{.Content}} placeholders (optional)")
		fmt.Println("  -title     Title for the HTML document")
		os.Exit(1)
	}

	err := ConvertMarkdown(*inputFile, *outputFile, *templateFile, *title)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func ConvertMarkdown(inputFile, outputFile, templateFile, title string) error {
	var content []byte
	var err error
	if inputFile == "" {
		content, err = io.ReadAll(os.Stdin)
	} else {
		content, err = os.ReadFile(inputFile)
	}
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	// Read or assign default template
	var html string
	var templateContent string
	if templateFile != "" {
		text, err := os.ReadFile(templateFile)
		if err != nil {
			return fmt.Errorf("error reading template file: %w", err)
		}
		templateContent = string(text)
	} else {
		templateContent = `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>{{ .Title }}</title></head>
<body>
{{ .Content }}
</body>
</html>`
	}

	html, err = ConvertMarkdownToHTML(string(content), templateContent, title)
	if err != nil {
		return err
	}

	// Write output
	if outputFile != "" {
		err = os.WriteFile(outputFile, []byte(html), 0644)
		if err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}
		fmt.Printf("HTML written to %s\n", outputFile)
	} else {
		fmt.Print(html)
	}

	return nil
}
