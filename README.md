# 📸 DeepPhotos — Self-Hosted Personal Media Cloud

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Svelte Version](https://img.shields.io/badge/Svelte-5.0_(Runes)-FF3E00?style=for-the-badge&logo=svelte)](https://svelte.dev)
[![SQLite Version](https://img.shields.io/badge/SQLite-3.40+-003B57?style=for-the-badge&logo=sqlite)](https://www.sqlite.org)
[![MinIO](https://img.shields.io/badge/Storage-MinIO_S3-C72C48?style=for-the-badge&logo=minio)](https://min.io)
[![Tailwind CSS](https://img.shields.io/badge/Styling-Tailwind_CSS_v4-38BDF8?style=for-the-badge&logo=tailwindcss)](https://tailwindcss.com)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg?style=for-the-badge)](LICENSE)

**DeepPhotos** is a high-performance, private, self-hosted photo, video, and document management platform built as a privacy-focused Google Photos alternative. It delivers edge-to-edge media streaming, EXIF metadata extraction, custom memory collections, passcode-protected vault folders, multi-select batch actions, user administration, and automated WebP thumbnail generation.

---

## 🧭 Table of Contents
1. [Why DeepPhotos? (Motivation & Core Value)](#-why-deepphotos-motivation--core-value)
2. [What is DeepPhotos? (Tech Stack & Architecture)](#-what-is-deepphotos-tech-stack--architecture)
3. [How DeepPhotos Works (System Data Flows)](#-how-deepphotos-works-system-data-flows)
4. [Private Object Storage & Avatar Bucket Hierarchy](#-private-object-storage--avatar-bucket-hierarchy)
5. [Page-to-API Endpoint Mapping Matrix](#-page-to-api-endpoint-mapping-matrix)
6. [Complete REST API Specification](#-complete-rest-api-specification)
7. [Key Features & Capabilities](#-key-features--capabilities)
8. [Getting Started & Development Setup](#-getting-started--development-setup)
9. [Production Deployment Guide](#-production-deployment-guide)

---

## ❓ Why DeepPhotos? (Motivation & Core Value)

Commercial cloud photo services (Google Photos, Apple iCloud, Amazon Photos) impose strict storage limits, downsample high-resolution RAW images, lock users into subscription tiers, and scan private user photos for telemetry.

**DeepPhotos solves these problems by providing:**
* **100% Privacy & Ownership**: You own your database (`photos.db`) and original uncompressed media files in local MinIO storage. No third-party tracking or mandatory subscriptions.
* **Sub-Millisecond Metadata Performance**: SQLite 3 with Write-Ahead Logging (WAL) and B-Tree indexes serves 500,000+ photo metadata records in **< 1ms**.
* **Decoupled Architecture**: High-speed Go microservice for API execution paired with Svelte 5 Runes reactive web client.
* **Universal Media Support**: Handles high-res images (`JPEG`, `PNG`, `WebP`, `HEIC`, `GIF`), 4K videos (`MP4`, `MOV`, `WebM`), and documents (`PDF`, `DOCX`, `TXT`).

---

## 🛠️ What is DeepPhotos? (Tech Stack & Architecture)

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
 │  Chi v5 Router • JWT Auth • bcrypt Security • Image Converter   │
 └─────────────────┬──────────────────────────────┬────────────────┘
                   │                              │
        SQLite CGo-Free Metadata             MinIO S3 / Local Disk
                   │                              │
 ┌─────────────────▼──────────────┐ ┌─────────────▼────────────────┐
 │    SQLITE 3 (data/photos.db)   │ │  MINIO BUCKET (deepphotos)   │
 │ WAL Mode • 5 B-Tree Indexes    │ │  image/ video/ avatars/      │
 └────────────────────────────────┘ └──────────────────────────────┘
```

### Technology Breakdown

| Component | Technology | Purpose & Rationale |
|---|---|---|
| **Frontend Framework** | **SvelteKit (Svelte 5 Runes)** | Fine-grained reactivity via `$state` and `$derived`. Zero virtual-DOM overhead, compiling to minimal client JS. |
| **Styling & UI** | **Tailwind CSS v4** | Utility-first responsive design with dynamic `dark:` class support for full light/dark theme switching. |
| **Backend API** | **Go 1.22+ & Chi v5** | Lightweight, high-concurrency HTTP routing with standard library execution speeds (< 5ms response time). |
| **Database** | **SQLite 3 (modernc.org/sqlite)** | CGo-free embedded SQL engine running in **WAL mode** with custom B-Tree indexes handling 500k+ rows smoothly. |
| **Object Storage** | **MinIO S3 API** | High-speed local S3-compatible blob storage with UUIDv4 hierarchical partitioning and dedicated `avatars/` bucket. |
| **Authentication** | **JWT & bcrypt** | Stateless JSON Web Token authentication with query-parameter token streaming (`?token=...`) for HTML5 `<video>` elements. |

---

## 🔄 How DeepPhotos Works (System Data Flows)

### 1. Ingestion & Upload Flow
```text
User File Drop / Upload
   │
   ├──> 1. Frontend constructs FormData / multipart payload
   ├──> 2. POST /api/media/upload -> Backend validates file MIME type & extension
   ├──> 3. Backend generates 3-level UUIDv4 partition: image/{uuid1}/{uuid2}/{uuid3}/photo.jpg
   ├──> 4. Original binary streamed to MinIO; WebP thumbnail created under thumbnails/
   ├──> 5. EXIF data (camera model, width, height, ISO, aperture) parsed
   └──> 6. Metadata row saved in SQLite photos.db -> Returns JSON -> Svelte invalidates cache
```

### 2. Authenticated Video & Photo Streaming
Browser `<video src="...">` and `<img>` tags cannot natively include HTTP headers (`Authorization: Bearer <jwt>`). DeepPhotos resolves this via JWT query parameters:
1. Svelte client calls `getMediaUrl("/api/media/{id}/file")`.
2. `getMediaUrl()` appends the user's active JWT token (`/api/media/{id}/file?token=<jwt>`).
3. Backend middleware inspects `r.URL.Query().Get("token")`, verifies signature, and streams the media stream.

---

## 🔒 Private Object Storage & Avatar Bucket Hierarchy

Media and profile photos are organized in MinIO / Local Disk Storage using 3-level UUIDv4 subfolder partitioning to prevent file system directory saturation:

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
└── avatars/                  <-- Dedicated User Profile Avatar Storage Bucket
    └── {user_id}/{filename}
```

---

## 🗺️ Page-to-API Endpoint Mapping Matrix

The following table documents exactly which backend REST API endpoint is invoked by each frontend view:

| Frontend Page / Component | Route | Backend API Endpoint(s) Called | Description / Purpose |
|---|---|---|---|
| **Gallery Page** | `/gallery` | `GET /api/media?type=gallery&deleted=false`<br>`GET /api/locked-folders`<br>`PUT /api/media/{id}`<br>`POST /api/media/batch-delete`<br>`POST /api/memories`<br>`POST /api/memories/{id}/photos` | Displays photo/video grid timeline, manages selection mode, batch moves items to bin, and assigns media to memories. |
| **Documents Page** | `/documents` | `GET /api/media?type=document&deleted=false`<br>`DELETE /api/media/{id}` | Fetches PDF and document files, handles document download and deletion. |
| **Memories Page** | `/memories` | `GET /api/memories`<br>`POST /api/memories`<br>`GET /api/memories/{id}`<br>`DELETE /api/memories/{id}`<br>`DELETE /api/memories/{id}/photos/{photoId}`<br>`GET /api/media?type=gallery&deleted=false` | Displays custom memory collections, creates new memories, views memory detail modal, and fetches automated highlights. |
| **Albums Page** | `/albums` | `GET /api/albums`<br>`POST /api/albums`<br>`GET /api/albums/{id}`<br>`PUT /api/albums/{id}`<br>`DELETE /api/albums/{id}`<br>`POST /api/albums/{id}/photos`<br>`GET /api/media?type=gallery&deleted=false` | Manages album list, creates new albums, fetches album photo lists, and updates album metadata. |
| **Locked Vault Page** | `/locked` | `GET /api/locked-folders`<br>`POST /api/locked-folders`<br>`POST /api/locked-folders/{id}/verify`<br>`DELETE /api/locked-folders/{id}`<br>`GET /api/media?locked_folder_id={id}` | Manages passcode-protected locked folders, verifies 4-digit passcodes, and renders vault media items. |
| **Bin / Trash Page** | `/bin` | `GET /api/media?deleted=true`<br>`POST /api/media/batch-restore`<br>`POST /api/media/batch-delete` | Renders soft-deleted items in bin, restores items back to timeline, and permanently purges server storage. |
| **Profile & Admin Page** | `/profile` | `GET /api/users`<br>`POST /api/users`<br>`PUT /api/users/{id}`<br>`PUT /api/users/{id}/password`<br>`PUT /api/users/{id}/role`<br>`DELETE /api/users/{id}`<br>`POST /api/users/{id}/avatar`<br>`GET /api/users/{id}/avatar`<br>`GET /api/audit-logs` | Updates user profile metadata, uploads profile avatars into `avatars/` bucket, manages accounts (Admin), and views security logs. |
| **TopHeader Component** | Global Header | `POST /api/media/upload`<br>`POST /api/media/upload-url`<br>`GET /api/users/{id}/avatar` | Handles direct multipart file uploads, ingests media from web image URLs, and renders user profile avatar. |
| **MediaViewer Component** | Global Lightbox | `GET /api/media/{id}/file`<br>`GET /api/media/{id}/thumbnail`<br>`PUT /api/media/{id}`<br>`DELETE /api/media/{id}` | Streams full-res original media (`?token=...`), toggles favorite status, moves item to locked vault, and deletes item. |

---

## 🌐 Complete REST API Specification

### Health & Authentication
* `GET /api/health` — Returns system status, SQLite connection mode, and MinIO connectivity.
* `POST /api/auth/login` — Authenticates user credentials (`{"email":"...", "password":"..."}`) and returns JWT token.

### Media Management (`/api/media`)
* `GET /api/media` — List media items (`?type=gallery`, `?type=photos`, `?type=video`, `?type=document`, `?deleted=false`, `?search=...`).
* `POST /api/media` & `POST /api/media/upload` — Multipart file upload endpoint.
* `POST /api/media/upload-url` — External web image/video ingestion (`{"url": "https://..."}`).
* `GET /api/media/{id}` — Fetch single media metadata and EXIF details.
* `PUT /api/media/{id}` — Update media title, favorite status, soft-delete state, or locked folder ID.
* `DELETE /api/media/{id}` — Permanently delete single media file.
* `POST /api/media/batch-delete` — Batch soft-delete or permanently purge items (`{"ids": ["..."]}`).
* `POST /api/media/batch-restore` — Batch restore soft-deleted items (`{"ids": ["..."]}`).
* `GET /api/media/{id}/file` — Stream original uncompressed file from storage (`?token=<jwt>`).
* `GET /api/media/{id}/thumbnail` — Stream generated WebP thumbnail image.

### Albums & Custom Memories
* `GET /api/albums` & `POST /api/albums` — List and create custom photo albums.
* `GET /api/albums/{id}` & `DELETE /api/albums/{id}` — Get detail or delete album.
* `POST /api/albums/{id}/photos` — Add photo IDs to an album.
* `GET /api/memories` & `POST /api/memories` — List and create custom memory collections.
* `POST /api/memories/{id}/photos` & `DELETE /api/memories/{id}/photos/{photoId}` — Add or remove photos in a memory.

### Passcode Locked Vault & User Management
* `GET /api/locked-folders` & `POST /api/locked-folders` — Vault folder management.
* `POST /api/locked-folders/{id}/verify` — Validate 4-digit passcode (`{"passcode": "1234"}`).
* `POST /api/users/{id}/avatar` — Upload user avatar to `avatars/{user_id}/{filename}` in MinIO.
* `GET /api/users/{id}/avatar` — Stream user avatar image.
* `GET /api/audit-logs` — Fetch authentication security audit logs *(Admin only)*.

---

## ✨ Key Features & Capabilities

### 📸 Reusable Media Viewer Lightbox (`MediaViewer.svelte`)
* **Universal Theater Lightbox**: Full-screen backdrop blur modal supporting both high-resolution photos and video streaming.
* **Centered Glassmorphism Play Overlay**: Large centered Play button overlay over video canvases with `Spacebar` play/pause shortcut.
* **Zoom & EXIF Details Drawer**: Zoom controls (+25%, -25%, 100%) and slide-over info drawer (camera model, dimension, file size).
* **Resized to Fit Badge**: Badge indicator with hover tooltip (`"Image is resized because of the screen size"`).

### 🎯 Multi-Select & Custom Memories
* **Selection Mode**: Multi-select toggle with visual checkmark cards.
* **Floating Action Toolbar**: Action bar at bottom for **Add to Memory**, **Move to Bin**, **Select All**, and **Deselect**.
* **Custom Collections**: Create named memory collections with optional descriptions.

### 🌓 Responsive Light/Dark Theme Engine
* **Dynamic Theme**: Full dynamic Tailwind `dark:` variant adaptation across all views, sidebars, modal backdrops, inputs, and toast notifications.

---

## ⚡ Getting Started & Development Setup

### Prerequisites
- **Go 1.22+**
- **Node.js 18+** & `npm`
- **SQLite 3**

### 1. Automated One-Command Launch
```bash
# Mac / Linux
./scripts/mac-linux/dev.sh

# Windows
scripts\windows\dev.bat
```

### 2. Manual Terminal Launch
```bash
# Terminal 1: Start Go Backend API (http://localhost:8080)
cd backend
go run ./cmd/server

# Terminal 2: Start SvelteKit Web App (http://localhost:5173)
cd frontend
npm run dev
```

---

## 🚀 Production Deployment Guide

```bash
# Build Backend Binary
cd backend
go build -o server ./cmd/server

# Build Frontend Production Assets
cd ../frontend
npm run build

# Start Production Server
./backend/server
```

---

## 📜 License

This project is open-source under the [Apache License 2.0](LICENSE).
