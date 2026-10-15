# DeepPhotos — Self-Hosted Personal Media Cloud

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Svelte Version](https://img.shields.io/badge/Svelte-5.0_(Runes)-FF3E00?style=for-the-badge&logo=svelte)](https://svelte.dev)
[![SQLite Version](https://img.shields.io/badge/SQLite-WAL_Mode-003B57?style=for-the-badge&logo=sqlite)](https://www.sqlite.org)
[![MinIO](https://img.shields.io/badge/Storage-MinIO_S3-C72C48?style=for-the-badge&logo=minio)](https://min.io)
[![Tailwind CSS](https://img.shields.io/badge/Styling-Tailwind_CSS_v4-38BDF8?style=for-the-badge&logo=tailwindcss)](https://tailwindcss.com)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg?style=for-the-badge)](LICENSE)

**DeepPhotos** is a private, self-hosted photo, video, and document manager built as a privacy-focused Google Photos alternative. You own the database and the original files. It provides a media timeline, albums, memory collections, a passcode-protected vault, a recycle bin, multi-select batch actions, user administration, and automatic WebP thumbnails.

---

## Table of Contents
1. [Why DeepPhotos?](#why-deepphotos)
2. [Features](#features)
3. [Quick Start](#quick-start)
4. [Architecture & Tech Stack](#architecture--tech-stack)
5. [Project Structure](#project-structure)
6. [Configuration](#configuration)
7. [Authentication & Roles](#authentication--roles)
8. [How It Works (Data Flows)](#how-it-works-data-flows)
9. [Object Storage Layout](#object-storage-layout)
10. [Page-to-API Mapping](#page-to-api-mapping)
11. [REST API Reference](#rest-api-reference)
12. [Development](#development)
13. [Production Deployment](#production-deployment)
14. [Security Checklist](#security-checklist)
15. [Troubleshooting](#troubleshooting)
16. [Contributing](#contributing)
17. [License](#license)

---

## Why DeepPhotos?

Commercial photo services (Google Photos, iCloud, Amazon Photos) impose storage limits, lock you into subscriptions, and process your private photos on their servers.

DeepPhotos gives you:

* **Privacy & ownership** — your metadata lives in a local SQLite file (`photos.db`) and your original, uncompressed files live in your own MinIO bucket (or on local disk). No third-party tracking, no mandatory subscription.
* **Fast metadata queries** — SQLite runs in Write-Ahead Logging (WAL) mode with indexes on the columns the gallery filters by (type, deleted state, owner, date, vault folder, favorite).
* **Decoupled architecture** — a lightweight Go REST API paired with a reactive Svelte 5 web client.
* **Broad media support** — images, videos (`MP4`, `MOV`, `WebM`, `MKV`, `AVI`), and documents (`PDF`, `DOC`/`DOCX`, `TXT`, `ZIP`).

---

## Features

| Area | What you get |
|---|---|
| **Gallery** | Photo and video timeline with multi-select, favorites, and batch move-to-bin. |
| **Media viewer** | Full-screen lightbox with zoom controls, a details drawer (camera model, dimensions, file size), video playback with `Space` to play/pause, and a "resized to fit" indicator. |
| **Albums** | Create, rename, and delete albums; add photos to them. |
| **Memories** | Named collections with optional descriptions; add or remove photos. |
| **Locked vault** | Passcode-protected folders (4-digit passcode, stored as a bcrypt hash). |
| **Bin** | Soft-delete, batch restore, and permanent purge (removes the object and its thumbnail from storage). |
| **Documents** | Dedicated view for PDFs and other documents. |
| **Uploads** | Drag-and-drop, file picker, or import from a URL. |
| **Users & roles** | Admin-managed accounts, role changes, password changes, avatars, and an authentication audit log. |
| **Theming** | Light and dark themes across every view. |

---

## Quick Start

### Option A — Docker (recommended)

Prerequisite: Docker with Compose v2.

```bash
git clone https://github.com/deep10dr/deepphotos.git
cd deepphotos

# Optional: customize credentials and secrets
cp .env.template .env

# Mac / Linux
./scripts/mac-linux/docker_run.sh

# Windows
scripts\windows\docker_run.bat
```

Then open:

| Service | URL |
|---|---|
| Web app | http://localhost:5173 |
| REST API | http://localhost:8080 |
| MinIO console | http://localhost:9001 |

Sign in with the admin account from your `.env` (default: `admin@deepphotos.local` / `deepphotos2026`). **Change these before exposing the app to anyone** — see the [Security Checklist](#security-checklist).

Stop or restart the stack:

```bash
./scripts/mac-linux/docker_run.sh stop
./scripts/mac-linux/docker_run.sh restart
```

### Option B — Local development

Prerequisites: **Go 1.25+**, **Node.js 18+** with `npm`, and (optionally) Docker for MinIO.

```bash
# Mac / Linux
./scripts/mac-linux/dev.sh

# Windows
scripts\windows\dev.bat
```

The script starts MinIO through Docker when it is available, runs the Go backend on `:8080`, and starts the Vite dev server on `:5173`. If Docker is not running, the backend falls back to **local disk storage** under `backend/data/storage/`.

---

## Architecture & Tech Stack

```text
 ┌─────────────────────────────────────────────────────────────────┐
 │                   BROWSER CLIENT (SvelteKit)                    │
 │  Svelte 5 Runes • Tailwind CSS v4 • Lucide Icons • Light/Dark   │
 └────────────────────────────────┬────────────────────────────────┘
                                  │
                       JSON REST API & JWT Bearer
                                  │
 ┌────────────────────────────────▼────────────────────────────────┐
 │                      GO BACKEND MICROSERVICE                    │
 │  Chi v5 Router • JWT Auth • bcrypt • WebP Thumbnail Generation  │
 └─────────────────┬──────────────────────────────┬────────────────┘
                   │                              │
        SQLite (CGo-free driver)         MinIO S3 / Local Disk
                   │                              │
 ┌─────────────────▼──────────────┐ ┌─────────────▼────────────────┐
 │    SQLITE 3 (data/photos.db)   │ │  MINIO BUCKET (deepphotos)   │
 │ WAL Mode • Indexed queries     │ │  image/ video/ document/ ... │
 └────────────────────────────────┘ └──────────────────────────────┘
```

| Component | Technology | Notes |
|---|---|---|
| **Frontend** | SvelteKit, Svelte 5 Runes | `$state` / `$derived` fine-grained reactivity; `adapter-node` for production. |
| **Styling** | Tailwind CSS v4, bits-ui, Lucide | Utility-first, with `dark:` variants for theming. |
| **Backend** | Go, Chi v5, chi/cors | Stateless REST API. |
| **Database** | SQLite via `modernc.org/sqlite` | Pure Go (no CGo), WAL mode, foreign keys on. |
| **Object storage** | MinIO (S3 API) | Automatic local-disk fallback when MinIO is unreachable. |
| **Auth** | JWT + bcrypt | Query-parameter token supported for `<img>` / `<video>` streaming. |

---

## Project Structure

```text
deepphotos/
├── .env.template          # Environment variable template
├── backend/               # Go REST API
│   ├── cmd/server/        #   Entry point and route wiring (main.go)
│   └── internal/
│       ├── config/        #   Environment-based configuration
│       ├── database/      #   SQLite setup, schema, indexes, admin seed
│       ├── handler/       #   HTTP handlers (auth, media, albums, memories, vault, users, audit, health)
│       ├── middleware/    #   JWT auth and admin-only guards
│       ├── model/         #   Request/response and domain types
│       ├── repository/    #   Data-access layer
│       ├── service/       #   Business logic (media ingestion, auth)
│       ├── storage/       #   MinIO client and local-disk fallback
│       └── utils/         #   Shared helpers
├── frontend/              # SvelteKit web app
│   └── src/
│       ├── lib/           #   API client, shared state, notifications, components
│       └── routes/(app)/  #   gallery, documents, albums, memories, vault, bin, profile, settings
├── docker/                # docker-compose.yaml and infrastructure notes
└── scripts/               # mac-linux/ and windows/ helper scripts
```

Each subproject also has its own README with more detail: [`backend/README.md`](backend/README.md), [`frontend/README.md`](frontend/README.md), and [`docker/README.md`](docker/README.md).

### Helper scripts

Each script exists in `scripts/mac-linux/` (`.sh`) and `scripts/windows/` (`.bat`):

| Script | Purpose |
|---|---|
| `dev` | Local development: MinIO (if Docker is available), Go backend, Vite dev server. Pass `--debug` (Mac/Linux) to run the backend under the Delve debugger on port 4000. |
| `build` | Build the backend binary and frontend assets. |
| `run` | Run the built application. |
| `build_and_run` | Build, then run. |
| `docker_build` | Build the Docker images. |
| `docker_run` | Start the Docker stack. Accepts `stop` and `restart`. |

---

## Configuration

The backend reads environment variables and falls back to the defaults below. Copy `.env.template` to `.env` to override them.

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Backend HTTP port. |
| `DB_PATH` | `data/photos.db` | SQLite database file. |
| `MINIO_ENDPOINT` | `127.0.0.1:9000` | MinIO host and port (`minio:9000` inside Docker). |
| `MINIO_ACCESS_KEY` | `minioadmin` | MinIO access key. |
| `MINIO_SECRET_KEY` | `minioadmin` | MinIO secret key. |
| `MINIO_BUCKET` | `deepphotos` | Bucket for media, thumbnails, and avatars. |
| `MINIO_USE_SSL` | `false` | Set to `true` to connect to MinIO over TLS. |
| `JWT_SECRET` | `deepphotos-secret-key-2026` | Secret used to sign login tokens. **Change in production.** |
| `ADMIN_EMAIL` | `admin@deepphotos.local` | Email of the seeded administrator. |
| `ADMIN_PASSWORD` | `deepphotos2026` | Password of the seeded administrator. **Change in production.** |
| `VITE_API_BASE_URL` | `http://localhost:8080` | Frontend: base URL of the backend API. |
| `VITE_DEV_USERNAME` / `VITE_DEV_PASSWORD` | — | Frontend: development convenience credentials used in `.env.template`. Do not set these in production. |

---

## Authentication & Roles

* `POST /api/auth/login` returns a JWT. Send it as `Authorization: Bearer <token>` on every other request.
* Browsers cannot attach headers to `<img>` and `<video>` tags, so media endpoints also accept `?token=<jwt>`.
* A default administrator is created on first run from `ADMIN_EMAIL` and `ADMIN_PASSWORD`. On every later start, the backend resets that account's password to the current `ADMIN_PASSWORD`, so a password changed through the UI for that account does not survive a restart. Change `ADMIN_PASSWORD` in your environment instead.
* Roles are `Administrator`, `Editor`, and `Viewer` (new users default to `Viewer`).
* Administrator-only actions: list and create users, change roles, delete users, and read audit logs. Any user can update their own profile, password, and avatar.

---

## How It Works (Data Flows)

### 1. Ingestion & upload
```text
User file drop / upload
   │
   ├──> 1. Frontend sends a multipart payload (or a URL for import)
   ├──> 2. POST /api/media/upload -> backend classifies the file (image / video / document)
   ├──> 3. Backend builds a partitioned key: image/{uuid1}/{uuid2}/{uuid3}/photo.jpg
   ├──> 4. Original streamed to storage; WebP thumbnail generated for images
   ├──> 5. EXIF metadata parsed where available
   └──> 6. Metadata row saved in SQLite -> JSON response -> Svelte refreshes the view
```

### 2. Authenticated media streaming
1. The Svelte client calls `getMediaUrl("/api/media/{id}/file")`.
2. `getMediaUrl()` appends the active JWT: `/api/media/{id}/file?token=<jwt>`.
3. Backend middleware verifies the token and streams the file.

### 3. Delete lifecycle
Items first move to the **Bin** (soft delete) and can be restored. Permanently deleting from the Bin removes the original object and its thumbnail from storage.

---

## Object Storage Layout

Media is partitioned by three levels of UUIDv4 folders to avoid very large directories:

```text
deepphotos/
├── image/
│   └── {uuid4_1}/{uuid4_2}/{uuid4_3}/photo.jpg
├── video/
│   └── {uuid4_1}/{uuid4_2}/{uuid4_3}/video.mp4
├── document/
│   └── {uuid4_1}/{uuid4_2}/{uuid4_3}/document.pdf
├── lockedfolder/
│   └── {uuid4_1}/{uuid4_2}/{uuid4_3}/vault_file.png
├── thumbnails/
│   └── {category}/{uuid4_1}/{uuid4_2}/{uuid4_3}/thumb_photo.webp
└── avatars/
    └── {user_id}/{filename}
```

When MinIO is unreachable, the backend stores objects on local disk under `data/storage/` instead.

---

## Page-to-API Mapping

Which backend endpoints each frontend view calls:

| Frontend Page / Component | Route | Backend API Endpoint(s) Called | Purpose |
|---|---|---|---|
| **Gallery** | `/gallery` | `GET /api/media?type=gallery&deleted=false`<br>`GET /api/locked-folders`<br>`PUT /api/media/{id}`<br>`POST /api/media/batch-delete`<br>`POST /api/memories`<br>`POST /api/memories/{id}/photos` | Photo/video timeline, selection mode, batch move to bin, add to memories. |
| **Documents** | `/documents` | `GET /api/media?type=document&deleted=false`<br>`DELETE /api/media/{id}` | List, download, and delete documents. |
| **Memories** | `/memories` | `GET /api/memories`<br>`POST /api/memories`<br>`GET /api/memories/{id}`<br>`DELETE /api/memories/{id}`<br>`DELETE /api/memories/{id}/photos/{photoId}`<br>`GET /api/media?type=gallery&deleted=false` | Create and browse memory collections. |
| **Albums** | `/albums` | `GET /api/albums`<br>`POST /api/albums`<br>`GET /api/albums/{id}`<br>`PUT /api/albums/{id}`<br>`DELETE /api/albums/{id}`<br>`POST /api/albums/{id}/photos`<br>`GET /api/media?type=gallery&deleted=false` | Manage albums and their photos. |
| **Locked Vault** | `/vault` | `GET /api/locked-folders`<br>`POST /api/locked-folders`<br>`POST /api/locked-folders/{id}/verify`<br>`DELETE /api/locked-folders/{id}`<br>`GET /api/media?locked_folder_id={id}` | Manage passcode-protected folders and view their media. |
| **Bin** | `/bin` | `GET /api/media?deleted=true`<br>`POST /api/media/batch-restore`<br>`POST /api/media/batch-delete` | Restore or permanently purge deleted items. |
| **Profile & Admin** | `/profile` | `GET /api/users`<br>`POST /api/users`<br>`PUT /api/users/{id}`<br>`PUT /api/users/{id}/password`<br>`PUT /api/users/{id}/role`<br>`DELETE /api/users/{id}`<br>`POST /api/users/{id}/avatar`<br>`GET /api/users/{id}/avatar`<br>`GET /api/audit-logs` | Profile editing, avatars, user management, and audit logs (admin). |
| **TopHeader** | Global | `POST /api/media/upload`<br>`POST /api/media/upload-url`<br>`GET /api/users/{id}/avatar` | Uploads, URL import, and avatar display. |
| **MediaViewer** | Global lightbox | `GET /api/media/{id}/file`<br>`GET /api/media/{id}/thumbnail`<br>`PUT /api/media/{id}`<br>`DELETE /api/media/{id}` | Full-resolution streaming, favorites, move to vault, delete. |

---

## REST API Reference

All routes except `/api/health` and `/api/auth/login` require authentication.

### Health & authentication
* `GET /api/health` — System status, database mode, and storage connectivity.
* `POST /api/auth/login` — Body: `{"email":"...","password":"..."}`. Returns a JWT.

### Media (`/api/media`)
* `GET /api/media` — List media. Filters: `type` (`gallery`, `photos`, `video`, `document`), `deleted`, `search`, `locked_folder_id`.
* `POST /api/media` and `POST /api/media/upload` — Multipart upload.
* `POST /api/media/upload-url` — Import from a URL: `{"url":"https://..."}`.
* `GET /api/media/{id}` — Single item metadata and EXIF details.
* `PUT /api/media/{id}` — Update title, favorite flag, soft-delete state, or vault folder.
* `DELETE /api/media/{id}` — Permanently delete one item.
* `POST /api/media/batch-delete` — Body: `{"ids":["..."]}`.
* `POST /api/media/batch-restore` — Body: `{"ids":["..."]}`.
* `GET /api/media/{id}/file` — Stream the original file (`?token=<jwt>` for browser tags).
* `GET /api/media/{id}/thumbnail` — Stream the WebP thumbnail.

Convenience aliases: `GET /api/gallery`, `GET /api/documents`, and `GET /api/photos`.

### Albums & memories
* `GET /api/albums`, `POST /api/albums` — List and create.
* `GET /api/albums/{id}`, `PUT /api/albums/{id}`, `DELETE /api/albums/{id}` — Detail, update, delete.
* `POST /api/albums/{id}/photos` — Add photos to an album.
* `GET /api/memories`, `POST /api/memories` — List and create.
* `GET /api/memories/{id}`, `DELETE /api/memories/{id}` — Detail and delete.
* `POST /api/memories/{id}/photos`, `DELETE /api/memories/{id}/photos/{photoId}` — Add or remove photos.

### Locked vault
* `GET /api/locked-folders`, `POST /api/locked-folders` — List and create.
* `POST /api/locked-folders/{id}/verify` — Body: `{"passcode":"1234"}`.
* `DELETE /api/locked-folders/{id}` — Delete a folder.

### Users & audit
* `PUT /api/users/{id}` — Update profile (own account, or any account as admin).
* `PUT /api/users/{id}/password` — Change password.
* `POST /api/users/{id}/avatar`, `GET /api/users/{id}/avatar` — Upload (10 MB limit) and stream avatars.
* `GET /api/users`, `POST /api/users`, `PUT /api/users/{id}/role`, `DELETE /api/users/{id}` — *Administrator only.*
* `GET /api/audit-logs` — Authentication audit log *(Administrator only)*.

---

## Development

```bash
# Terminal 1 — backend (http://localhost:8080)
cd backend
go run ./cmd/server

# Terminal 2 — frontend (http://localhost:5173)
cd frontend
npm install
npm run dev
```

Useful commands:

```bash
# Frontend type and Svelte checks
cd frontend && npm run check

# Debug the backend with Delve (Mac/Linux)
./scripts/mac-linux/dev.sh --debug     # debugger listens on :4000

# Start only MinIO
docker compose -f docker/docker-compose.yaml up -d minio
```

---

## Production Deployment

### With Docker Compose

```bash
cp .env.template .env      # set strong secrets first
docker compose -f docker/docker-compose.yaml up -d --build
```

Data is kept in two named volumes: `minio_data` (media) and `sqlite_data` (the database). Back both up regularly.

### Without Docker

```bash
# Build the backend
cd backend
go build -o server ./cmd/server

# Build the frontend
cd ../frontend
npm install
npm run build

# Run the backend (serve the frontend build with Node or a reverse proxy)
./backend/server
```

Put the app behind a reverse proxy that terminates HTTPS, and set `VITE_API_BASE_URL` to the public API address before building the frontend.

---

## Security Checklist

Before exposing DeepPhotos beyond your own machine:

- [ ] Set a long, random `JWT_SECRET`.
- [ ] Change `ADMIN_EMAIL` and `ADMIN_PASSWORD`.
- [ ] Change `MINIO_ACCESS_KEY` and `MINIO_SECRET_KEY`, and do not publish MinIO ports `9000`/`9001` publicly.
- [ ] Serve everything over HTTPS.
- [ ] Remove `VITE_DEV_USERNAME` and `VITE_DEV_PASSWORD` from production builds.
- [ ] Keep regular backups of the `minio_data` and `sqlite_data` volumes.

---

## Troubleshooting

| Symptom | Likely cause and fix |
|---|---|
| Backend logs "MinIO storage offline … local disk fallback active" | MinIO is not running or `MINIO_ENDPOINT` is wrong. Start it with `docker compose -f docker/docker-compose.yaml up -d minio`, or keep using local disk storage. |
| Images or videos do not load in the browser | The media URL is missing `?token=<jwt>`, or the token has expired. Sign in again. |
| Frontend cannot reach the API | Check that `VITE_API_BASE_URL` points to the running backend and that port `8080` is free. |
| Port 8080 already in use | Stop the other process. `dev.sh` stops an existing listener on `8080` automatically on Mac/Linux. |
| Admin password reverts after a restart | The backend resets the admin account's password to `ADMIN_PASSWORD` on every start. Update `ADMIN_PASSWORD` in `.env` and restart. |

---

## Contributing

1. Fork the repository and create a branch from `dev`.
2. Keep changes focused and run `npm run check` for frontend changes.
3. Open a pull request against `dev` describing what changed and why.

---

## License

This project is open-source under the [Apache License 2.0](LICENSE).
