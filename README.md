# Kartala Go Task Management API

API Manajemen Task sederhana berbasis Golang dengan Clean Architecture.

## 🚀 Fitur Utama
- **CRUD Tasks**: Kelola daftar tugas harian Anda.
- **User Authentication**: Register & Login menggunakan JWT.
- **API Key Security**: Pengamanan endpoint menggunakan API Key.
- **Clean Architecture**: Kode yang terstruktur, testable, dan mudah di-maintain.
- **Database Migrations**: Sinkronisasi schema database menggunakan **Atlas**.
- **Interactive Documentation**: Swagger UI terintegrasi.

## 🛠️ Tech Stack
- **Language**: Go 1.25+
- **Framework**: Gin Gonic
- **ORM**: GORM (PostgreSQL)
- **Migration**: Atlas
- **Security**: JWT & Bcrypt
- **Documentation**: Swaggo (Swagger 2.0)
- **Testing**: Testify & SQLMock

## 🏁 Memulai

### 1. Prasyarat
- Go installed
- PostgreSQL (NeonDB direkomendasikan)
- Atlas CLI

### 2. Konfigurasi Environment
Buat file `.env` di root directory (copy dari `.env.example`):
```env
DB_HOST=ep-frosty-wind-a1eqid67-pooler.ap-southeast-1.aws.neon.tech
DB_PORT=5432
DB_USER=neondb_owner
DB_PASSWORD=npg_Ypv9uTMCmA4a
DB_NAME=kartala_dev
DB_SSLMODE=require

# URL format used by Atlas and Makefile (harus sinkron dengan di atas)
DB_URL=postgres://neondb_owner:npg_Ypv9uTMCmA4a@ep-frosty-wind-a1eqid67-pooler.ap-southeast-1.aws.neon.tech:5432/kartala_dev?sslmode=require

JWT_SECRET=rahasia-anda
API_KEY=dev-api-key-123
```

> [!IMPORTANT]
> Pastikan `DB_NAME` dan database di dalam `DB_URL` sama (misal: `kartala_dev`). Proyek ini menggunakan `DB_NAME` untuk aplikasi dan `DB_URL` untuk tool Atlas.

### 3. Database Migrations (Atlas)
Proyek ini menggunakan Atlas untuk manajemen database deklaratif.

```bash
# Cek status migrasi
make migrate-status

# Jalankan migrasi ke database (Target Up)
make migrate-up-atlas

# Update checksum jika mengedit file migrasi
make migrate-hash
```

### 4. Menjalankan Seeder
Populate database dengan data awal untuk pengembangan:
```bash
go run cmd/seeder/main.go
```

### 5. Menjalankan Aplikasi
Gunakan `Air` untuk live-reload saat development:
```bash
# Install Air jika belum ada
# go install github.com/air-verse/air@latest+

make dev
```

## 📖 Dokumentasi API
Setelah aplikasi berjalan, buka documentation Swagger di:
[http://localhost:8082/swagger/index.html](http://localhost:8082/swagger/index.html)

Anda juga bisa melihat detail endpoint di [docs/API_DOCS.md](./docs/API_DOCS.md).

## 🧪 Testing
```bash
# Jalankan semua test
make test

# Unit test saja
make test-unit
```

## 🏗️ Folder Structure
- `cmd/`: Entry point aplikasi.
- `internal/`: Logika bisnis (Clean Architecture layers).
- `pkg/`: Utility packages (Database, Logger, Response).
- `migrations/`: File migrasi Atlas.
- `docs/`: Dokumentasi API.
