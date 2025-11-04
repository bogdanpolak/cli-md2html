package main

import (
	"testing"

	"github.com/bogdanpolak/cli-md2html/steps"
	"github.com/cucumber/godog"
)

// Initialize the conversion function in the steps package
func init() {
	steps.ConvertMarkdownToHTML = ConvertMarkdownToHTML
}

func Test_BDD_TemplateProcessing(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: steps.InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/template_processing.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func Test_BDD_CommandLine(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: steps.InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/cli_interface.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
