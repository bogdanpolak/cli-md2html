package main

import (
	"os/exec"
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

func buildCliBinary(t *testing.T) {
	t.Helper()

	cmd := exec.Command("go", "build", "-o", "bin/md2html", ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build md2html test binary: %v\n%s", err, output)
	}
}

func Test_BDD_CommandLine(t *testing.T) {
	buildCliBinary(t)

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
