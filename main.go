package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	var help = flag.Bool("help", false, "Show help information")
	var inputFile = flag.String("input", "", "Input Markdown file (stdin if not specified)")
	var outputFile = flag.String("output", "", "Output HTML file (stdout if not specified)")
	var templateFile = flag.String("template", "", "HTML template file with %title% and %content% placeholders (optional)")
	var title = flag.String("title", "", "Title for the HTML document (optional)")
	var preview = flag.Bool("preview", false, "Open converted HTML in default browser")
	flag.Parse()

	if *help {
		fmt.Println("Usage: md2html -input <markdown-file> [-output <html-file>] [-template <template-file>] [-title <title>] [-preview]")
		fmt.Println("  -input     Input Markdown file (stdin if not specified)")
		fmt.Println("  -output    Output HTML file (stdout if not specified)")
		fmt.Println("  -template  HTML template file with {{.Title}} and {{.Content}} placeholders (optional)")
		fmt.Println("  -title     Title for the HTML document")
		fmt.Println("  -preview   Open converted HTML in default browser")
		os.Exit(1)
	}

	err := ConvertMarkdown(*inputFile, *outputFile, *templateFile, *title, *preview)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func ConvertMarkdown(inputFile, outputFile, templateFile, title string, preview bool) error {
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
		templateContent = "<!DOCTYPE html>\n" +
			"<html>\n<head>\n  <meta charset=\"UTF-8\">\n  <title>{{ .Title }}</title>\n</head>\n" +
			"<body>{{ .Content }}</body>\n</html>"
	}

	html, err = ConvertMarkdownToHTML(string(content), templateContent, title)
	if err != nil {
		return err
	}

	if preview {
		// Create temporary file
		tempFile, err := os.CreateTemp("", "md2html-preview-*.html")
		if err != nil {
			return fmt.Errorf("error creating temp file: %w", err)
		}
		defer tempFile.Close()

		// Write HTML to temp file
		_, err = tempFile.WriteString(html)
		if err != nil {
			return fmt.Errorf("error writing to temp file: %w", err)
		}

		// Get absolute path for browser
		tempPath, err := filepath.Abs(tempFile.Name())
		if err != nil {
			return fmt.Errorf("error getting temp file path: %w", err)
		}

		// Open in default browser
		cmd := exec.Command("open", tempPath)
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("error opening browser: %w", err)
		}

		fmt.Printf("Preview opened in browser: %s\n", tempPath)
		return nil
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
