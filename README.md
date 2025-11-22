# KerjaKuy – Project Management SaaS Backend (Go + Gin + GORM + PostgreSQL)

Backend multi-tenant untuk manajemen proyek dan kolaborasi (board/kanban, tugas, komentar, lampiran) dengan autentikasi JWT + sesi. Arsitektur sekarang berupa **modular monolith berbasis feature** (auth, user, workspace, project, task) agar komposisi dependency jelas tanpa mengubah logika yang sudah ada.

## Fitur Singkat
- Auth & User: register, login, refresh, logout, JWT + session store, profil user.
- Workspace: create/update, role member (owner/admin/member), invite/remove member.
- Project/Board/Column/Task: CRUD lengkap, assignee, komentar, lampiran.
- Health: `/api/v1/ping`.

## Arsitektur Singkat
- **Modules**: `internal/modules/{auth,user,workspace,project,task}` merangkai repo/service/handler dan mendaftarkan route-nya.
- **Composition root**: `internal/app/app.go` memuat config, init DB, bangun modul, dan menyerahkan router ke `cmd/web`.
- **Layer**: handler (Gin) → service (bisnis) → repository (GORM) → models.
- **Config & infra**: `pkg/config` (env/.env), `pkg/database` (Postgres via GORM).

Struktur utama:
```
cmd/            # Entrypoint (web server, migrasi)
internal/
  app/          # Composition root modular monolith
  modules/      # Modul feature (auth, user, workspace, project, task)
  handler/      # Handler HTTP
  service/      # Service bisnis
  repository/   # Repository GORM
  models/       # Entity GORM
  router/v1/    # Router API v1 (compose modul)
pkg/
  config/       # Load env
  database/     # Init Postgres
db/migrations/  # SQL migrasi (manual)
docs/           # OpenAPI/Postman
```

## Jalankan secara lokal
Prasyarat: Go 1.22+, PostgreSQL berjalan, dan variabel lingkungan terisi.

1) Salin `.env.example` (jika ada) atau set env berikut:
```
DB_HOST=127.0.0.1
DB_USER=postgres
DB_PASS=postgres
DB_NAME=kerjakuy
DB_PORT=5432
DB_SSL=disable
APP_PORT=8080
GIN_MODE=debug
JWT_SECRET=supersecret
JWT_ISSUER=kerjakuy
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
```
2) Jalankan migrasi:
```
go run ./cmd/migrate
```
3) Jalankan server:
```
go run ./cmd/web
```
API tersedia di `http://localhost:8080/api/v1`.

## API Docs
- OpenAPI: `docs/openapi.yaml` (bisa di-import ke Swagger UI/Insomnia/Stoplight)
- Ping: `GET /api/v1/ping`

## Catatan Pengembangan
- Modul modular monolith: tambah feature cukup buat modul baru di `internal/modules` dan daftarkan di `internal/app/app.go`.
- Business logic tidak diubah saat modularisasi; semua handler/service/repo tetap sama, hanya wiring yang dipisah.
