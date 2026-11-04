# AGENTS.md

Guidance for agents working in the `sqlgo` repository.

## Project Overview

`sqlgo` is an SQL query builder library for Go (module `gitlab.com/threetopia/sqlgo/v2`, package `sqlgo`). All code lives in the root package.

Current implementation scope: **PostgreSQL**. Placeholders use the `$N` format, parameters are passed through `github.com/lib/pq`, and PostgreSQL-specific builders are available: `tsvector`/`tsquery`, `similarity` (pg_trgm), and embedding (pgvector `<=>`). The `MySQLDialect`/`PostgreSQLDialect` constants exist in `sqlgo.go`, but there is no MySQL-specific builder branch yet.

The API is chain-based:

```go
sql := sqlgo.NewSQLGo().
    SQLSchema("schema").
    SQLSelect(sqlgo.SetSQLSelect("t.id", "id")).
    SQLFrom("table", "t").
    SQLWhere(sqlgo.SetSQLWhere("AND", "t.id", "=", 1))

query := sql.BuildSQL()                           // SQL string
args := sql.GetSQLGoParameter().GetSQLParameter() // []interface{}
```

## Commands

```sh
go build ./...
go test ./...
go vet ./...
gofmt -l .
```

`go.mod` targets `go 1.18`; do not use language features newer than that without bumping the module version. There is no CI, linter, or Makefile configuration in the repo — the minimum validation is build, test, vet, and gofmt.

## File Structure

| File | Contents |
|---|---|
| `sqlgo.go` | `SQLGo` facade, `NewSQLGo()`, `BuildSQL()` composition order, dialect constants |
| `sqlgo_<clause>.go` | One file per clause: `insert`, `update`, `delete`, `select`, `from`, `join`, `where`, `values`, `group_by`, `order`, `offset_limit`, `schema` |
| `sqlgo_parameter.go` | Parameter registry and placeholder generator |
| `util.go` | `combineSQL` and `hash` helpers |
| `sqlgo_test.go` | SQL and parameter construction tests |

## Architecture

- Each clause has a public `SQLGo<Clause>` interface, a private `sqlGo<Clause>` implementation, and implements `SQLGoBase` (`GetSQLGoParameter`, `BuildSQL`) plus the `SetSQLGoParameter` setter.
- There are two usage styles: prepend/chain methods (`SQLSelect(...)`, `SQLWhere(...)`) and one-at-a-time `SetSQL*` methods (`SetSQLSelect`, `SetSQLWhere`, ...). Both must produce identical SQL.
- The package also exposes value constructor functions, e.g. `SetSQLSelect`, `SetSQLWhere`, `SetSQLOrder`, `AddSQLGoValue`, `AddSQLGoValueArray`.
- The `sqlGo.BuildSQL()` facade composes clauses in a fixed order: insert → update → delete → select → from → join → where → values → group by → order → offset/limit. Empty clauses return an empty string and are skipped by `combineSQL`.
- Parameters are registered via `SQLGoParameter.SetSQLParameter`; identical values (SHA-1 hash of JSON + type) use the same placeholder number. The order of `GetSQLParameter()` is first-appearance order.
- The schema prefix (`SetSQLSchema`/`SQLSchema`) is only used by insert, update, delete, from, and join.

## Code Conventions

- Public interfaces and types use the `SQLGo` prefix; private types use the `sqlGo` prefix.
- Builder methods return the clause/facade interface so they can be chained.
- Do not insert user values directly into SQL strings. All dynamic values must go through `SetSQLParameter` and then `GetSQLParameterSign`.
- To add a new clause: create `sqlgo_<clause>.go` (interface, struct, constructor, `BuildSQL`), register the field and `SetSQLGo<Clause>` on `SQLGOModuleSetter`/`sqlGo` in `sqlgo.go`, then call that module in `sqlGo.BuildSQL()` according to SQL order.
- Follow the existing `fmt.Sprintf` pattern so spacing between clauses stays consistent; tests compare strings exactly.

## Testing

- Tests compare the SQL string exactly (including spaces) against query constants, then check `GetSQLGoParameter()`.
- Every builder change must run `go test ./...` and update the affected query constants.
- **Baseline status:** `go test ./...` currently fails on `TestToTsQuery` and `TestToTsQueryPrepend`. Both expect `$2`, while parameter dedup returns `$1` for the same value. Do not assume this failure is caused by your change; only fix the test/expectation if the task actually touches it.

## Documentation

- `docs/brief/`, `docs/plans/`, `docs/specs/`, `docs/test-cases/` contain templates. Duplicate a template and name it `YYYY-MM-DD_<TYPE>_<slug>.md`.
- `README.md` is still the default GitLab template and does not describe the library; do not use it as a behavior reference.
- `CLAUDE.md` only contains a reference to this file.

## Constraints

- Public v2 library: avoid breaking changes to existing signatures. Add new methods/types instead of changing existing ones.
- The only runtime dependency is `github.com/lib/pq`.
- Do not change the order in which placeholder parameters are registered without adjusting tests and callers, because the order of `GetSQLParameter()` is part of the usage contract.
