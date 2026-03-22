---
postId: safe-refactoring-parameter-object
language: en
---

# Bug

**Panic on Empty Lines in List-Nested Fenced Code Blocks**


# Example 1 - Empty Lines in Nested Ordered List

1. **Update the tests.**
    - Update tests.
    - Evaluate the size of the migration.
    - If available, enable automatic test execution in watch mode.
    ```js
    test("movePlayer - Obi-Wan passing Start and buys location", () => {
        const game = createTestGame();

        movePlayer_new(game);

        assert.equal(money, 1640);
    });
    ```
2. **Migrate call sites in production code.**
   - Update every call
   - Again, if the migration is large, keep both variants and migrate gradually, running tests after each change.

# Example 2 - Unordered List with Fenced Code Block

- update code:
```js
movePlayer_new(game);

assert.equal(money, 1640);
});
```
- migrate production code.
