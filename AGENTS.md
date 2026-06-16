# Koito Proxy

<!-- CODEGRAPH_START -->
## CodeGraph

In repositories indexed by CodeGraph (a `.codegraph/` directory exists at the repo root), reach for it BEFORE grep/find or reading files when you need to understand or locate code:

- **MCP tools** (when available): `codegraph_explore` answers most code questions in one call — the relevant symbols' verbatim source plus the call paths between them. `codegraph_node` returns one symbol's source + callers, or reads a whole file with line numbers. If the tools are listed but deferred, load them by name via tool search.
- **Shell** (always works): `codegraph explore "<symbol names or question>"` and `codegraph node <symbol-or-file>` print the same output.

If there is no `.codegraph/` directory, skip CodeGraph entirely — indexing is the user's decision.
<!-- CODEGRAPH_END -->

## Project

A metadata-correction transparent proxy for Koito that intercepts music scrobble requests, applies user-defined rules, and forwards to the upstream Koito service.

Built with Go, Gin, GORM, and SQLite.

## Architecture

Handler → Service → Repository

- Handlers: parse input, call service, return output.
- Services: business logic, rule engine sync.
- Repositories: GORM-backed persistence.
- Do not access repositories from handlers.

## Testing

- Ginkgo/Gomega for BDD-style tests.
- `httptest` for HTTP handler tests.
- Mock interfaces (repository), not implementations.
- Test success paths, failure paths, and edge cases.
