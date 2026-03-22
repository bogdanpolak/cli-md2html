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
EXPECTED02='<div class="code">│<pre><code>npm install│npm start</code></pre>│</div>'
cat test02.md | ./bin/md2html | tr '\n' '│' | grep "$EXPECTED02" > /dev/null && echo "✅ OK "$N02 || echo "❌ Failed "$N02
rm test02.md

# ---------------------------------------------------------
T03_ID="[TEST 03]"
T03_NAME="$T03_ID Regression - Panic on empty lines in lists"
SAMPLE03_MD="./samples/bug007-article.md"
if (cat "$SAMPLE03_MD" > /dev/null); then
  cat "$SAMPLE03_MD" | ./bin/md2html | tr '\n' '│' | grep "panic: runtime error: slice bounds out of range" > /dev/null && echo "❌ Failed "$T03_NAME" - panic occurred" || echo "✅ OK "$T03_NAME" - no panic"
  EXPECTED03_Sample01='│            <li>If available, enable automatic test execution in watch mode.│                <div class="code">│                <pre><code>test(&quot;movePlayer - Obi-Wan passing Start and buys location&quot;, () =&gt; {│    const game = createTestGame();││    movePlayer_new(game);││    assert.equal(money, 1640);│});</code></pre>│                </div>│            </li>│'
  cat "$SAMPLE03_MD" | ./bin/md2html | tr '\n' '│' | grep "$EXPECTED03_Sample01" > /dev/null && echo "✅ OK "$T03_NAME" - Example 1" || echo "❌ Failed "$T03_NAME" - Example 1"
  EXPECTED03_Sample02='│    <li>update code:│        <div class="code">│        <pre><code>movePlayer_new(game);││assert.equal(money, 1640);│});</code></pre>│        </div>│    </li>│'
  cat "$SAMPLE03_MD" | ./bin/md2html | tr '\n' '│' | grep "$EXPECTED03_Sample02" > /dev/null && echo "✅ OK "$T03_NAME" - Example 2" || echo "❌ Failed "$T03_NAME" - Example 2"
else
  echo "❌ Failed "$T03_ID" - sample file missing: "$SAMPLE03_MD
fi 
