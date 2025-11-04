CURRENT_DIR=$(pwd)
echo "Running E2E tests in "$CURRENT_DIR
cd "$CURRENT_DIR"

# Compile the tool
go build -o bin/md2html

# ---------------------------------------------------------
N01="[TEST] Convert Title and 3 section document"
cat > test01.md << 'EOF'
# Long Test Document

This is a longer markdown document to test the piping functionality.

## Section 1

Some content here.

## Section 2

More content.

### Subsection

Even more content.

## Section 3

And finally, some more text.

This should be enough to test the head command.
EOF
EXPECTED01='<h3>Subsection</h3>││<p>Even more content.</p>││<h2>Section 3</h2>││<p>And finally, some more text.</p>'
cat test01.md | ./bin/md2html | tr '\n' '│' | grep "$EXPECTED01" > /dev/null && echo "✅ OK "$N01 || echo "❌ Failed "$N01
rm test01.md

# ---------------------------------------------------------
N02="[TEST] Convert Code Blocks"
cat > test02.md << 'EOF'
## Test Code Blocks

If you're using Go Modules and have a `go.mod` file in your project's root, you can import `clerk-sdk-go` directly:
```
import (
  "github.com/clerk/clerk-sdk-go/v2"
)
go get -u github.com/clerk/clerk-sdk-go/v2
```

And another block:

```
npm install
npm start
```
EOF
EXPECTED02='<section class="code">│<pre><code>npm install│npm start</code></pre>│</section>'
cat test02.md | ./bin/md2html | tr '\n' '│' | grep "$EXPECTED02" > /dev/null && echo "✅ OK "$N02 || echo "❌ Failed "$N02
rm test02.md
