Feature: CLI Interface
  As a user of the md2html command-line tool
  I want to use various input and output options
  So that I can integrate it into my workflow

  Scenario: CLI 001 Convert file to file
    Given I have a markdown file "input.md" with content "# Hello"
    When I run the command "md2html -input input.md -output output.html"
    Then a file "output.html" should be created
    And the file should contain "<h1>Hello</h1>"

  Scenario: CLI 002 Convert from stdin to stdout
    Given I have markdown content "# From Stdin"
    When I pipe the content to md2html
    Then I should get HTML output containing "<h1>From Stdin</h1>"

  Scenario: CLI 003 Convert file to stdout
    Given I have a markdown file "test.md" with content "## Test Header"
    When I run the command "md2html -input test.md"
    Then I should get HTML output containing "<h2>Test Header</h2>"

  Scenario: CLI 004 Use custom title
    Given I have a markdown file "doc.md" with content "# Document"
    When I run the command "md2html -input doc.md -title "My Custom Title""
    Then the HTML output should contain a title "My Custom Title"

  Scenario: CLI 005 Use custom template file
    Given I have a markdown file "content.md" with content "# Content"
    And I have a template file "template.html" with content:
      """
      <html><body>{{.Content}}</body></html>
      """
    When I run the command "md2html -input content.md -template template.html"
    Then the HTML output should contain "<html><body>"
    And the HTML output should contain "<h1>Content</h1>"

  Scenario: CLI 006 Show help information
    When I run the command "md2html -help"
    Then I should see help text containing "Usage: md2html"
    And I should see help text containing "-input"
    And I should see help text containing "-output"
    And I should see help text containing "-template"
    And I should see help text containing "-title"

  Scenario: CLI 007 Handle missing input file
    When I run the command "md2html -input nonexistent.md"
    Then I should get an error message
    And the command should exit with code 1

  Scenario: CLI 008 Handle invalid template file
    Given I have a markdown file "valid.md" with content "# Test"
    When I run the command "md2html -input valid.md -template invalid-template.html"
    Then I should get an error message about template parsing
    And the command should exit with code 1