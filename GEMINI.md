# Koito Proxy's Agents MD

<!-- CODEGRAPH_START -->
## CodeGraph

In repositories indexed by CodeGraph (a `.codegraph/` directory exists at the repo root), reach for it BEFORE grep/find or reading files when you need to understand or locate code:

- **MCP tools** (when available): `codegraph_explore` answers most code questions in one call — the relevant symbols' verbatim source plus the call paths between them. `codegraph_node` returns one symbol's source + callers, or reads a whole file with line numbers. If the tools are listed but deferred, load them by name via tool search.
- **Shell** (always works): `codegraph explore "<symbol names or question>"` and `codegraph node <symbol-or-file>` print the same output.

If there is no `.codegraph/` directory, skip CodeGraph entirely — indexing is the user's decision.
<!-- CODEGRAPH_END -->

## Core Principles

Priority order:

1. Correctness
2. Context propagation
3. Testability
4. Readability
5. Simplicity
6. Performance

When uncertain, ask questions instead of making assumptions.

---

## Context Propagation

* Always propagate `context.Context`.
* Pass context through handlers, services, repositories, database calls (GORM), cache calls, and external API calls.
* Avoid `context.Background()` and `context.TODO()` in request flows.
* `context.Background()` is acceptable during application bootstrap and intentionally detached background workers.

---

## Go Style

* Prefer explicit code over magic.
* Prefer standard library before introducing dependencies.
* Avoid premature optimization.
* Avoid `make([]T, 0)` unless a capacity is provided or benchmarked performance justifies it.
* Prefer constructor dependency injection.
* Avoid package-level mutable state.

---

## Architecture

Preferred flow:

Gin Handler → Service → Repository

* Keep handlers thin. Use Gin for request/response handling.
* Business logic belongs in services (e.g. metadata transformation rules).
* Repositories handle persistence concerns (SQLite via GORM).
* Do not access repositories directly from handlers.

---

## Error Handling

* Wrap errors with contextual information.
* Return actionable errors.
* Do not silently ignore errors.

Example:

```go
return fmt.Errorf("get user: %w", err)
```

---

## Testing

All new business logic should include tests.

Preferred:

* Table-driven tests.
* `httptest` for HTTP handlers (Gin).
* Mock interfaces, not implementations.

Test:

* Success paths.
* Failure paths.
* Edge cases.

---

## Git

Use Conventional Commits.

Examples:

* feat(rules): add rule processing engine
* fix(proxy): handle nil request body
* refactor(repo): propagate context in GORM
* test(api): add Gin handler tests

Avoid vague commit messages.

---

## Generated Code Expectations

Generated code should:

* Compile.
* Include tests when appropriate.
* Propagate context correctly.
* Handle errors properly.
* Minimize dependencies.

---

## Clarification Policy

Do not invent:

* API contracts
* Database schemas
* Business rules
* Requirements

Ask for clarification when information is missing.

