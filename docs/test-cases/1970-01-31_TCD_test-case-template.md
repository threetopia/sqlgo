# TCD: {Nama Fitur atau Aplikasi}

> Template. Ganti setiap `{placeholder}` dan baris contoh dengan data nyata, lalu perbarui status sebelum dokumen disetujui. Contoh di bawah memakai fitur fiktif "Manajemen Pelanggan" agar bentuk pengisiannya terlihat. Aturan lengkap ada di [skill writing-test-case](../../.agents/skills/writing-test-case/SKILL.md).

## 1. Informasi Dokumen

| Field | Isi |
|---|---|
| Judul | TCD: Manajemen Pelanggan |
| Tipe dokumen | TCD (Test Case Document) |
| Tanggal dibuat | 1970-01-31 |
| Terakhir diperbarui | 1970-01-31 |
| Status | Draft |
| Fitur / aplikasi / modul | `api_customer`, modul manajemen pelanggan |
| Ruang lingkup | API dan alur UI pembuatan, pencarian, dan pembaruan data pelanggan |
| Eksklusi | Migrasi data historis, laporan analitik, integrasi pihak ketiga di luar `api_customer` |
| FSD | [`{yyyy-mm-dd}_FSD_{slug}.md`](../specs/{yyyy-mm-dd}_FSD_{slug}.md) |
| TSD | [`{yyyy-mm-dd}_TSD_{slug}.md`](../specs/{yyyy-mm-dd}_TSD_{slug}.md) |
| DRD | [`{yyyy-mm-dd}_DRD_{slug}.md`](../specs/{yyyy-mm-dd}_DRD_{slug}.md) |
| Brief / Plan | [`{yyyy-mm-dd}_{brief}.md`](../brief/{yyyy-mm-dd}_{brief}.md) |
| Pemilik | {nama atau role penanggung jawab} |

Aturan penamaan file: `{yyyy-mm-dd}_TCD_{slug}.md`, dengan tanggal dan slug yang sama seperti dokumen spesifikasi pasangannya. TCD yang dibuat belakangan tetap memakai tanggal dan slug set spesifikasi tersebut.

## 2. Tujuan dan Ruang Lingkup Pengujian

### 2.1 Perilaku yang Diverifikasi

- {perilaku utama 1, misal: pelanggan valid tersimpan dan mengembalikan id unik}
- {perilaku utama 2, misal: pencarian mengembalikan hasil sesuai filter dan pagination}
- {perilaku utama 3, misal: input tidak valid ditolak tanpa mengubah data}

### 2.2 Requirement yang Dicakup

- `FR-001` sampai `FR-004`: {ringkasan}
- `FR-010`: {ringkasan}
- {requirement lain sesuai dokumen FSD/TSD/DRD}

### 2.3 Eksklusi

- {yang sengaja tidak diuji, beserta alasannya}

### 2.4 Dependensi, Prasyarat, dan Asumsi Environment

- Database uji terpisah dengan migrasi terbaru sudah dijalankan.
- Seed data: {jumlah dan isi data minimum}.
- Variabel environment: {daftar env yang relevan, tanpa nilai rahasia}.
- Versi service, model, atau prompt yang memengaruhi hasil: {isi atau `tidak relevan`}.

## 3. Ketertelusuran Requirement

| Requirement ID | Sumber | Ringkasan | Test Case | Cakupan |
|---|---|---|---|---|
| FR-001 | FSD | Pelanggan dibuat dengan data valid | TC-CUSTOMER-001 | Positif |
| FR-002 | FSD | Data wajib tidak boleh kosong | TC-CUSTOMER-002 | Negatif |
| FR-003 | FSD | Email unik per tenant | TC-CUSTOMER-003 | Boundary |
| FR-010 | FSD | Agent tidak mengarang data pelanggan | TC-AGENT-001 | AI Agent |

Requirement yang belum punya test case dicatat eksplisit di bagian [Asumsi, Pertanyaan Terbuka, dan Gap Cakupan](#9-asumsi-pertanyaan-terbuka-dan-gap-cakupan). Jangan membuat requirement ID baru di TCD.

## 4. Test Case

### 4.1 Field Wajib

| Field | Keterangan |
|---|---|
| Test Case ID | Identifier unik dan stabil, format `TC-{DOMAIN}-{NNN}`. |
| Requirement ID | Requirement FSD, acceptance criterion, atau defect ID yang diuji. |
| Title | Ringkasan singkat perilaku yang diuji. |
| Category | Positive, Negative, Boundary, Error Handling, Security, Regression, Performance, Concurrency, AI Agent. |
| Priority | Critical, High, Medium, Low. |
| Preconditions | Kondisi awal, permission, dependensi, dan setup yang dibutuhkan. |
| Test Data | Input, fixture, dan kondisi data yang relevan. |
| Test Steps | Langkah terurut yang dapat dieksekusi ulang. |
| Expected Result | Hasil yang spesifik, teramati, dan dapat diverifikasi. |
| Actual Result | Hasil yang dicatat setelah eksekusi. |
| Status | Not Run, Pass, Fail, atau Blocked. Default: `Not Run`. |
| Evidence / Defect ID | Log, artefak, tangkapan layar, atau referensi defect. |

Field tambahan bila relevan: environment, expected HTTP status, assertion database, expected tool calls, versi model, versi prompt, dan kriteria evaluasi.

### 4.2 Ringkasan

| Test Case ID | Requirement | Judul | Kategori | Prioritas | Status |
|---|---|---|---|---|---|
| TC-CUSTOMER-001 | FR-001 | Membuat pelanggan dengan data valid | Positive | Critical | Not Run |
| TC-CUSTOMER-002 | FR-002 | Menolak pembuatan pelanggan tanpa nama | Negative | High | Not Run |
| TC-CUSTOMER-003 | FR-003 | Menolak email duplikat dalam tenant yang sama | Boundary | High | Not Run |
| TC-AGENT-001 | FR-010 | Agent tidak mengarang data pelanggan | AI Agent | High | Not Run |

### 4.3 Detail Test Case

#### TC-CUSTOMER-001: Membuat pelanggan dengan data valid

| Field | Isi |
|---|---|
| Requirement ID | FR-001, FR-004 |
| Kategori | Positive |
| Prioritas | Critical |
| Prasyarat | Service berjalan, database uji kosong untuk tenant `tenant-a`, aktor punya permission `customer:create`. |
| Test Data | `{"name": "Budi Santoso", "email": "budi@example.com", "tenant_id": "tenant-a"}` |
| Expected Result | HTTP 201, body berisi `id` unik dan `name` yang sama, satu baris baru di `customers` dengan `email` ternormalisasi, tidak ada error di log. |
| Actual Result | (diisi saat eksekusi) |
| Status | Not Run |
| Evidence / Defect ID | - |

Langkah:
1. Kirim `POST /customers` dengan Test Data di atas.
2. Periksa status code dan body respons.
3. Query tabel `customers` untuk `email` tersebut.
4. Periksa log service untuk error atau panic.

#### TC-CUSTOMER-002: Menolak pembuatan pelanggan tanpa nama

| Field | Isi |
|---|---|
| Requirement ID | FR-002 |
| Kategori | Negative |
| Prioritas | High |
| Prasyarat | Sama seperti TC-CUSTOMER-001. |
| Test Data | `{"name": "", "email": "tanpa-nama@example.com", "tenant_id": "tenant-a"}` |
| Expected Result | HTTP 422 dengan pesan validasi yang menyebut field `name`, tidak ada baris baru di `customers`, dan tidak ada perubahan pada data lain. |
| Actual Result | (diisi saat eksekusi) |
| Status | Not Run |
| Evidence / Defect ID | - |

Langkah:
1. Kirim `POST /customers` dengan Test Data di atas.
2. Periksa status code dan pesan error.
3. Pastikan jumlah baris `customers` tidak bertambah.

#### TC-CUSTOMER-003: Menolak email duplikat dalam tenant yang sama

| Field | Isi |
|---|---|
| Requirement ID | FR-003 |
| Kategori | Boundary |
| Prioritas | High |
| Prasyarat | Sudah ada pelanggan `budi@example.com` pada `tenant-a`. |
| Test Data | `{"name": "Budi Lain", "email": "BUDI@example.com", "tenant_id": "tenant-a"}` |
| Expected Result | HTTP 409 dengan kode error `EMAIL_ALREADY_EXISTS`, tidak ada baris baru, dan data pelanggan lama tidak berubah. |
| Actual Result | (diisi saat eksekusi) |
| Status | Not Run |
| Evidence / Defect ID | - |

Langkah:
1. Kirim `POST /customers` dengan Test Data di atas.
2. Periksa status code dan kode error.
3. Pastikan hanya satu baris dengan email tersebut pada `tenant-a`.

#### TC-AGENT-001: Agent tidak mengarang data pelanggan

| Field | Isi |
|---|---|
| Requirement ID | FR-010 |
| Kategori | AI Agent, Negative |
| Prioritas | High |
| Prasyarat | Agent aktif dengan tool `search_customer`, dataset uji tidak memuat `INVALID-001`. |
| Test Data | Prompt: `"Cari pelanggan dengan id INVALID-001."` |
| Expected Result | Agent memanggil `search_customer` dengan `{"customer_id": "INVALID-001"}` tepat satu kali, menyatakan pelanggan tidak ditemukan, tidak menyebut nama atau email yang tidak ada di dataset, dan tidak memanggil tool lain. |
| Actual Result | (diisi saat eksekusi) |
| Status | Not Run |
| Evidence / Defect ID | - |

Langkah:
1. Kirim prompt di atas ke agent.
2. Rekam daftar tool call beserta argumennya.
3. Bandingkan respons dengan isi dataset.
4. Catat versi model, versi prompt, dan konfigurasi tool.

## 5. Matriks Cakupan

Isi kolom dengan jumlah test case, atau `n/a` bila kategori tidak relevan beserta alasannya.

| Kategori | Wajib dipertimbangkan | Jumlah | Catatan |
|---|---|---|---|
| Positive | Ya | 1 | Alur sukses utama |
| Negative | Ya | 2 | Validasi dan duplikasi |
| Boundary dan validation | Bila ada batas nilai | 1 | Panjang, rentang, nullability |
| Error handling dan resilience | Bila ada dependensi eksternal | 0 | {alasan} |
| Security | Bila ada auth, izin, atau data sensitif | 0 | {alasan} |
| Regression | Untuk bug fix atau perubahan perilaku | 0 | {alasan} |
| Performance dan concurrency | Bila ada target atau risiko | 0 | {alasan} |
| AI Agent | Bila fitur memakai agent atau LLM | 1 | Grounding dan tool call |

Jangan membuat semua kategori secara mekanis. Sertakan kategori sesuai requirement, risiko, dan mode kegagalan yang nyata.

## 6. Pengujian AI Agent

Isi bagian ini hanya bila fitur melibatkan agent, LLM, atau tool calling. Nilai setiap kategori sebagai relevan, tidak relevan, atau belum diputuskan.

1. Interpretasi intent dan kepatuhan instruksi.
2. Pemilihan tool dan argumen tool call.
3. Eksekusi tool, penanganan error, dan pemulihan.
4. Kebenaran respons dan grounding pada bukti yang tersedia.
5. Informasi yang hilang dan pencegahan halusinasi.
6. Kepatuhan format dan schema output.
7. Penyelesaian alur multi langkah.
8. Penanganan prompt injection dan input tidak tepercaya.
9. Autentikasi, otorisasi, dan isolasi data.
10. Batas token, timeout, retry, dan rate limit.
11. Regresi akibat perubahan model, prompt, tool, atau retrieval.

Aturan expected result untuk AI agent:

- Assert nama tool dan argumen secara eksak hanya bila requirement menetapkan demikian.
- Verifikasi tool wajib benar-benar dipanggil dan tool terlarang tidak dipanggil.
- Verifikasi fakta dalam respons didukung data yang tersedia.
- Verifikasi output terstruktur terhadap schema.
- Tetapkan kriteria relevansi, kelengkapan, grounding, dan refusal sebelum evaluasi.
- Jangan menuntut kesamaan teks natural language kecuali requirement memang memintanya.
- Catat versi model, prompt, dataset, dan konfigurasi yang memengaruhi reproduksibilitas.

Dataset uji sebaiknya mencakup permintaan normal, ambigu, informasi hilang atau bertentangan, argumen tool tidak valid, tool error atau dependensi mati, input prompt injection, dan permintaan di luar wewenang agent.

## 7. Skenario BDD

Gunakan Given-When-Then bila membantu komunikasi lintas tim. Gherkin bersifat opsional dan tidak boleh menambah framework baru tanpa alasan yang jelas.

```gherkin
Feature: Manajemen pelanggan

  Scenario: Pelanggan valid berhasil dibuat
    Given tenant "tenant-a" belum memiliki pelanggan "budi@example.com"
    When pengguna mengirim data pelanggan yang valid
    Then sistem mengembalikan HTTP 201
    And sistem menyimpan satu pelanggan baru

  Scenario: Email duplikat ditolak
    Given tenant "tenant-a" sudah memiliki pelanggan "budi@example.com"
    When pengguna mengirim pelanggan baru dengan email "BUDI@example.com"
    Then sistem mengembalikan HTTP 409
    And jumlah pelanggan pada tenant tersebut tidak berubah
```

## 8. Strategi Otomasi per Layer

| Layer | Tujuan | Pendekatan |
|---|---|---|
| Unit | Verifikasi aturan bisnis dan fungsi terisolasi | Framework test bawaan bahasa |
| Integration | Verifikasi database, API, adapter, dan tool | Framework test dengan fixture dan dependensi terkontrol |
| Acceptance / BDD | Verifikasi perilaku yang teramati dari luar | Given-When-Then, Godog atau Cucumber bila diperlukan |
| LLM evaluation | Evaluasi kualitas respons dan perilaku agent | Harness evaluasi yang repeatable dengan assertion dan kriteria eksplisit |
| Security | Verifikasi batas akses dan input adversarial | Test otomatis dan manual yang terarah |

Untuk proyek Go: pakai `testing`, table-driven test, dan `httptest` bila sesuai. TCD menetapkan perilaku yang diharapkan; tidak semua kasus wajib punya test otomatis. Catat kasus yang diverifikasi manual beserta alasannya.

## 9. Asumsi, Pertanyaan Terbuka, dan Gap Cakupan

| Jenis | Isi | Dampak |
|---|---|---|
| Asumsi | {asumsi yang dipakai saat menulis test case} | {dampak bila salah} |
| Pertanyaan terbuka | {keputusan yang belum ada di spesifikasi} | {test case yang tertahan} |
| Gap cakupan | {requirement tanpa test case} | {risiko} |

Jangan mengarang business rule, schema, atau acceptance criteria. Tandai ketidakpastian di sini.

## 10. Definisi Selesai dan Checklist

TCD siap direview bila:

- [ ] Ruang lingkup dan dokumen terkait sudah diidentifikasi.
- [ ] Kasus positif dan negatif yang relevan sudah didefinisikan.
- [ ] Skenario boundary, error, dan security sudah dipertimbangkan.
- [ ] Expected result spesifik dan dapat diverifikasi.
- [ ] Ketertelusuran requirement tersedia.
- [ ] Perilaku khusus AI agent tercakup bila relevan.
- [ ] Kebutuhan verifikasi otomatis dan manual sudah dibedakan.
- [ ] Asumsi, pertanyaan terbuka, dan gap cakupan tercatat.
- [ ] Semua status default `Not Run`, dan tidak ada klaim `Pass` tanpa eksekusi yang terverifikasi.

## 11. Panduan Penggunaan Template

Untuk manusia:

1. Ganti seluruh `{placeholder}` dan contoh "Manajemen Pelanggan" dengan fitur nyata.
2. Ambil requirement dari `docs/specs/` dan pertahankan ID aslinya.
3. Simpan file sebagai `docs/test-cases/{yyyy-mm-dd}_TCD_{slug}.md` dengan tanggal dan slug yang sama seperti set spesifikasi.
4. Perbarui `Actual Result`, `Status`, dan `Evidence` setelah eksekusi.

Untuk AI agent:

1. Baca FSD, TSD, DRD, brief, dan plan terkait sebelum mengisi.
2. Turunkan test case dari requirement yang ada. Jangan menambah requirement baru.
3. Ikuti format field dan penamaan `TC-{DOMAIN}-{NNN}` tanpa mengubah ID yang sudah dipakai.
4. Jangan menandai `Pass` tanpa menjalankan test dan memverifikasi hasilnya.
5. Laporkan asumsi, gap cakupan, dan tes yang belum dijalankan di bagian 9.
