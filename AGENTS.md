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

`go.mod` targets `go 1.18`; do not use language features newer than that without bumping the module version. There is no CI, linter, or Makefile configuration in the repo: the minimum validation is build, test, vet, and gofmt.

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
- Plans created with the plan tool are written to `~/.commandcode/plans/`. When a plan is finished, copy it into `docs/plans/` immediately, renamed to the `YYYY-MM-DD_PLAN_<slug>.md` convention:

  ```sh
  cp ~/.commandcode/plans/<plan-name>.md docs/plans/YYYY-MM-DD_PLAN_<slug>.md
  ```

- `README.md` is still the default GitLab template and does not describe the library; do not use it as a behavior reference.
- `CLAUDE.md` only contains a reference to this file.

## Skills

Skills live in `.agents/skills/`. Load one with the skill tool before starting work it covers. The tables below list each skill, when to use it, and what it is for.

### Specifications and planning

| Skill | Use when | Purpose |
|---|---|---|
| `brainstorming` | Before any creative work: new features, components, or behavior changes | Explores intent, requirements, and design before implementation. |
| `writing-specs` | Writing or revising specs in `docs/specs/` (FSD, TSD, DRD, diagrams) | Enforces separate document types, shared date/slug naming, and traceable requirement IDs. Not for briefs, plans, README, or code-only changes. |
| `writing-plans` | You have a spec or requirements for a multi-step task | Produces a step-by-step implementation plan before touching code. |
| `executing-plans` | Executing a written plan in a separate session with review checkpoints | Keeps execution aligned with the plan and adds review gates. |
| `writing-test-case` | Creating or revising TCDs, or defining acceptance and edge-case coverage | Defines verifiable positive, negative, and boundary cases with expected results. |

### Development workflow

| Skill | Use when | Purpose |
|---|---|---|
| `test-driven-development` | Implementing a feature or bugfix | Write the test first, watch it fail, then implement the minimal fix. |
| `systematic-debugging` | Any bug, test failure, or unexpected behavior | Investigates root cause before proposing fixes. |
| `using-git-worktrees` | Feature work that needs isolation from the current workspace | Creates a separate worktree with safety checks. |
| `dispatching-parallel-agents` | Two or more independent tasks with no shared state or sequential dependency | Runs the tasks concurrently. |
| `subagent-driven-development` | Executing a plan of independent tasks in the current session | Dispatches a subagent per task with review between tasks. |
| `karpathy-guidelines` | Writing, reviewing, or refactoring code | Reduces over-engineering and keeps changes surgical and verifiable. |

### Review and completion

| Skill | Use when | Purpose |
|---|---|---|
| `requesting-code-review` | Completing tasks, shipping a major feature, or preparing to merge | Verifies the work against requirements before merge. |
| `receiving-code-review` | Acting on review feedback, especially when it is unclear or questionable | Requires verification and technical rigor instead of blind agreement. |
| `verification-before-completion` | About to claim work is complete, fixed, or passing | Runs verification commands and confirms output before making success claims. |
| `finishing-a-development-branch` | Implementation is complete and all tests pass | Presents structured options for merge, PR, or cleanup. |

### Anti-slop suite

| Skill | Use when | Purpose |
|---|---|---|
| `antislop` | Always, as the core filter | Stops generic AI slop in code and output. |
| `antislop-code` | Writing or reviewing code comments | Removes generic AI-slop comments while keeping valuable ones; does not change code. |
| `antislop-copywriting` | Writing or editing prose such as headlines, tone, and CTAs | Applies anti-AI-writing patterns to copy. |
| `antislop-ui` | Building or editing any interface | Covers color, layout, components, and motion. |
| `antislop-human` | UI work involving accessibility | Covers contrast, keyboard, focus, and states; includes a contrast checker. |
| `antislop-layoutmobile` | Layouts that reflow from phone to desktop | Covers grids, overflow, and tap targets. |

### Project knowledge and tooling

| Skill | Use when | Purpose |
|---|---|---|
| `using-superpowers` | Starting a conversation or deciding which skill to use | Explains how to discover and invoke skills. |
| `graphify` | Questions about the codebase, architecture, file relationships, or the project knowledge graph | Builds, queries, or updates the project knowledge graph. |
| `writing-skills` | Creating, editing, or verifying a skill | Enforces skill structure and testing before deployment. |
| `deploy-server` | Deploying, releasing, rolling back, or restarting Rakkume apps on the systemd server | Build, transfer, and restart runbook for Rakkume apps; not applicable to this library. |

## Constraints

- Public v2 library: avoid breaking changes to existing signatures. Add new methods/types instead of changing existing ones.
- The only runtime dependency is `github.com/lib/pq`.
- Do not change the order in which placeholder parameters are registered without adjusting tests and callers, because the order of `GetSQLParameter()` is part of the usage contract.
