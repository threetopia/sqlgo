# Plan: sqlgo FSD/TSD/DRD specs + README rewrite

## Goal

1. Create an English specification set for `sqlgo` in `docs/specs/`: FSD, TSD, and DRD that share one date and slug.
2. Replace the default GitLab `README.md` with a user-facing guide that explains installation, quick start, core usage, and PostgreSQL-specific features.

No source or test files are changed by this plan. The existing baseline test failures (`TestToTsQuery`, `TestToTsQueryPrepend`) are documented, not fixed.

## Decisions

- **Format:** separate FSD + TSD + DRD files, as requested. Do not combine them.
- **Language:** English for both specs and README.
- **Date/slug:** `2026-11-04` / `sqlgo` for all three files (same date and slug, TYPE differs), per `.agents/skills/writing-specs`.
- **No ERD/DFD:** `sqlgo` owns no schema and the user did not request diagrams; the TSD gets one Mermaid sequence diagram of the query-build flow.
- **DRD adaptation:** `sqlgo` is schema-agnostic and never executes SQL, so the DRD documents the library's data model (parameter registry, accepted Go value types, PostgreSQL extension requirements) and states explicitly that no database schema is owned or migrated.
- **README depth:** quick start + API overview + link to the full spec, not a full API reference.

## Files

Create:

- `docs/specs/2026-11-04_FSD_sqlgo.md`
- `docs/specs/2026-11-04_TSD_sqlgo.md`
- `docs/specs/2026-11-04_DRD_sqlgo.md`

Modify:

- `README.md` (full replacement)

Leave unchanged:

- `docs/specs/1970-01-31_TYPE_spec-template.md` (template stays as-is)
- `.agents/skills/**`, `AGENTS.md`, source and test files

## Source references for factual accuracy

Read these before writing; all behavior below must match the code.

- `sqlgo.go` — `SQLGo` facade, composition order in `BuildSQL()`, dialect constants.
- `sqlgo_<clause>.go` — clause interfaces, constructors, and emitted SQL.
- `sqlgo_parameter.go` — dedup via `hash(value, fmt.Sprintf("%T", value))`, first-appearance order.
- `sqlgo_where.go` — `specialOperator` map, slice handling, `BETWEEN`, `TO TSQUERY`, `similarity`.
- `sqlgo_test.go` — verified SQL strings to reuse in README examples.

## Document contents

### 1. `2026-11-04_FSD_sqlgo.md`

Header: title, document type (FSD), date, status (`Draft`), scope, and links to TSD/DRD (`./2026-11-04_TSD_sqlgo.md`, `./2026-11-04_DRD_sqlgo.md`).

Sections:

1. **Background and purpose** — query builder for Go, currently PostgreSQL-only, prevents manual placeholder bookkeeping.
2. **Goals and success criteria** — chainable construction, deterministic SQL, safe parameter binding, PostgreSQL feature coverage.
3. **Scope and exclusions**
   - In scope: SELECT/FROM/JOIN/WHERE/GROUP BY/ORDER BY/OFFSET-LIMIT, INSERT/UPDATE/DELETE/VALUES, schema prefix, parameter registry, tsvector/tsquery, pg_trgm similarity, pgvector embeddings.
   - Out of scope: MySQL builder implementation, connection management, query execution, migrations, transaction handling, ORM mapping.
4. **Actors** — Go application developers and services using the library.
5. **Functional requirements** with IDs `FR-001` … `FR-019`, each with behavior and acceptance criteria:
   - `FR-001` SELECT projection with optional aliases.
   - `FR-002` FROM with alias and optional column list.
   - `FR-003` JOIN with `ON` conditions and subquery joins.
   - `FR-004` WHERE with `AND`/`OR`, comparison operators, and special operators (`ANY`, `ILIKE ANY`, `LIKE ANY`, `IN`, `NOT IN`, `LIKE`, `NOT LIKE`, `ILIKE`, `NOT ILIKE`).
   - `FR-005` WHERE `BETWEEN` and grouped conditions.
   - `FR-006` Parameter binding and placeholder generation.
   - `FR-007` Schema prefix for insert/update/delete/from/join.
   - `FR-008` GROUP BY.
   - `FR-009` ORDER BY including embedding distance ordering.
   - `FR-010` OFFSET/LIMIT and page-based helper.
   - `FR-011` INSERT with one or more value rows.
   - `FR-012` UPDATE including `to_tsvector`.
   - `FR-013` DELETE with WHERE.
   - `FR-014` standalone VALUES rows.
   - `FR-015` `DISTINCT ON`.
   - `FR-016` full-text search (`ts_rank`, `to_tsquery`, `to_tsvector`).
   - `FR-017` trigram similarity projection and filter.
   - `FR-018` pgvector embedding projection, array projection, and ordering.
   - `FR-019` two equivalent usage styles (prepend/chain vs `Set*`) and fixed clause composition order.
6. **Main flow** — create facade, add clauses, call `BuildSQL()`, pass `GetSQLParameter()` to `database/sql`.
7. **Business rules and validation**
   - Empty clauses emit an empty string and are skipped.
   - Schema prefix is first-write-wins (`SQLSchema` only sets when empty).
   - Identical values share a placeholder; order is first appearance.
   - Slice values for `IN`/`NOT IN` are enumerated and deduplicated; other slice operators use `pq.Array`.
   - `SetSQLPage` derives offset from the current limit, so limit must be set first.
   - `OFFSET` is emitted when offset ≥ 0 and limit > 0.
   - `SetSQLJoinWhere` writes a literal value, not a bound parameter.
8. **Empty, error, and edge states** — no Go error is returned; invalid usage yields incomplete SQL; `nil` values are ignored by the registry and should be avoided; unknown placeholder lookups return `$-1`.
9. **Acceptance criteria** — testable statements tied to the existing test constants.
10. **Traceability** — table mapping `FR-*` to TSD sections and DRD data elements.

### 2. `2026-11-04_TSD_sqlgo.md`

Header with links to FSD/DRD and the same date/slug.

Sections:

1. **Architecture overview** — facade (`sqlGo`), clause modules (`SQLGo<Clause>` / `sqlGo<Clause>`), parameter registry, utilities.
2. **Component responsibilities** — table per file/module.
3. **Public interface contracts** — `SQLGo`, `SQLGOModuleSetter`, `SQLGoBase`, clause interfaces, `SQLGoParameter`, and exported constructors (`NewSQLGo`, `SetSQLSelect`, `SetSQLWhere`, `SetSQLJoinWhere`, `SetSQLInsertToTsVector`, `SetSQLUpdateToTsVector`, `AddSQLGoValue`, `AddSQLGoValueArray`).
4. **Internal data structures** — clause structs and the value wrappers (`sqlGoWhereBetween`, `sqlGoWhereToTsQuery`, `sqlGoSimilarity`, `sqlGoSelectTsRank`, `sqlGoSelectEmbedding`, `sqlGoSelectEmbeddingArray`, `sqlGoSelectSimilarity`, `sqlGoUpdateToTsVector`, `sqlGoInsertToTsVector`).
5. **Clause composition order** — table matching `sqlGo.BuildSQL()`: insert → update → delete → select → from → join → where → values → group by → order → offset/limit.
6. **Parameter binding design** — SHA-1 over `json.Marshal(value)` plus `%T`; dedup map and ordered slice; `$N` generation; first-appearance order contract.
7. **PostgreSQL feature design** — `to_tsvector`/`to_tsquery`/`ts_rank`, `similarity`, pgvector operators (`<->`, `<=>`, `<#>`), `pq.Array`, supported slice types.
8. **Technical flow** — Mermaid `sequenceDiagram`: application → facade → clause `BuildSQL()` → parameter registry → SQL string + args → `database/sql` + `lib/pq` → PostgreSQL. Keep it consistent with the prose.
9. **Concurrency and lifecycle** — builder holds mutable state; not safe for concurrent use; `BuildSQL()` should be treated as a terminal call.
10. **Error handling** — no errors returned; invalid usage surfaces as malformed SQL; callers validate identifiers.
11. **Nonfunctional requirements** — Go 1.18 compatibility, single runtime dependency `github.com/lib/pq`, allocation-light string assembly, deterministic output.
12. **Verification strategy** — `go build ./...`, `go test ./...`, `go vet ./...`, `gofmt -l .`; reference existing tests; record the known failing tests and why.
13. **Traceability** — `FR-*` to components/tests.

### 3. `2026-11-04_DRD_sqlgo.md`

Header with links to FSD/TSD and the same date/slug.

Sections:

1. **Scope and data ownership** — the library owns no database schema; it only formats SQL and parameter values. No migration, backfill, or persistence.
2. **Data model** — logical entities: parameter registry entry (value, Go type, hash, ordinal), clause value wrappers, generated statement.
3. **Data dictionary** — accepted Go types and behavior:
   - scalars: `string`, numeric types, `bool`, `time.Time`
   - slices: `[]string`, `[]int`, `[]int32`, `[]int64`, `[]float32`, `[]float64`, `[]bool`, `[]time.Time`
   - `pq.Array` for non-`IN` slice operators
   - `nil` handling
   - wrapper types and their SQL meaning.
4. **Integrity rules** — placeholder ordinals are contiguous per build; identical `(value, type)` pairs share an ordinal; array element order is preserved after dedup; no uniqueness constraint beyond the hash map.
5. **Access patterns** — write path registers a parameter; read path resolves the sign; callers read the ordered slice.
6. **PostgreSQL data requirements** — built-in `tsvector`/`tsquery`; `pg_trgm` for `similarity`; `pgvector` for embeddings; caller-owned table/column names.
7. **Impact on existing data** — none; no schema changes.
8. **Open questions** — confirm whether identical-value deduplication is intentional (the two failing tests expect separate placeholders); decide if `TO TSQUERY` should remain in `specialOperator`; document nil-value behavior as unsupported or add explicit handling.
9. **Traceability** — `FR-006`, `FR-016`, `FR-017`, `FR-018` mapped to data elements.

### 4. `README.md`

Full replacement in English:

1. **Title and one-line description** — chainable PostgreSQL query builder for Go.
2. **Features**
3. **Requirements** — Go 1.18+, PostgreSQL; optional `pg_trgm` and `pgvector` for their features.
4. **Installation** — `go get gitlab.com/threetopia/sqlgo/v2`.
5. **Quick start** — SELECT with schema, join, WHERE, group, order, pagination; print `BuildSQL()` and `GetSQLParameter()`; run with `db.Query`.
6. **Building statements** — compact examples for INSERT, UPDATE, DELETE, each showing the produced SQL.
7. **Parameter binding** — explain `BuildSQL()` + `GetSQLParameter()`, first-appearance order, and identical-value deduplication.
8. **Two usage styles** — chain/prepend vs `Set*`; both produce identical SQL.
9. **PostgreSQL features** — full-text search, trigram similarity, pgvector embedding select/order examples.
10. **Clause composition order** — table.
11. **API overview** — grouped method/constructor table.
12. **Limitations** — PostgreSQL only, `MySQLDialect` reserved but unimplemented, no execution/connection handling, no errors returned.
13. **Documentation** — links to `docs/specs/2026-11-04_FSD_sqlgo.md`, TSD, and DRD.
14. **License** — GNU GPL v3, link to `LICENSE`.

README examples must use only behavior verified by `sqlgo_test.go` or direct inspection. Prefer examples whose expected SQL is already covered by tests (select, insert, update, delete, tsquery, similarity, embedding).

## Execution steps

1. Write `docs/specs/2026-11-04_FSD_sqlgo.md` with the sections and `FR-001`…`FR-019` above.
2. Write `docs/specs/2026-11-04_TSD_sqlgo.md`, referencing the FSD requirement IDs and including the Mermaid sequence diagram.
3. Write `docs/specs/2026-11-04_DRD_sqlgo.md`, referencing the FSD/TSD and clearly stating no schema is owned.
4. Replace `README.md` with the outline above.
5. Cross-check every documented signature, operator, and SQL shape against the source files.
6. Verify relative links resolve (`docs/specs` siblings, `../../../LICENSE` from docs, `LICENSE` from README).
7. Report assumptions/open questions (notably the dedup-vs-test conflict) instead of silently deciding them.

## Verification

- Filenames match `{yyyy-mm-dd}_{TYPE}_{slug}.md` and all three share `2026-11-04` + `sqlgo`.
- FSD and TSD are separate files; DRD exists and states the library owns no database schema.
- All `FR-*` referenced by TSD/DRD exist in the FSD.
- Mermaid block is valid `sequenceDiagram` syntax.
- Every README code sample compiles conceptually and its documented SQL matches the implementation; for the main samples, add a temporary `readme_examples_test.go` in the repo, run `go test -run TestReadme`, compare outputs, then delete the file.
- `go build ./...`, `go vet ./...`, `gofmt -l .` remain clean.
- `go test ./...` still shows only the two known `TestToTsQuery*` failures; do not fix them in this task.
- `git diff --check` reports no whitespace errors.

## Risks and notes

- `go test ./...` is not green on `main`; the specs and README must not claim otherwise.
- The DRD is an adaptation, not an application schema document; say so in its scope section to avoid implying entities that do not exist.
- Avoid documenting MySQL support beyond the reserved constants.
- Do not copy the 18-section service template wholesale; most service sections do not apply to a library.
- Keep the existing template file untouched so future specs can still use it.
