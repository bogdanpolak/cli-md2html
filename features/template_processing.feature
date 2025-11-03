Feature: Template Processing
  As a user of the md2html CLI tool
  I want to use custom HTML templates
  So that I can control the output format and styling

  Scenario: Processing 001 Use default template
    Given I have a markdown content:
      """
      # Hello World
      """
    When I convert it to HTML
    Then I should get HTML containing "<h1>Hello World</h1>"
	
  Scenario: Processing 002 Use custom template with title
    Given I have a title "Generated document v1" with a template content:
      """
      <!DOCTYPE html>
      <html>
      <head><title>{{.Title}}</title></head>
      <body><article>{{.Content}}</article></body>
      </html>
      """
    And I have a markdown content:
      """
      This is a sample markdown line.
      """
    When I convert it to HTML
    Then I should get HTML content:
      """
      <!DOCTYPE html>
      <html>
      <head><title>Generated document v1</title></head>
      <body><article><p>This is a sample markdown line.</p></article></body>
      </html>
      """

  Scenario: Processing 003 Template with complex content
    Given I have a markdown content: 
      """
      # Hello World
      
      Introduction to the subject

      - Why point one
      - Act now

      This is the article content with **bold** text.
      """
    When I convert it to HTML
    Then I should get HTML containing "<h1>Hello World</h1>"
    And I should get HTML containing "<p>This is the article content with <strong>bold</strong> text.</p>"