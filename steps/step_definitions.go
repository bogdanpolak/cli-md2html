package steps

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

// Context holds the test state for BDD scenarios
type Context struct {
	Title      string
	Template   string
	Markdown   string
	ResultHTML string

	// CLI testing context
	Files         map[string]string // Mock filesystem: filename -> content
	StdinContent  string            // Content to pipe to stdin
	CommandOutput string            // Captured stdout from command
	CommandError  string            // Captured stderr from command
	ExitCode      int               // Exit code from command
	LastFile      string            // Last file that was checked/created
}

// ConvertMarkdownToHTML is the conversion function (will be called from main package)
var ConvertMarkdownToHTML func(markdown, template, title string) (string, error)

// Helpers
func assertTextEqual(actual, expected string) error {
	if actual == expected {
		return nil
	}

	return fmt.Errorf("HTML content mismatch.\nExpected: %s\nActual: %s",
		expected, actual)
}

func assertTextContains(actual, expected string) error {
	if !strings.Contains(actual, expected) {
		return fmt.Errorf("expected '%s' to contain '%s'", actual, expected)
	}
	return nil
}

// Step Implementations

func (c *Context) GivenIHaveATemplateContent(templateDoc *godog.DocString) error {
	c.Template = templateDoc.Content
	return nil
}

func (c *Context) GivenIHaveATitleWithATemplateContent(title string, templateDoc *godog.DocString) error {
	c.Title = title
	c.Template = templateDoc.Content
	return nil
}

func (c *Context) GivenIHaveMarkdownContent(markdownDoc *godog.DocString) error {
	c.Markdown = markdownDoc.Content
	return nil
}

func (c *Context) WhenIConvertItToHTML() error {
	if c.Template == "" {
		c.Title = "Test Document"
		c.Template = `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>{{ .Title }}</title></head>
<body>
{{ .Content }}
</body>
</html>`
	}

	// Use the conversion function
	html, err := ConvertMarkdownToHTML(c.Markdown, c.Template, c.Title)
	if err != nil {
		return err
	}

	c.ResultHTML = html
	return nil
}

func (c *Context) ThenIShouldGetHtmlContent(htmlDoc *godog.DocString) error {
	return assertTextEqual(c.ResultHTML, htmlDoc.Content)
}

func (c *Context) ThenIShouldGetHtmlContaining(expected string) error {
	return assertTextContains(c.ResultHTML, expected)
}

// CLI Testing Step Implementations

func (c *Context) GivenIHaveAMarkdownFileWithContent(filename, content string) error {
	c.Files[filename] = content
	return nil
}

func (c *Context) GivenIHaveMarkdownContentString(content string) error {
	c.StdinContent = content
	return nil
}

func (c *Context) GivenIHaveATemplateFileWithContent(filename string, content *godog.DocString) error {
	c.Files[filename] = content.Content
	return nil
}

func (c *Context) WhenIRunTheCommand(command string) error {
	// Parse the command string to extract arguments
	args := parseCommand(command)

	// Skip the command name if it's "md2html"
	if len(args) > 0 && args[0] == "md2html" {
		args = args[1:]
	}

	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "md2html-test-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Write mock files to temp directory
	for filename, content := range c.Files {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filename, err)
		}
	}

	// Build the command - use absolute path to binary
	wd, _ := os.Getwd()
	binaryPath := filepath.Join(wd, "bin", "md2html")
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = tempDir

	// Set up stdin if we have content to pipe
	if c.StdinContent != "" {
		cmd.Stdin = strings.NewReader(c.StdinContent)
	}

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err = cmd.Run()

	// Store results
	c.CommandOutput = stdout.String()
	c.CommandError = stderr.String()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			c.ExitCode = exitErr.ExitCode()
		} else {
			c.ExitCode = 1
		}
	} else {
		c.ExitCode = 0
	}

	// Read any output files that were created
	// First, check for files that were expected to be created (in the Files map)
	for filename := range c.Files {
		if strings.HasSuffix(filename, ".html") {
			filePath := filepath.Join(tempDir, filename)
			if content, err := os.ReadFile(filePath); err == nil {
				c.Files[filename] = string(content)
			}
		}
	}

	// Also check for any .html files that might have been created
	entries, err := os.ReadDir(tempDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".html") {
				filename := entry.Name()
				if _, exists := c.Files[filename]; !exists {
					filePath := filepath.Join(tempDir, filename)
					if content, err := os.ReadFile(filePath); err == nil {
						c.Files[filename] = string(content)
					}
				}
			}
		}
	}

	return nil
}

func (c *Context) WhenIPipeTheContentToMd2html() error {
	// This is similar to running a command but specifically for piping content
	return c.WhenIRunTheCommand("./md2html")
}

func (c *Context) ThenAFileShouldBeCreated(filename string) error {
	c.LastFile = filename
	if _, exists := c.Files[filename]; !exists {
		return fmt.Errorf("file %s was not created", filename)
	}
	return nil
}

func (c *Context) ThenTheFileShouldContain(expectedContent string) error {
	filename := c.LastFile
	if filename == "" {
		return fmt.Errorf("no file was previously checked for creation")
	}
	actualContent, exists := c.Files[filename]
	if !exists {
		return fmt.Errorf("file %s does not exist", filename)
	}
	if !strings.Contains(actualContent, expectedContent) {
		return fmt.Errorf("file %s content does not contain expected content.\nExpected to contain: %s\nActual content: %s", filename, expectedContent, actualContent)
	}
	return nil
}

func (c *Context) ThenIShouldGetHtmlOutputContaining(expected string) error {
	if !strings.Contains(c.CommandOutput, expected) {
		return fmt.Errorf("expected output to contain '%s', but got: %s", expected, c.CommandOutput)
	}
	return nil
}

func (c *Context) ThenTheHtmlOutputShouldContainATitle(expectedTitle string) error {
	expected := fmt.Sprintf("<title>%s</title>", expectedTitle)
	return c.ThenIShouldGetHtmlOutputContaining(expected)
}

func (c *Context) ThenTheHtmlOutputShouldContain(expected string) error {
	return c.ThenIShouldGetHtmlOutputContaining(expected)
}

func (c *Context) ThenIShouldSeeHelpTextContaining(expected string) error {
	output := c.CommandOutput + c.CommandError
	if !strings.Contains(output, expected) {
		return fmt.Errorf("expected help text to contain '%s', but got: %s", expected, output)
	}
	return nil
}

func (c *Context) ThenIShouldGetAnErrorMessage() error {
	if c.CommandError == "" && c.ExitCode == 0 {
		return fmt.Errorf("expected an error message, but command succeeded with output: %s", c.CommandOutput)
	}
	return nil
}

func (c *Context) ThenTheCommandShouldExitWithCode1() error {
	if c.ExitCode != 1 {
		return fmt.Errorf("expected exit code 1, but got %d. Output: %s, Error: %s", c.ExitCode, c.CommandOutput, c.CommandError)
	}
	return nil
}

func (c *Context) ThenIShouldGetAnErrorMessageAboutTemplateParsing() error {
	errorMsg := c.CommandError + c.CommandOutput
	if !strings.Contains(errorMsg, "template") {
		return fmt.Errorf("expected template parsing error, but got: %q", errorMsg)
	}
	return nil
}

// Helper function to parse command string into arguments
func parseCommand(command string) []string {
	// Simple parsing - split by spaces, but handle quoted strings
	var args []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(command); i++ {
		char := command[i]

		switch {
		case !inQuotes && (char == '"' || char == '\''):
			inQuotes = true
			quoteChar = char
		case inQuotes && char == quoteChar:
			inQuotes = false
			quoteChar = 0
		case !inQuotes && char == ' ':
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(char)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

// InitializeScenario registers all step definitions
func InitializeScenario(ctx *godog.ScenarioContext) {
	// Create a new context for each scenario
	scenarioContext := &Context{
		Files: make(map[string]string),
	}

	// Reset context before each scenario
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		scenarioContext.Files = make(map[string]string)
		scenarioContext.StdinContent = ""
		scenarioContext.CommandOutput = ""
		scenarioContext.CommandError = ""
		scenarioContext.ExitCode = 0
		scenarioContext.LastFile = ""
		return ctx, nil
	})

	// Register Given steps
	ctx.Given(`^I have a template content:$`, scenarioContext.GivenIHaveATemplateContent)
	ctx.Given(`^I have a title "([^"]*)" with a template content:$`, scenarioContext.GivenIHaveATitleWithATemplateContent)
	ctx.Given(`^I have a markdown content:$`, scenarioContext.GivenIHaveMarkdownContent)

	// Register When steps
	ctx.When(`^I convert it to HTML$`, scenarioContext.WhenIConvertItToHTML)

	// Register Then steps
	ctx.Then(`^I should get HTML content:$`, scenarioContext.ThenIShouldGetHtmlContent)
	ctx.Then(`^I should get HTML containing "([^"]*)"$`, scenarioContext.ThenIShouldGetHtmlContaining)

	// Register CLI testing steps
	ctx.Given(`^I have a markdown file "([^"]*)" with content "([^"]*)"$`, scenarioContext.GivenIHaveAMarkdownFileWithContent)
	ctx.Given(`^I have markdown content "([^"]*)"$`, scenarioContext.GivenIHaveMarkdownContentString)
	ctx.Given(`^I have a template file "([^"]*)" with content:$`, scenarioContext.GivenIHaveATemplateFileWithContent)
	ctx.When(`^I run the command "([^"]*)"$`, scenarioContext.WhenIRunTheCommand)
	ctx.When(`^I pipe the content to md2html$`, scenarioContext.WhenIPipeTheContentToMd2html)
	ctx.Then(`^a file "([^"]*)" should be created$`, scenarioContext.ThenAFileShouldBeCreated)
	ctx.Then(`^the file should contain "([^"]*)"$`, scenarioContext.ThenTheFileShouldContain)
	ctx.Then(`^I should get HTML output containing "([^"]*)"$`, scenarioContext.ThenIShouldGetHtmlOutputContaining)
	ctx.Then(`^the HTML output should contain a title "([^"]*)"$`, scenarioContext.ThenTheHtmlOutputShouldContainATitle)
	ctx.Then(`^the HTML output should contain "([^"]*)"$`, scenarioContext.ThenTheHtmlOutputShouldContain)
	ctx.Then(`^I should see help text containing "([^"]*)"$`, scenarioContext.ThenIShouldSeeHelpTextContaining)
	ctx.Then(`^I should get an error message$`, scenarioContext.ThenIShouldGetAnErrorMessage)
	ctx.Then(`^the command should exit with code 1$`, scenarioContext.ThenTheCommandShouldExitWithCode1)
	ctx.Then(`^I should get an error message about template parsing$`, scenarioContext.ThenIShouldGetAnErrorMessageAboutTemplateParsing)
}
