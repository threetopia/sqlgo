# Plan: sqlgo Test Case Document (TCD)

## Goal

Create one Test Case Document for the whole `sqlgo` library in `docs/test-cases/`, covering every functional requirement in the existing specification set (`FR-001` to `FR-019`) plus the documented edge, error, security, and regression behavior.

No source, test, spec, or template file is modified. The task produces one documentation file. Automated test implementation is out of scope unless requested later.

## Decisions

- **Format:** one TCD for the whole library, not one file per clause. The spec set already uses one date and slug (`2026-11-04_sqlgo`), so the TCD reuses them.
- **Filename:** `docs/test-cases/2026-11-04_TCD_sqlgo.md`.
- **Language:** English, matching `2026-11-04_FSD_sqlgo.md`, `2026-11-04_TSD_sqlgo.md`, and `2026-11-04_DRD_sqlgo.md`. The TCD template is Indonesian, but the content follows the spec set and the project documentation preference.
- **Requirement source:** only the existing `FR-001` to `FR-019` IDs. No new requirement IDs are invented.
- **Test case IDs:** `TC-{DOMAIN}-{NNN}`. Domains are listed in the inventory below. IDs are sequential per domain and stable once written.
- **Status:** every test case starts as `Not Run`. No `Pass` claim is made without an executed and verified test.
- **Dedup conflict:** expected results follow `FR-006` (identical value and type pairs share one placeholder). The existing `TestToTsQuery` and `TestToTsQueryPrepend` expectation of `$2` is recorded as an open question and a regression case, not as the expected behavior.
- **Automation:** the TCD defines expected behavior and maps cases to existing Go tests. New automated tests are not written in this task.
- **Template section 11:** the template's usage guide is not copied into the TCD. The TCD contains sections 1 to 10 only.

## Files

Create:

- `docs/test-cases/2026-11-04_TCD_sqlgo.md`

Leave unchanged:

- `docs/test-cases/1970-01-31_TCD_test-case-template.md`
- `docs/specs/2026-11-04_FSD_sqlgo.md`, `2026-11-04_TSD_sqlgo.md`, `2026-11-04_DRD_sqlgo.md`
- All `sqlgo_*.go`, `sqlgo.go`, `util.go`, `sqlgo_test.go`
- `README.md`, `AGENTS.md`, `.agents/skills/**`

## Source references for factual accuracy

Read these before writing. Every expected SQL shape, placeholder, and parameter order must match the code.

- `sqlgo.go` - facade, clause methods, `BuildSQL()` composition order, module setters, dialect constants.
- `sqlgo_select.go` - projection, alias, `DISTINCT ON`, `ts_rank`, `similarity`, embedding, embedding array, subquery value.
- `sqlgo_from.go` - table, alias, column list, subquery source.
- `sqlgo_join.go` - join type casing, table or subquery, alias, condition, `ON` replacement for literal tables.
- `sqlgo_where.go` - connectors, `specialOperator` map, `IN`/`NOT IN`, `pq.Array`, `BETWEEN`, groups, `to_tsquery`, `similarity`, non-parameter conditions.
- `sqlgo_insert.go`, `sqlgo_update.go`, `sqlgo_delete.go`, `sqlgo_values.go` - write statements and `to_tsvector`.
- `sqlgo_group_by.go`, `sqlgo_order.go`, `sqlgo_offset_limit.go`, `sqlgo_schema.go` - grouping, ordering, pagination, schema prefix.
- `sqlgo_parameter.go` - hash, dedup map, ordinal assignment, `$-1` fallback, nil handling.
- `util.go` - `combineSQL` spacing and `hash`.
- `sqlgo_test.go` - exact SQL constants and the current active test names.
- `docs/specs/2026-11-04_FSD_sqlgo.md`, `2026-11-04_TSD_sqlgo.md`, `2026-11-04_DRD_sqlgo.md` - requirements, contracts, data rules, known limitations, open questions.

## TCD document structure

Follow the template layout and fill every placeholder with real data.

### 1. Document information

- Title: TCD: sqlgo query builder.
- Document type: TCD.
- Created and last updated: `2026-11-04`.
- Status: `Draft`.
- Feature or module: `sqlgo`, root Go package.
- Scope: SQL generation and parameter binding for SELECT, FROM, JOIN, WHERE, GROUP BY, ORDER BY, OFFSET/LIMIT, INSERT, UPDATE, DELETE, VALUES, schema prefix, and the PostgreSQL full-text, trigram, and vector features.
- Exclusions: query execution, connection management, transactions, migrations, ORM mapping, MySQL-specific builder, real PostgreSQL extension behavior, performance, and concurrency testing.
- Links: `../specs/2026-11-04_FSD_sqlgo.md`, `../specs/2026-11-04_TSD_sqlgo.md`, `../specs/2026-11-04_DRD_sqlgo.md`, `../plans/2026-11-04_PLAN_sqlgo-specs-and-readme.md`.
- Owner: Repository maintainer.
- Brief: none exists for `sqlgo`; state `none` rather than linking the template.

### 2. Test objectives and scope

- Verify exact SQL text and ordered parameter slices for every clause.
- Verify the two usage styles and the fixed composition order.
- Verify PostgreSQL feature constructors.
- Verify documented edge and error behavior.
- State preconditions: Go 1.18, no database connection, PostgreSQL `$N` placeholders, `lib/pq` available, tests run with the standard `testing` package.
- State that `pg_trgm` and `pgvector` are only relevant to SQL text, not to a live database in this TCD.

### 3. Requirements traceability

Map `FR-001` to `FR-019` to test case IDs. Identify the known partial coverage for `FR-003` (subquery join keyword) and `FR-006` (dedup conflict) and link them to section 9.

### 4. Test cases

Use the template's required fields for each case: ID, requirement, title, category, priority, preconditions, test data, steps, expected result, actual result, status, evidence. Keep shared preconditions at domain level to avoid repetition.

### 5. Coverage matrix

Count cases per category and mark not-applicable categories with a reason.

### 6. AI agent testing

Mark every category as not applicable. `sqlgo` contains no LLM, agent, tool calling, or retrieval logic. State the reason once instead of mechanically filling 11 sub-sections.

### 7. BDD scenarios

Include two Gherkin scenarios only as documentation, with no new framework:

- A SELECT with parameters builds the expected SQL and ordered argument slice.
- An INSERT with a repeated value and type uses one placeholder.

### 8. Automation strategy per layer

- Unit: Go `testing`, table-driven cases, exact string and slice assertions. This is the primary layer.
- Integration: not applicable, the library never opens a database connection.
- Acceptance/BDD: optional documentation only.
- Security: unit-level assertion that dynamic values never appear in the SQL text.
- LLM evaluation: not applicable.

### 9. Assumptions, open questions, and coverage gaps

Record at least:

- Dedup versus `TestToTsQuery` and `TestToTsQueryPrepend`.
- Subquery JOIN appends the condition with `WHERE` instead of `ON`.
- `nil` values render `$-1` and are unsupported.
- `SetSQLPage` must be called after the limit is set.
- Identifiers, aliases, operators, and language names are interpolated without quoting or validation.
- Hash collisions are not checked against the original value.
- No active automated coverage for OFFSET/LIMIT, standalone VALUES, DISTINCT ON, subquery FROM/JOIN, empty builder, schema first-wins, and nil handling.
- Performance, concurrency, live PostgreSQL extension behavior, and MySQL are not covered and why.

### 10. Definition of done checklist

Mark the drafting items as done and keep execution-related items unchecked. All statuses remain `Not Run`.

## Test case inventory

The exact wording is finalized while writing, but the scope and ID ranges are fixed here.

| Domain | IDs | Requirement | Cases | Key scenarios |
|---|---|---|---|---|
| SELECT | `TC-SELECT-001..004` | `FR-001`, `FR-015` | 4 | Single and multiple items, alias and empty alias, subquery projection, `DISTINCT ON`. |
| FROM | `TC-FROM-001..004` | `FR-002` | 4 | Table with alias, column list, no alias, subquery source. |
| JOIN | `TC-JOIN-001..005` | `FR-003` | 5 | Single and multiple joins, join type casing, schema prefix, subquery and literal conditions. |
| WHERE | `TC-WHERE-001..008` | `FR-004`, `FR-005` | 8 | Connector rules, special operator matrix, `IN`/`NOT IN`, `pq.Array`, `BETWEEN`, group and nested group, non-parameter condition, empty slice boundary. |
| PARAM | `TC-PARAM-001..007` | `FR-006` | 7 | Sequential placeholders, dedup, type separation, first-appearance order, nil, unknown lookup, empty builder. |
| SCHEMA | `TC-SCHEMA-001..003` | `FR-007` | 3 | Applied clauses, first-write-wins, no prefix on SELECT/WHERE/ORDER. |
| GROUP BY | `TC-GROUP-001..002` | `FR-008` | 2 | Single and multiple values, omitted when empty. |
| ORDER BY | `TC-ORDER-001..003` | `FR-009` | 3 | `ASC`/`DESC`, empty direction, embedding ordering. |
| OFFSET/LIMIT | `TC-LIMIT-001..005` | `FR-010` | 5 | Offset with limit, limit only, offset-only not emitted, page helper, page boundary and call order. |
| INSERT | `TC-INSERT-001..004` | `FR-011`, `FR-016` | 4 | Single and multiple rows, dedup across rows, `to_tsvector`, columns without values and missing columns. |
| UPDATE | `TC-UPDATE-001..002` | `FR-012`, `FR-016` | 2 | Assignments, `to_tsvector`, omitted when empty. |
| DELETE | `TC-DELETE-001..003` | `FR-013` | 3 | Table, schema prefix, facade `SQLDelete` forwarding WHERE. |
| VALUES | `TC-VALUES-001..002` | `FR-014` | 2 | Single and multiple rows, dedup. |
| FTS | `TC-FTS-001..002` | `FR-016` | 2 | `ts_rank` projection, `@@ to_tsquery` filter. |
| TRGM | `TC-TRGM-001..002` | `FR-017` | 2 | Similarity projection, similarity filter with threshold parameter. |
| VECTOR | `TC-VECTOR-001..004` | `FR-018` | 4 | Embedding projection with and without prefix, array projection, embedding ordering. |
| COMPOSE | `TC-COMPOSE-001..003` | `FR-019` | 3 | Prepend/chain versus `Set*` versus module setters, fixed order, empty clause spacing. |
| SECURITY | `TC-SEC-001..002` | `FR-004`, `FR-006` | 2 | Dynamic value stays out of SQL text, identifier interpolation risk documented. |
| REGRESSION/EDGE | `TC-REG-001..004` | `FR-006`, `FR-010`, `FR-019` | 4 | Dedup conflict baseline, pagination call order, builder reuse, invalid identifier passthrough. |

Total: about 65 cases. Merge table-driven variants for special operators and slice types into one case per group so the document stays readable.

## Existing automated coverage mapping

Reference these tests in the traceability and automation sections. Do not modify them.

- `TestSelectQueryPrepend`, `TestSelectQueryPipeline`, `TestSelectQuerySetByVariable` - SELECT, FROM, JOIN, WHERE, GROUP BY, ORDER BY, schema, two styles.
- `TestDeleteQueryPrepend`, `TestDeleteQueryPipeline` - DELETE and WHERE forwarding.
- `TestInsert`, `TestInsertToTsVector` - INSERT rows, dedup, `to_tsvector`.
- `TestUpdate`, `TestUpdateToTsVector`, `TestUpdateToTsVectorPrepend` - UPDATE and `to_tsvector`.
- `TestWhereGroupQueryPrepend`, `TestWhereGroupPipeline`, `TestBetween` - WHERE group and `BETWEEN`.
- `TestToTsQuery`, `TestToTsQueryPrepend`, `TestToTsQueryEmbedding`, `TestToTsQueryEmbeddingPrepend` - full-text and embedding ordering.
- `TestToSelectEmbedding`, `TestToSelectEmbeddingArray`, `TestSelectSimilarity` - embedding projection and similarity.

Mark cases that have no active test as manual or documented-only in the automation section.

## Execution steps

1. Re-read the FSD, TSD, DRD, source files, and `sqlgo_test.go` (already read during planning; confirm again before writing).
2. Write `docs/test-cases/2026-11-04_TCD_sqlgo.md` with sections 1 to 10 from the template.
3. Fill the document information, scope, exclusions, related links, and owner.
4. Write the traceability table covering `FR-001` to `FR-019`.
5. Write the summary table and detailed cases for the inventory above.
6. Write the coverage matrix, AI-agent not-applicable section, BDD scenarios, and automation strategy.
7. Write the assumptions, open questions, and coverage gaps section with the known conflicts.
8. Write the definition of done checklist with execution items unchecked.
9. Cross-check every SQL shape and placeholder against the source and `sqlgo_test.go`.
10. Report the created file, the requirement coverage, and the open questions. Do not fix code or tests.

## Verification

- File exists at `docs/test-cases/2026-11-04_TCD_sqlgo.md` and follows `{yyyy-mm-dd}_TCD_{slug}.md`.
- The date and slug match `2026-11-04_sqlgo`.
- Relative links to FSD, TSD, DRD, and the plan resolve.
- Every `FR-001` to `FR-019` appears in the traceability table and maps to at least one test case.
- Every test case ID is unique and follows `TC-{DOMAIN}-{NNN}`.
- Every case has a specific, observable expected result.
- All statuses are `Not Run`; no unverified `Pass` claim exists.
- Categories that do not apply are marked `n/a` with a reason, not filled mechanically.
- The template file and all source, test, and spec files are unchanged.
- `git diff --check` reports no whitespace errors.
- Optional baseline check: run `go test ./...` to confirm the two known `TestToTsQuery*` failures referenced by the open questions. Do not fix them. Build, vet, and gofmt are not required for a docs-only change, but run them if the user wants a full check.

## Risks and notes

- `go test ./...` is not green on `main`. The TCD must say so and must not claim a passing suite.
- The dedup conflict is a requirement-versus-test decision. The TCD follows `FR-006` and records the conflict. If the user wants the tests to define behavior instead, the expected result for `TC-REG-001` must be flipped.
- The subquery JOIN uses `WHERE` instead of `ON`. The TCD documents the current behavior and flags `ON` as the desired outcome pending a requirement decision.
- The TCD will be long because the library has many clauses. Keep it readable with domain-level preconditions and table-driven groups.
- Do not document MySQL support beyond the reserved constants.
- Do not copy the template's section 11 into the TCD.
- Do not modify `.commandcode/taste/**` or any learned preference file.
