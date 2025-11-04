package steps

import (
	"context"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// Context holds the test state for BDD scenarios
type Context struct {
	Title      string
	Template   string
	Markdown   string
	ResultHTML string
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
		return fmt.Errorf("Expected '%s' to contain '%s'", actual, expected)
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

// Placeholder for unimplemented steps
var quietMode = true

func stepNotImplementedNoParams() error {
	if quietMode {
		return godog.ErrSkip
	}
	return fmt.Errorf("Step not implemented")
}

func stepNotImplementedParams1(param1 string) error {
	if quietMode {
		return godog.ErrSkip
	}
	return fmt.Errorf("Step not implemented: %s", param1)
}

func stepNotImplementedParams2(param1, param2 string) error {
	if quietMode {
		return godog.ErrSkip
	}
	return fmt.Errorf("Step not implemented: %s, %s", param1, param2)
}

func stepNotImplementedParams1AndDocString(param1 string, docString *godog.DocString) error {
	if quietMode {
		return godog.ErrSkip
	}
	return fmt.Errorf("Step not implemented: %s, %s", param1, docString.Content)
}

// InitializeScenario registers all step definitions
func InitializeScenario(ctx *godog.ScenarioContext) {
	// Create a new context for each scenario
	scenarioContext := &Context{}

	// Reset context before each scenario
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		scenarioContext = &Context{}
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

	// Register placeholder steps (not yet implemented)
	ctx.Given(`^I have a markdown file "([^"]*)" with content "([^"]*)"$`, stepNotImplementedParams2)
	ctx.Given(`^I have markdown content "([^"]*)"$`, stepNotImplementedParams1)
	ctx.Given(`^I have a template file "([^"]*)" with content:$`, stepNotImplementedParams1AndDocString)
	ctx.When(`^I run the command "([^"]*)"$`, stepNotImplementedParams1)
	ctx.When(`^I pipe the content to md2html$`, stepNotImplementedParams1)
	ctx.Then(`^a file "([^"]*)" should be created$`, stepNotImplementedParams1)
	ctx.Then(`^the file should contain "([^"]*)"$`, stepNotImplementedParams1)
	ctx.Then(`^I should get HTML output containing "([^"]*)"$`, stepNotImplementedParams1)
	ctx.Then(`^the HTML output should contain a title "([^"]*)"$`, stepNotImplementedParams1)
	ctx.Then(`^the HTML output should contain "([^"]*)"$`, stepNotImplementedParams1)
	ctx.Then(`^I should see help text containing "([^"]*)"$`, stepNotImplementedParams1)
	ctx.Then(`^I should get an error message$`, stepNotImplementedNoParams)
	ctx.Then(`^the command should exit with code 1$`, stepNotImplementedNoParams)
	ctx.Then(`^I should get an error message about template parsing$`, stepNotImplementedNoParams)
}
