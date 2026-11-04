# sqlgo

A chainable SQL query builder for Go. `sqlgo` builds PostgreSQL statements and keeps the SQL text and the bound argument slice in sync, so you do not manage `$1`, `$2`, and values by hand.

## Features

- One builder for SELECT, INSERT, UPDATE, DELETE, and VALUES.
- Parameterized values with automatic `$N` placeholders and deduplication.
- Schema prefixes, aliases, joins, grouped conditions, pagination, and ordering.
- PostgreSQL full-text search with `tsvector` and `tsquery`.
- PostgreSQL trigram similarity with `pg_trgm`.
- PostgreSQL vector distance with `pgvector`.
- Two equivalent usage styles: chain methods and `Set*` methods.

## Requirements

- Go 1.18 or newer.
- PostgreSQL.
- The `pg_trgm` extension when you use similarity.
- The `pgvector` extension when you use embedding projections or ordering.

## Installation

```sh
go get gitlab.com/threetopia/sqlgo/v2
```

## Quick Start

```go
package main

import (
    "database/sql"
    "fmt"

    "gitlab.com/threetopia/sqlgo/v2"
    _ "github.com/lib/pq"
)

func main() {
    query := sqlgo.NewSQLGo().
        SQLSchema("public").
        SQLSelect(
            sqlgo.SetSQLSelect("u.id", "id"),
            sqlgo.SetSQLSelect("u.name", "name"),
        ).
        SQLFrom("users", "u").
        SQLJoin(
            sqlgo.SetSQLJoin("INNER", "orders", "o",
                sqlgo.SetSQLJoinWhere("AND", "o.user_id", "=", "u.id"),
            ),
        ).
        SQLWhere(
            sqlgo.SetSQLWhere("AND", "u.active", "=", true),
            sqlgo.SetSQLWhere("AND", "o.total", ">=", 100),
        ).
        SQLGroupBy("u.id", "u.name").
        SQLOrder(sqlgo.SetSQLOrder("u.name", "ASC")).
        SQLPageLimit(1, 20)

    sqlText := query.BuildSQL()
    args := query.GetSQLGoParameter().GetSQLParameter()

    fmt.Println(sqlText)
    fmt.Println(args)

    db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/app?sslmode=disable")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    rows, err := db.Query(sqlText, args...)
    if err != nil {
        panic(err)
    }
    defer rows.Close()
}
```

The builder produces:

```sql
SELECT u.id AS id, u.name AS name FROM public.users AS u INNER JOIN public.orders AS o ON o.user_id=u.id WHERE u.active=$1 AND o.total>=$2 GROUP BY u.id, u.name ORDER BY u.name ASC OFFSET 0 LIMIT 20
```

## Building Statements

### INSERT

```go
query := sqlgo.NewSQLGo().
    SQLInsert("users",
        sqlgo.SetSQLInsertColumn("name", "active"),
        sqlgo.SetSQLInsertValue("Alice", true),
        sqlgo.SetSQLInsertValue("Bob", false),
    )
```

```sql
INSERT INTO users (name, active) VALUES ($1, $2), ($3, $4)
```

Arguments: `["Alice", true, "Bob", false]`.

### UPDATE

```go
query := sqlgo.NewSQLGo().
    SQLUpdate("users", sqlgo.SetSQLUpdateValue("name", "Alice")).
    SQLWhere(sqlgo.SetSQLWhere("AND", "id", "=", 1))
```

```sql
UPDATE users SET name=$1 WHERE id=$2
```

### DELETE

```go
query := sqlgo.NewSQLGo().
    SQLDelete("users", sqlgo.SetSQLWhere("AND", "id", "=", 1))
```

```sql
DELETE FROM users WHERE id=$1
```

## Parameter Binding

`BuildSQL()` returns the SQL string. `GetSQLGoParameter().GetSQLParameter()` returns the argument slice in first-appearance order, which matches `$1`, `$2`, and so on.

Identical values of the same Go type share one placeholder:

```go
query := sqlgo.NewSQLGo().
    SQLSelect(sqlgo.SetSQLSelect("id", "id")).
    SQLFrom("users", "u").
    SQLWhere(sqlgo.SetSQLWhere("AND", "status", "=", "active")).
    SQLWhere(sqlgo.SetSQLWhere("AND", "role", "=", "active"))
```

```sql
SELECT id AS id FROM users AS u WHERE status=$1 AND role=$1
```

Arguments: `["active"]`. The value is sent once because both conditions use it.

Values of different types do not share a placeholder. For example, `int(1)` and `float64(1)` are registered separately.

## Two Usage Styles

Prepend style passes clause values directly:

```go
a := sqlgo.NewSQLGo().
    SQLSelect(sqlgo.SetSQLSelect("id", "id")).
    SQLFrom("users", "u").
    SQLWhere(sqlgo.SetSQLWhere("AND", "id", "=", 1))
```

Setter style configures one item at a time:

```go
b := sqlgo.NewSQLGo().
    SetSQLSelect("id", "id").
    SetSQLFrom("users", "u").
    SetSQLWhere("AND", "id", "=", 1)
```

Both produce `SELECT id AS id FROM users AS u WHERE id=$1`.

## PostgreSQL Features

### Full-Text Search

```go
query := sqlgo.NewSQLGo().
    SQLSelect(sqlgo.SetSQLSelectTsRank("search_vector", "english", "postgres & sql", "rank")).
    SQLFrom("articles", "a").
    SQLWhere(sqlgo.SetSQLWhereToTsQuery("AND", "search_vector", "english", "postgres & sql"))
```

```sql
SELECT ts_rank(search_vector, to_tsquery('english', $1)) AS rank FROM articles AS a WHERE search_vector @@ to_tsquery('english', $1)
```

Arguments: `["postgres & sql"]`. The repeated value is deduplicated.

For inserts and updates, `SetSQLInsertToTsVector` and `SetSQLUpdateToTsVector` wrap a value in `to_tsvector`.

### Trigram Similarity

```go
query := sqlgo.NewSQLGo().
    SQLSelect(
        sqlgo.SetSQLSelect("id", "id"),
        sqlgo.SetSQLSelectSimilarity("full_name", "John Doe", "similarity"),
    ).
    SQLFrom("users", "u").
    SQLWhere(sqlgo.SetSQLWhereSimilarity("AND", "full_name", "John Doe", ">=", 0.75))
```

```sql
SELECT id AS id, similarity(full_name, $1) AS similarity FROM users AS u WHERE similarity(full_name, $1) >= $2
```

Arguments: `["John Doe", 0.75]`.

### Vector Distance

```go
query := sqlgo.NewSQLGo().
    SQLSelect(sqlgo.SetSQLSelectEmbedding("", "embedding", "<=>", "[0.1,0.2,0.3]", "distance")).
    SQLFrom("documents", "d").
    SQLOrder(sqlgo.SetSQLOrderEmbedding("embedding", "<=>", "[0.1,0.2,0.3]")).
    SQLOffsetLimit(0, 10)
```

```sql
SELECT (embedding <=> $1) AS distance FROM documents AS d ORDER BY embedding <=> $1 OFFSET 0 LIMIT 10
```

Arguments: `["[0.1,0.2,0.3]"]`. Pass the vector in the text form your pgvector setup expects. The operator is caller supplied, for example `<->`, `<=>`, or `<#>`.

## Clause Composition Order

The builder composes clauses in a fixed order, regardless of the order you call the methods:

| Order | Clause |
|---|---|
| 1 | INSERT |
| 2 | UPDATE |
| 3 | DELETE |
| 4 | SELECT |
| 5 | FROM |
| 6 | JOIN |
| 7 | WHERE |
| 8 | VALUES |
| 9 | GROUP BY |
| 10 | ORDER BY |
| 11 | OFFSET/LIMIT |

## API Overview

| Area | Methods and constructors |
|---|---|
| Facade | `NewSQLGo`, `BuildSQL`, `GetSQLGoParameter` |
| Schema | `SQLSchema`, `SetSQLSchema` |
| Select | `SQLSelect`, `SetSQLSelect`, `SetSQLSelectDistinct`, `SetSQLSelectTsRank`, `SetSQLSelectSimilarity`, `SetSQLSelectEmbedding`, `SetSQLSelectEmbeddingArray` |
| From | `SQLFrom`, `SetSQLFrom`, `SetSQLFromColumn` |
| Join | `SQLJoin`, `SetSQLJoin`, `SetSQLJoinWhere` |
| Where | `SQLWhere`, `SetSQLWhere`, `SetSQLWhereBetween`, `SQLWhereGroup`, `SetSQLWhereGroup`, `SetSQLWhereToTsQuery`, `SetSQLWhereSimilarity` |
| Group by | `SQLGroupBy`, `SetSQLGroupBy` |
| Order | `SQLOrder`, `SetSQLOrder`, `SetSQLOrderEmbedding` |
| Offset and limit | `SQLOffsetLimit`, `SetSQLOffset`, `SetSQLLimit`, `SQLPageLimit`, `SetSQLPage` |
| Insert | `SQLInsert`, `SetSQLInsert`, `SetSQLInsertColumn`, `SetSQLInsertValue`, `SetSQLInsertToTsVector` |
| Update | `SQLUpdate`, `SetSQLUpdate`, `SetSQLUpdateValue`, `SetSQLUpdateValueSlice`, `SetSQLUpdateToTsVector` |
| Delete | `SQLDelete`, `SetSQLDelete` |
| Values | `SQLValues`, `SetSQLValues`, `SetSQLValuesValue` |
| Parameters | `SetSQLParameter`, `GetSQLParameterCount`, `GetSQLParameterSign`, `GetSQLParameter` |
| Helpers | `AddSQLGoValue`, `AddSQLGoValueArray` |

## Limitations

- PostgreSQL only. `MySQLDialect` is defined but no MySQL builder branch exists.
- The library does not open connections, run queries, manage transactions, or run migrations.
- Builder methods do not return errors. Invalid identifiers or unsupported operators produce malformed SQL or a driver error.
- A builder holds mutable state. Create one builder per query and do not share it across goroutines.
- Identifiers such as table and column names are written as provided and are not quoted or escaped.
- `SetSQLPage` reads the current limit, so set the limit before the page.

## Documentation

- [Functional Specification](./docs/specs/2026-11-04_FSD_sqlgo.md)
- [Technical Specification](./docs/specs/2026-11-04_TSD_sqlgo.md)
- [Data Requirements](./docs/specs/2026-11-04_DRD_sqlgo.md)

## License

GNU General Public License v3. See [LICENSE](./LICENSE).
