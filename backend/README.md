# ⚙️ DeepPhotos Backend Microservice

A high-performance, lightweight REST API built in **Go 1.22+** with embedded **SQLite 3** (`photos.db`) metadata storage, self-hosted **MinIO S3** object storage, and JWT token authentication.

---

## 🧭 Table of Contents
1. [Architecture & Design Rationale](#-architecture--design-rationale)
2. [Database Schema & B-Tree Indexes](#-database-schema--b-tree-indexes)
3. [MinIO S3 Object Storage Layout](#-minio-s3-object-storage-layout)
4. [Package Directory Structure](#-package-directory-structure)
5. [Complete REST API Reference](#-complete-rest-api-reference)
6. [Environment Variables Configuration](#-environment-variables-configuration)
7. [Build & Execution Commands](#-build--execution-commands)

---

## 🏗️ Architecture & Design Rationale

DeepPhotos decouples **structured metadata queries** from **heavy binary media blobs**:

```text
 ┌──────────────────────┐
 │  SvelteKit Frontend  │
 └──────────┬───────────┘
            │
            │ REST API (JSON / HTTP Stream)
            ▼
 ┌──────────────────────┐
 │    Go Backend API    │  <-- Chi v5 Router (Auth Middleware)
 └──────┬────────┬──────┘
        │        │
  Metadata      Files
        ▼        ▼
  ┌──────────┐ ┌──────────┐
  │  SQLite  │ │  MinIO   │
  │photos.db │ │  Bucket  │
  └──────────┘ └──────────┘
```

### Why Go + SQLite 3 (WAL Mode) + MinIO?
1. **CGo-Free SQLite Engine (`modernc.org/sqlite`)**: Compiles into a single static Go binary with zero external OS library dependencies.
2. **WAL (Write-Ahead Logging) Mode**: `PRAGMA journal_mode=WAL;` allows concurrent readers to query metadata while write operations (uploads/updates) take place without locking.
3. **B-Tree Indexing for 500,000+ Records**: Indexed columns allow instant `< 1ms` query lookups over 500,000 metadata rows.
4. **Hierarchical S3 Storage**: Prevents OS file system degradation by storing files under `image/`, `video/`, `document/`, `lockedfolder/`, `thumbnails/`, and `avatars/` partitions.

---

## 📊 Database Schema & B-Tree Indexes

The database schema (`data/photos.db`) contains 7 core tables:

### 1. `users` Table
Stores user accounts, bcrypt hashed passwords, roles (*Administrator*, *Editor*, *Viewer*), and avatar URLs.

### 2. `photos` Table
Stores primary media metadata:
```sql
CREATE TABLE IF NOT EXISTS photos (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    filename TEXT NOT NULL,
    storage_path TEXT NOT NULL,
    thumbnail_path TEXT,
    file_type TEXT NOT NULL DEFAULT 'image',
    mime_type TEXT NOT NULL,
    size INTEGER NOT NULL,
    width INTEGER,
    height INTEGER,
    is_favorite BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    locked_folder_id TEXT,
    exif_make TEXT,
    exif_model TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 3. High-Performance B-Tree Indexes (Optimized for 500,000+ Rows)
```sql
CREATE INDEX IF NOT EXISTS idx_photos_type_deleted ON photos(file_type, is_deleted);
CREATE INDEX IF NOT EXISTS idx_photos_user_id ON photos(user_id);
CREATE INDEX IF NOT EXISTS idx_photos_created_at ON photos(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_photos_locked_folder ON photos(locked_folder_id);
CREATE INDEX IF NOT EXISTS idx_photos_favorite ON photos(is_favorite);
```

---

## 🔒 MinIO S3 Object Storage Layout

Media and profile photos are organized in MinIO / Local Disk Storage using 3-level UUIDv4 subfolder partitioning:

```text
deepphotos/
├── image/
│   └── {uuid4_1}/{uuid4_2}/{uuid4_3}/photo.jpg
├── video/
│   └── {uuid4_1}/{uuid4_2}/{uuid4_3}/video.mp4
├── document/
│   └── {uuid4_1}/{uuid4_2}/{uuid4_3}/document.pdf
├── lockedfolder/
│   └── {uuid4_1}/{uuid4_2}/{uuid4_3}/secret.png
├── thumbnails/
│   └── {category}/{uuid4_1}/{uuid4_2}/{uuid4_3}/thumb_photo.webp
└── avatars/                  <-- Dedicated User Profile Avatar Storage Bucket Path
    └── {user_id}/{filename}
```

---

## 📁 Package Directory Structure

```text
backend/
├── cmd/
│   └── server/
│       └── main.go           <-- Application entrypoint & Chi v5 route configuration
├── internal/
│   ├── config/
│   │   └── config.go         <-- Environment variable loading & defaults
│   ├── database/
│   │   └── db.go             <-- SQLite initialization, WAL mode, schema migrations & indexes
│   ├── handler/
│   │   ├── auth.go           <-- Authentication handlers
│   │   ├── photos.go         <-- Media ingestion, list, stream, update, & delete handlers
│   │   ├── albums.go         <-- Album management handlers
│   │   ├── memories.go       <-- Custom memory collections handlers
│   │   ├── users.go          <-- User management & avatar upload handlers
│   │   ├── locked.go         <-- Passcode-protected vault handlers
│   │   ├── audit.go          <-- Security audit log handlers
│   │   └── health.go         <-- Health check endpoint handler
│   ├── middleware/
│   │   └── auth.go           <-- JWT bearer token & query parameter auth middleware
│   ├── model/
│   │   └── models.go         <-- Go structs & JSON DTO models
│   ├── repository/
│   │   ├── photo_repo.go     <-- SQLite database queries for media
│   │   ├── user_repo.go      <-- User repository queries
│   │   ├── album_repo.go     <-- Album repository queries
│   │   ├── memory_repo.go    <-- Memory repository queries
│   │   ├── locked_repo.go    <-- Vault repository queries
│   │   └── audit_repo.go     <-- Audit repository queries
│   ├── service/
│   │   ├── auth_service.go   <-- Authentication business logic & JWT signing
│   │   └── photo_service.go  <-- File upload processing & EXIF extraction
│   ├── storage/
│   │   └── minio.go          <-- MinIO S3 object storage client & key generation
│   └── utils/
│       └── env.go            <-- Environment helper utilities
└── go.mod                    <-- Go module dependencies
```

---

## 🛠️ Complete REST API Reference

### 1. Media Operations (`/api/media`)

| Method | Endpoint | Query / Body Params | Description |
|---|---|---|---|
| `GET` | `/api/media` | `?type=gallery`, `?type=document`, `?type=photos`, `?type=video`, `?deleted=false`, `?search=query` | List media with optional filters |
| `POST` | `/api/media` | `multipart/form-data` (`files`) | Upload single or multiple media files |
| `POST` | `/api/media/upload` | `multipart/form-data` (`files`) | Upload single or multiple media files |
| `POST` | `/api/media/upload-url` | `{"url": "https://..."}` | Ingest photo/video from web URL |
| `GET` | `/api/media/{id}` | - | Fetch single media metadata & EXIF |
| `PUT` | `/api/media/{id}` | `{"title":"...", "is_favorite":true, "is_deleted":true, "locked_folder_id":"..."}` | Update media properties |
| `DELETE` | `/api/media/{id}` | - | Permanently purge single media file |
| `POST` | `/api/media/batch-delete` | `{"ids": ["id1", "id2"]}` | Batch soft-delete or purge items |
| `POST` | `/api/media/batch-restore` | `{"ids": ["id1", "id2"]}` | Batch restore items from bin |
| `GET` | `/api/media/{id}/file` | `?token=<jwt>` | Stream high-res original file |
| `GET` | `/api/media/{id}/thumbnail` | `?token=<jwt>` | Stream WebP thumbnail image |

### 2. Gallery & Documents Aliases
* `GET /api/gallery` — Equivalent to `GET /api/media?type=gallery&deleted=false`.
* `GET /api/documents` — Equivalent to `GET /api/media?type=document&deleted=false`.

### 3. User Avatars (`/api/users/{id}/avatar`)
* `POST /api/users/{id}/avatar` — Multipart upload saving avatar to `avatars/{user_id}/{filename}` in MinIO.
* `GET /api/users/{id}/avatar` — Streams user profile avatar image.

---

## ⚡ Environment Variables Configuration

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP API listening port |
| `DB_PATH` | `data/photos.db` | SQLite database file location |
| `MINIO_ENDPOINT` | `localhost:9000` | MinIO storage endpoint |
| `MINIO_ACCESS_KEY` | `minioadmin` | MinIO access key |
| `MINIO_SECRET_KEY` | `minioadmin` | MinIO secret key |
| `MINIO_BUCKET` | `deepphotos` | MinIO storage bucket name |
| `JWT_SECRET` | `deepphotos-secret-key-2026` | HMAC signing secret for JWT tokens |

---

## 🚀 Build & Execution Commands

```bash
# Build binary
go build -o server ./cmd/server

# Run server
./server
```
