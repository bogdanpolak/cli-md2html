package main

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

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

var testContext *md2htmlContext

type md2htmlContext struct {
	title    string
	template string
	markdown string
	// ---
	resultHtml string
}

func givenIHaveATemplateContent(templateDoc *godog.DocString) error {
	testContext = &md2htmlContext{
		template: templateDoc.Content}
	return nil
}

func givenIHaveATitleWithATemplateContent(title string, templateDoc *godog.DocString) error {
	testContext = &md2htmlContext{
		title:    title,
		template: templateDoc.Content}
	return nil
}

func givenIHaveMarkdownContentDocString(markdownDoc *godog.DocString) error {
	if testContext == nil {
		testContext = &md2htmlContext{}
	}
	testContext.markdown = markdownDoc.Content
	return nil
}

func whenIConvertItToHTML() error {
	if testContext == nil {
		return godog.ErrPending
	}
	if testContext.template == "" {
		testContext.title = "Test Document"
		testContext.template = `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>{{ .Title }}</title></head>
<body>
{{ .Content }}
</body>
</html>`
	}

	// Use the existing conversion function with the stored template
	html, err := ConvertMarkdownToHTML(testContext.markdown, testContext.template, testContext.title)
	if err != nil {
		return err
	}

	testContext.resultHtml = html
	return nil
}

func thenIShouldGetHtmlContent(htmlDoc *godog.DocString) error {
	if testContext == nil {
		return godog.ErrPending
	}
	return assertTextEqual(testContext.resultHtml, htmlDoc.Content)
}

func thenIShouldGetHtmlContaining(expected string) error {
	if testContext == nil {
		return godog.ErrPending
	}
	return assertTextContains(testContext.resultHtml, expected)
}

var quietMode bool = true

func stepNotImplementedNoParams() error {
	if quietMode {
		return godog.ErrSkip
	}
	return fmt.Errorf("Step not implemented")
}
func stepNotImplementedParams1(param1 string) error {
	if quietMode {
		return godog.ErrPending
	}
	return fmt.Errorf("Step not implemented: %s", param1)
}
func stepNotImplementedParams2(param1, param2 string) error {
	if quietMode {
		return godog.ErrPending
	}
	return fmt.Errorf("Step not implemented: %s, %s", param1, param2)
}
func stepNotImplementedParams1AndDocString(param1 string, docString *godog.DocString) error {
	if quietMode {
		return godog.ErrPending
	}
	return fmt.Errorf("Step not implemented: %s, %s", param1, docString.Content)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	// Reset test context before each scenario
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		testContext = nil
		return ctx, nil
	})

	// whenIConvertMarkupToHTML(markdownDoc *godog.DocString) error {

	// ctx.Given()
	ctx.Given(`^I have a template content:$`, givenIHaveATemplateContent)
	ctx.Given(`^I have a title "([^"]*)" with a template content:$`, givenIHaveATitleWithATemplateContent)
	ctx.Given(`^I have a markdown content:$`, givenIHaveMarkdownContentDocString)
	//
	ctx.When(`^I convert it to HTML$`, whenIConvertItToHTML)
	// ----
	ctx.Step(`^I should get HTML content$`, thenIShouldGetHtmlContent)
	ctx.Step(`^I should get HTML containing "([^"]*)"$`, thenIShouldGetHtmlContaining)
	// ---
	ctx.Given(`^I should get HTML containing "([^"]*)"$`, stepNotImplementedParams1)
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
	ctx.Then(`^the command should exit with code 1$`, stepNotImplementedNoParams)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
