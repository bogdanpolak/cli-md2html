# MD2HTML - Simple Markdown to HTML Converter

A lightweight Go CLI tool that converts Markdown documents to simple HTML without CSS styling.

## Features

- ✅ **Headers** (`#` and `##` → `<h1>` and `<h2>`)
- ✅ **Code blocks** (``` → `<pre><code>`)
- ✅ **Inline code** (` → `<code>`)
- ✅ **Ordered lists** (`1.` → `<ol><li>`)
- ✅ **Unordered lists** (`-` → `<ul><li>`)
- ✅ **Links** (`[text](url)` and auto-detect URLs → `<a href="">`)
- ✅ **Paragraphs** (regular text → `<p>`)
- ✅ **List grouping** (consecutive list items are grouped properly)

## Installation

```bash
go build -o md2html main.go
```

## Usage

```bash
# Convert to HTML file (basic)
./md2html -input input.md -output output.html

# Convert with custom template
./md2html -input input.md -template template.html -title "My Document" -output output.html

# Convert to stdout
./md2html -input input.md

# Show help
./md2html
```

### Template Support

You can use a custom HTML template with Go template syntax:
- `{{.Title}}` - Replaced with the document title
- `{{.Content}}` - Replaced with the converted markdown content

**Example template:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container">
        {{.Content}}
    </div>
</body>
</html>
```

The template uses Go's standard `text/template` package, so you can use any Go template features like conditionals, loops, etc.

## Example

**Input Markdown:**
```markdown
## CLI `go-migrate` tool

1. Create serverless DB
    - https://console.neon.tech
    - Define environment variable with connection string

2. Create `001` migration
    ```
    migrate create -ext sql -dir ./migrations -seq "Create Sets table"
    ```
```

**Output HTML:**
```html
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Converted Document</title>
</head>
<body>
<h2>CLI <code>go-migrate</code> tool</h2>

<ol>
    <li>Create serverless DB</li>
</ol>
<ul>
    <li><a href="https://console.neon.tech">https://console.neon.tech</a></li>
    <li>Define environment variable with connection string</li>
</ul>

<ol>
    <li>Create <code>001</code> migration</li>
</ol>
<pre><code>migrate create -ext sql -dir ./migrations -seq "Create Sets table"</code></pre>
</body>
</html>
```

## Limitations

This is a simple converter focused on basic Markdown elements. It does not support:
- Tables
- Images
- Complex nested lists
- Bold/italic formatting
- Blockquotes
- Horizontal rules

## Notes

- The output HTML has no CSS styling - it's plain semantic HTML
- HTML characters are properly escaped
- URLs are automatically converted to clickable links
- Code blocks preserve formatting and indentation