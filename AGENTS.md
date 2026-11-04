# AGENTS.md

Panduan kerja agen di repository `sqlgo`.

## Ringkasan Project

`sqlgo` adalah library query builder SQL untuk Go (module `gitlab.com/threetopia/sqlgo/v2`, package `sqlgo`). Semua kode ada di root package.

Scope implementasi saat ini: **PostgreSQL**. Placeholder memakai format `$N`, parameter dikirim lewat `github.com/lib/pq`, dan tersedia builder khusus PostgreSQL: `tsvector`/`tsquery`, `similarity` (pg_trgm), serta embedding (pgvector `<=>`). Konstanta `MySQLDialect`/`PostgreSQLDialect` ada di `sqlgo.go`, tetapi belum ada cabang builder khusus MySQL.

API berbentuk chain:

```go
sql := sqlgo.NewSQLGo().
    SQLSchema("schema").
    SQLSelect(sqlgo.SetSQLSelect("t.id", "id")).
    SQLFrom("table", "t").
    SQLWhere(sqlgo.SetSQLWhere("AND", "t.id", "=", 1))

query := sql.BuildSQL()                          // string SQL
args := sql.GetSQLGoParameter().GetSQLParameter() // []interface{}
```

## Perintah

```sh
go build ./...
go test ./...
go vet ./...
gofmt -l .
```

`go.mod` menargetkan `go 1.18`; jangan memakai fitur bahasa di atas itu tanpa menaikkan versi module. Tidak ada konfigurasi CI, linter, atau Makefile di repo — validasi minimal adalah build, test, vet, dan gofmt.

## Struktur File

| File | Isi |
|---|---|
| `sqlgo.go` | Facade `SQLGo`, `NewSQLGo()`, urutan komposisi `BuildSQL()`, konstanta dialect |
| `sqlgo_<clause>.go` | Satu file per clause: `insert`, `update`, `delete`, `select`, `from`, `join`, `where`, `values`, `group_by`, `order`, `offset_limit`, `schema` |
| `sqlgo_parameter.go` | Registry parameter dan generator placeholder |
| `util.go` | Helper `combineSQL` dan `hash` |
| `sqlgo_test.go` | Tes pembentukan SQL dan parameter |

## Arsitektur

- Setiap clause punya interface publik `SQLGo<Clause>`, implementasi privat `sqlGo<Clause>`, dan mengimplementasikan `SQLGoBase` (`GetSQLGoParameter`, `BuildSQL`) plus setter `SetSQLGoParameter`.
- Ada dua gaya pemakaian: method prepend/chain (`SQLSelect(...)`, `SQLWhere(...)`) dan method `SetSQL*` satu per satu (`SetSQLSelect`, `SetSQLWhere`, ...). Keduanya harus menghasilkan SQL identik.
- Paket juga mengekspos fungsi konstruktor value, mis. `SetSQLSelect`, `SetSQLWhere`, `SetSQLOrder`, `AddSQLGoValue`, `AddSQLGoValueArray`.
- Facade `sqlGo.BuildSQL()` merangkai clause dengan urutan tetap: insert → update → delete → select → from → join → where → values → group by → order → offset/limit. Clause kosong mengembalikan string kosong dan dilewati oleh `combineSQL`.
- Parameter didaftarkan lewat `SQLGoParameter.SetSQLParameter`; nilai yang identik (hash SHA-1 dari JSON + tipe) memakai nomor placeholder yang sama. Urutan `GetSQLParameter()` adalah urutan kemunculan pertama.
- Prefix schema (`SetSQLSchema`/`SQLSchema`) hanya dipakai oleh insert, update, delete, from, dan join.

## Konvensi Kode

- Interface dan tipe publik memakai prefix `SQLGo`; tipe privat memakai prefix `sqlGo`.
- Method builder mengembalikan interface clause/facade agar bisa di-chain.
- Jangan menyisipkan nilai user langsung ke string SQL. Semua nilai dinamis harus lewat `SetSQLParameter` lalu `GetSQLParameterSign`.
- Menambah clause baru: buat file `sqlgo_<clause>.go` (interface, struct, constructor, `BuildSQL`), daftarkan field dan `SetSQLGo<Clause>` pada `SQLGOModuleSetter`/`sqlGo` di `sqlgo.go`, lalu panggil modul tersebut di `sqlGo.BuildSQL()` sesuai urutan SQL.
- Ikuti pola `fmt.Sprintf` yang sudah ada agar spasi antar clause konsisten; tes membandingkan string secara persis.

## Testing

- Tes membandingkan string SQL persis (termasuk spasi) terhadap konstanta query, lalu memeriksa `GetSQLGoParameter()`.
- Setiap perubahan builder wajib menjalankan `go test ./...` dan menyesuaikan konstanta query yang terdampak.
- **Status baseline:** `go test ./...` saat ini gagal pada `TestToTsQuery` dan `TestToTsQueryPrepend`. Keduanya mengharapkan `$2`, sedangkan dedup parameter mengembalikan `$1` untuk nilai yang sama. Jangan menganggap kegagalan ini disebabkan perubahanmu; perbaiki tes/ekspektasi hanya jika task memang menyentuhnya.

## Dokumentasi

- `docs/brief/`, `docs/plans/`, `docs/specs/`, `docs/test-cases/` berisi template. Duplikat template dan beri nama `YYYY-MM-DD_<TYPE>_<slug>.md`.
- `README.md` masih template default GitLab dan belum mendeskripsikan library; jangan jadikan acuan perilaku.
- `CLAUDE.md` hanya berisi referensi ke file ini.

## Batasan

- Library publik versi v2: hindari breaking change pada signature yang sudah ada. Tambahkan method/tipe baru alih-alih mengubah yang lama.
- Satu-satunya dependency runtime adalah `github.com/lib/pq`.
- Jangan mengubah urutan pendaftaran parameter placeholder tanpa menyesuaikan tes dan pemanggil, karena urutan `GetSQLParameter()` adalah kontrak pemakaian.
