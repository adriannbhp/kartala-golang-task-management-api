# Kartala Task Management API

REST API modern untuk sistem manajemen tugas (task management) yang dibangun menggunakan Go dengan arsitektur bersih (Clean Architecture). Didesain untuk performa tinggi, keamanan berlapis, dan kemudahan deployment.

---

## 🚀 Fitur Utama

- **User Authentication** – Register & Login menggunakan JWT (Access 15m + Refresh 7d).
- **Task Management** – Operasi CRUD lengkap untuk tugas pribadi user.
- **Keamanan Berlapis**:
    - **Rate Limiting**: Proteksi spam berbasis IP (2 req/sec, burst 5).
    - **API Key**: Semua endpoint `/api/v1` wajib menggunakan `X-API-KEY`.
    - **Hardening**: Pembatasan ukuran request (2MB) & timeout otomatis (60s).
    - **Security Headers**: Dilengkapi CSP, HSTS, No-Sniff, dll.
- **Validation**: Validasi input ketat di layer DTO menggunakan `ozzo-validation`.
- **Infrastructure**:
    - **GORM**: ORM canggih untuk interaksi database.
    - **Atlas**: Manajemen migrasi database secara deklaratif & aman.
    - **Cloud Native**: Siap dideploy ke Docker, Docker Compose, atau Google Cloud Run.
- **Interactive Docs**: Swagger UI otomatis ter-generate.

---

## 📋 Persyaratan Sistem

- **Go**: Versi 1.25 atau lebih baru.
- **PostgreSQL**: Versi 14 atau lebih baru.
- **Atlas CLI**: [Install Guide](https://atlasgo.io/getting-started/) (untuk migrasi).
- **Make**: Untuk menjalankan shortcut perintah.
- **Docker & Docker Compose**: (Opsional) Untuk containerization.

---

## ⚙️ Opsi Setup & Instalasi

Pilih metode yang paling sesuai dengan kebutuhan Anda:

### 1. Lokal (Tanpa Docker) - *Cepat & Mudah*

1. **Clone & Install**:
   ```bash
   git clone https://github.com/adriannbhp/kartala-golang-task-management-api.git
   cd kartala-golang-task-management-api
   go mod tidy
   ```

2. **Environment**:
   Salin `.env.example` menjadi `.env` dan isi data database Anda:
   ```bash
   cp .env.example .env
   ```

3. **Database Migration**:
   ```bash
   make migrate-up-atlas
   ```

4. **Run Application**:
   ```bash
   make dev   # Menggunakan live-reload (Air)
   # ATAU
   make run   # Langsung jalankan tanpa reload
   ```

---

### 2. Docker Compose - *Satu Command Selesai*

Metode ini akan menjalankan API di dalam container.

1. **Setup Env**: Pastikan `.env` sudah terisi dengan benar.
2. **Start Services**:
   ```bash
   make docker-up
   ```
3. **Check Logs**:
   ```bash
   make docker-logs
   ```
4. **Stop**:
   ```bash
   make docker-down
   ```

---

### 3. Google Cloud Run - *Siap Produksi*

Proyek ini sudah dilengkapi pipeline otomatis untuk deployment ke GCP.

1. **Authentication**: Login ke gcloud CLI.
   ```bash
   gcloud auth login
   gcloud auth configure-docker
   ```

2. **GCS Credentials (Opsional)**:
   Jika Anda membutuhkan akses ke Google Cloud Storage, encode Service Account JSON Anda ke Base64 dan masukkan ke `.env` pada field `GCS_CREDENTIALS_BASE64`. Jika dikosongkan, aplikasi tetap akan berjalan normal.

3. **Deploy (Metode Rekomendasi)**:
   Gunakan Cloud Build agar proses build image tidak membebani komputer Anda:
   ```bash
   make deploy-cloudbuild GCP_PROJECT=id-project-anda GCP_REGION=asia-southeast2 SERVICE_NAME=kartala-api
   ```

---

## � Dokumentasi & Testing

### Swagger UI
Akses dokumentasi interaktif di: `http://localhost:8080/swagger/index.html`

### Menjalankan Test
Aplikasi ini memiliki cakupan unit test yang luas dan integration test yang sangat stabil (**Total Coverage: 99.9%**):
```bash
make test-unit           # Jalankan test logika bisnis (Cepat)
make test-integration    # Jalankan integration test (Butuh DB Test)
make test-coverage-html  # Lihat laporan visual cakupan kode
```

---

## 🔒 Keamanan & Header Penting

Setiap request ke API **WAJIB** menyertakan header berikut:

| Header | Nilai | Keterangan |
|---|---|---|
| `X-API-KEY` | *API_KEY_ANDA* | Token rahasia internal. |
| `Authorization` | `Bearer <JWT_TOKEN>` | Digunakan untuk endpoint privat (/tasks, /users/me). |
| `Content-Type` | `application/json` | Format data standar. |

Jika terjadi rate limiting, Anda akan menerima response `429 Too Many Requests` dengan format standar project.

---

## �️ Ringkasan Makefile (Cheat Sheet)

| Bagian | Command | Kegunaan |
|---|---|---|
| **Dev** | `make dev` | Live reload development. |
| **Data** | `make migrate-up-atlas` | Sinkronisasi database. |
| **Data** | `make seed` | Isi database dengan data contoh. |
| **Test** | `make test` | Jalankan seluruh suite pengujian. |
| **Clean** | `make clean` | Bersihkan build & log. |
| **Cloud** | `make deploy-cloudbuild` | Build & Deploy ke Cloud Run. |

---

**© 2026 Adrian Bimo Hernawan Pratama**
