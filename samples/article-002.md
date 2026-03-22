---
postId: safe-refactoring-parameter-object
language: en
---

## Refactoring Step by Step

1. **Create the `game` object.**
   - Create a a new function that returns a `game` object, and test it.
    ```js
    test("createGame returns expected shape", () => {
      const game = createGame(["Luke", "Leia"]);

      expect(game.players).toHaveLength(2);
    });
    ```
