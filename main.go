package main

import (
	"flag"
	"fmt"
	"os"
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

	err := ConvertMarkdown(*inputFile, *outputFile, *templateFile, *title)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
