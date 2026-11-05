# AGENTS.md

Operating guide for AI agents in the `sqlgo` repo. This document defines how work is done here. It does not duplicate the specifications. Behavior, contracts, and data have one source of truth: `docs/specs/`.

## 1. Ground rules

- Read this document and the relevant skill before changing code.
- Every piece of work follows `docs/specs/`. Product and technical detail is not restated here.
- Do not invent requirements, numbers, contracts, schemas, or behavior. If a specification does not exist, write it first with the `writing-specs` skill.
- Every change must trace to a requirement ID (`FR-xxx`) in the FSD and to its design in the TSD/DRD.
- If a change touches behavior, contracts, or data, the specification update belongs in the same plan. Do not defer specs as separate work.
- For non-trivial work, write the plan in `docs/plans/` first and wait for an explicit go-ahead before code execution.
- Agent tools may draft plans in `~/.commandcode/plans/`. That directory is untracked tool storage, not a deliverable. Copy the plan to `docs/plans/{yyyy-mm-dd}_{feature-name}.md` immediately after it is written, before implementation and without waiting to be asked. `docs/plans/` is the plan of record: when a copied plan changes, update it there in the same change.
- Commit and push documents (brief, spec, plan) before code execution so they can be reviewed.
- During execution, commit and push each finished item before starting the next one.
- An implementation gap is not an exception to the rules. Record the gap and close it against the target contract.

## 2. Project scope

This repo implements **sqlgo**, a chainable SQL query builder library for Go (module `gitlab.com/threetopia/sqlgo/v2`, package `sqlgo`). All code lives in the root package. The builder composes PostgreSQL statement text and keeps the `$N` placeholder list in sync; it does not open connections, execute queries, manage transactions, or run migrations.

Current scope is **PostgreSQL only**. `MySQLDialect` is defined in `sqlgo.go` but no MySQL builder branch exists. PostgreSQL-specific support covers `tsvector`/`tsquery` full-text search, `similarity` from `pg_trgm`, and vector distance from `pgvector`. The only runtime dependency is `github.com/lib/pq`.

Specification sets in `docs/specs/` are mandatory references:

| Set | Documents | Used for |
|---|---|---|
| `sqlgo` | `2026-11-04_FSD_sqlgo.md`, `2026-11-04_TSD_sqlgo.md`, `2026-11-04_DRD_sqlgo.md` | Query-building requirements, public contracts, internal design, parameter registry, PostgreSQL data requirements |
| `sqlgo` test cases | `../test-cases/2026-11-04_TCD_sqlgo.md` | Positive, negative, boundary, and acceptance coverage with expected results |

There is no DFD: the library owns no schema, owns no data store, and never executes SQL.

Repository layout:

| Path | Contents |
|---|---|
| `sqlgo.go` | `SQLGo` facade, `NewSQLGo()`, `BuildSQL()` composition order, dialect constants |
| `sqlgo_<clause>.go` | One file per clause: `insert`, `update`, `delete`, `select`, `from`, `join`, `where`, `values`, `group_by`, `order`, `offset_limit`, `schema` |
| `sqlgo_parameter.go` | Parameter registry and `$N` placeholder generator |
| `util.go` | Shared helpers such as `combineSQL` |
| `sqlgo_test.go` | SQL string and parameter construction tests |
| `docs/` | `brief/`, `plans/`, `specs/`, `test-cases/` with templates and their dated instances |
| `.agents/skills/` | Agent skills listed in section 3 |

Reading order:

1. Before changing a clause, read the `sqlgo` spec set.
2. Then read the related feature spec at `docs/specs/<date>_{FSD,TSD,DRD}_<slug>.md` and the matching test cases.
3. Adoption notes or gap notes inside a spec do not replace the target contract.
4. If a feature spec does not exist, create it with the `writing-specs` skill before writing code.

Do not copy sqlgo detail into this document. Link to `docs/specs/` so the two cannot drift apart.

## 3. Skills

Skills live in `.agents/skills/<name>/SKILL.md`. Read a skill file before using it. The "Required" column means the skill must not be skipped under the stated condition.

An agent can use more than one skill at the same time. Load every skill whose condition applies to the task, and combine them when their concerns overlap. For example, use `antislop` together with `antislop-copywriting` while writing a spec, or `test-driven-development` together with `systematic-debugging` while fixing a failing test. Using one skill does not exclude another.

### 3.1 Core: output and writing quality

| Skill | Required | When to use | What it is for |
|---|---|---|---|
| `antislop` | Yes, on every output | Any document, UI text, comment, or user-facing code you write | Core anti-slop filter: rejects technique without purpose, requires a written reason, keeps the result alive |
| `antislop-copywriting` | Yes, when writing proses | Documents, briefs, specs, plans, UI copy, and product text | Removes AI writing patterns, chooses specific words, preserves the author's voice |
| `antislop-code` | Yes, when touching comments | Writing or editing code comments | Removes slop comments, keeps comments that carry information, never changes the code |
| `antislop-ui` | Yes, when building UI | Pages, components, color, layout, and motion | Rejects default AI visuals and enforces hierarchy and visual identity |
| `antislop-human` | Yes, when the UI has interaction | Forms, navigation, color, states, and interactive controls | Protects contrast, keyboard access, focus, and states so the UI works for everyone |
| `antislop-layoutmobile` | Yes, when the layout must reflow | Layouts that run from phone to desktop | Enforces breakpoints, grids, overflow, tap targets, and mobile navigation |

### 3.2 Meta and process

| Skill | Required | When to use | What it is for |
|---|---|---|---|
| `using-superpowers` | Yes, at session start | Starting a new conversation | Decides how to find and invoke the right skill before responding |
| `brainstorming` | Yes, before creative work | Before building a feature, component, or behavior change | Explores intent, requirements, and design before implementation |
| `writing-specs` | Yes, when writing specs | Writing or revising specs in `docs/specs/` | Produces FSD, TSD, DRD, DFD, and diagrams with consistent structure and naming |
| `writing-test-case` | Yes, when creating a TCD | Writing a Test Case Document or agent evaluation scenarios | Defines positive, negative, acceptance coverage, and expected results |
| `writing-plans` | Yes, for multi-step tasks | After requirements or specs exist, before touching code | Produces an implementation plan that can be executed and reviewed |
| `executing-plans` | Conditional | Executing a written plan in a separate session with review checkpoints | Runs the plan in stages while keeping review at each checkpoint |
| `subagent-driven-development` | Yes, for approved plans | Executing a plan with independent tasks in the same session | Default method: controller, one implementer per task, spec review, then code-quality review |
| `dispatching-parallel-agents` | Conditional | Two or more independent tasks with no shared state | Runs tasks in parallel without waiting on each other |
| `using-git-worktrees` | Conditional | A feature needs isolation from the workspace, or before executing a plan | Creates an isolated worktree with directory selection and safety checks |

### 3.3 Quality and delivery

| Skill | Required | When to use | What it is for |
|---|---|---|---|
| `test-driven-development` | Yes, before implementation code | Working on a feature or bugfix | Writes the test first, watches it fail, then writes the minimal code |
| `systematic-debugging` | Yes, on any failure | A bug, failing test, or unexpected behavior appears, before proposing a fix | Finds the root cause before changing code |
| `verification-before-completion` | Yes, before claiming completion | Before saying work is done, committing, or opening a PR | Runs verification and shows output as evidence, not assumptions |
| `requesting-code-review` | Conditional | Finishing a task, a major feature, or before merging | Checks that the result matches the requirements |
| `receiving-code-review` | Yes, when receiving review | Receiving review feedback, before applying suggestions | Judges feedback on technical merit instead of agreeing or applying it blindly |
| `finishing-a-development-branch` | Conditional | Implementation is complete and all tests pass | Chooses the integration step: merge, PR, or cleanup |
| `karpathy-guidelines` | Yes, when touching code | Writing, reviewing, or refactoring code | Keeps changes small, avoids overcomplication, and states assumptions |

### 3.4 Domain

| Skill | Required | When to use | What it is for |
|---|---|---|---|
| `graphify` | Conditional, always through the wrapper | Architecture or file-relationship questions, or building and updating the knowledge graph | Answers codebase questions from the graph and keeps the graph current |
| `deploy-server` | Yes, when deploying | Deploying, releasing, rolling back, or restarting a Rakkume app on the systemd production server | Builds the binary, rsyncs, restarts the unit, runs migrations, and verifies health |
| `writing-skills` | Conditional | Creating, editing, or verifying a skill | Keeps agent skill structure and quality consistent |

## 4. Graphify

The `graphify` skill expects a `scripts/graphify` wrapper at the repo root. This repo does not currently have one. If the wrapper is absent, report it as a missing prerequisite; do not call the bare `graphify` executable, `python -m graphify`, or another project's wrapper, and do not create a wrapper unless requested.

If a wrapper is added later, these rules apply:

- Always run `scripts/graphify <subcommand> ...`.
- Never write the graph to `graphify-out/` or any other directory. The official output is `docs/graphify/`.
- If `docs/graphify/graph.json` already exists and the question is natural language, query it directly. Do not rebuild the graph.
- After a large code change, update the graph incrementally.
- If the wrapper is missing or unusable, report it as a missing prerequisite. Do not create a new wrapper without a request.

Examples once the wrapper exists:

```bash
scripts/graphify query "how does parameter deduplication work"
scripts/graphify path "SQLGoWhere" "SQLGoParameter"
scripts/graphify explain "SQLGoParameter"
scripts/graphify update .
scripts/graphify god-nodes --top 10
```

## 5. Workflow

1. Establish scope. Read `docs/specs/` as described in section 2.
2. Write a brief or plan in `docs/brief/` or `docs/plans/`.
3. If the plan was drafted in `~/.commandcode/plans/`, copy it to `docs/plans/` immediately after it is written, before implementation.
4. Include specification updates in the plan when the change touches behavior, contracts, or data.
5. Commit and push the documents, then wait for the execution go-ahead.
6. Execute with TDD. For an approved plan, use `subagent-driven-development`.
7. Run verification: `go build ./...`, `go test ./...`, `go vet ./...`, and `gofmt -l .`. `go.mod` targets `go 1.18`; do not use newer language features without bumping the module.
8. Review the result, then commit and push.
9. Update the graph with `scripts/graphify update .` when the change affects large-scale structure, if the wrapper exists.

Known baseline: `go test ./...` currently fails on `TestToTsQuery` and `TestToTsQueryPrepend`. Both expect `$2` while parameter deduplication returns `$1` for the same value. Do not assume this failure was caused by your change; fix the test expectation only when the task actually touches that behavior.

## 6. Prohibitions

- Do not copy `docs/specs/` content into this document or into another document that can drift.
- Do not add dependencies, abstractions, or configuration beyond the task's needs. The only runtime dependency is `github.com/lib/pq`.
- Do not change existing public v2 signatures; add new methods or types instead.
- Do not change the order in which placeholders are registered without adjusting tests and callers, because `GetSQLParameter()` order is part of the usage contract.
- Do not commit secrets, tokens, or credentials.
- Do not leave a dev server or background process running after testing.
