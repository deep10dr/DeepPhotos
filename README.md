# 🖼️ DeepPhotos

A lightweight, fast, self-hosted photo management application inspired by Google Photos.

DeepPhotos is designed to provide a fast and simple way to upload, organize, browse, and manage photos while keeping the application lightweight enough to run on a personal server, home lab, or low-resource machine.

---

## 🚧 Project Status

> **Current Phase: 🏗️ Foundation & Setup**
>
> DeepPhotos is actively under development. The frontend is built on **SvelteKit** (Svelte 5 + Vite), and the backend is being architected with **Go**, **SQLite**, and **MinIO**.

---

## ✨ Features & Capabilities

### 📸 Photo Management
* **Timeline Browsing**: View photos chronologically in an interactive grid.
* **Full-Screen Viewer**: High-resolution viewer with metadata overlay.
* **Upload & Download**: Fast photo ingestion and original photo retrieval.
* **Photo Organization**: Mark favorites, delete to trash, and restore photos.

### 📁 Albums & Collections
* Create and manage custom photo albums.
* Organize photos into multiple collections.

### 🖼️ Optimized Media Pipeline
* **Asynchronous Thumbnails**: Background worker generates WebP thumbnails to keep UI fast.
* **Original Storage**: High-res originals stored safely in MinIO object storage.
* **Metadata Extraction**: EXIF data parsing (dimensions, timestamps, camera metadata, GPS).

### 🔍 Search & Filtering
* Instant search by filename, date, and EXIF tags.
* *Planned*: AI-powered semantic search and tag indexing.

### ⚡ Performance & Self-Hosting
* **Lightweight Frontend**: Fluid Svelte 5 single-page application experience.
* **Low Footprint Backend**: High-performance Go microservice with low memory usage.
* **Smart Storage**: SQLite for metadata & MinIO for media objects.

---

## 🏗️ Architecture

```text
                    ┌──────────────────────┐
                    │  SvelteKit Frontend  │
                    │      (Svelte 5)      │
                    │                      │
                    │  Photos  │  Albums   │
                    │  Search  │  Viewer   │
                    └──────────┬───────────┘
                               │
                               │ REST API
                               ▼
                    ┌──────────────────────┐
                    │       Go Backend     │
                    │                      │
                    │ Authentication       │
                    │ Photo Management     │
                    │ Album Management     │
                    │ Search & Metadata    │
                    │ Thumbnail Processing │
                    └───────┬───────┬──────┘
                            │       │
                  Metadata  │       │ Files
                            ▼       ▼
                     ┌──────────┐ ┌──────────┐
                     │  SQLite  │ │  MinIO   │
                     │          │ │  Object  │
                     │ Metadata │ │ Storage  │
                     │  Albums  │ │Originals │
                     │  Users   │ │Thumbnails│
                     └──────────┘ └──────────┘
```

---

## 🛠️ Tech Stack

| Layer | Technology | Description |
|---|---|---|
| **Frontend** | [Svelte 5](https://svelte.dev/) / [SvelteKit](https://svelte.dev/docs/kit) | Modern, fast reactive UI framework |
| **Styling & Tooling** | TypeScript, Vite, Tailwind CSS | Type-safe front-end tooling |
| **Backend** | [Go (Golang)](https://golang.org/) | High-concurrency, low-footprint REST API |
| **Database** | [SQLite](https://www.sqlite.org/) | Embedded zero-config metadata store |
| **Object Storage** | [MinIO](https://min.io/) | S3-compatible self-hosted object storage |
| **Containerization** | Docker & Docker Compose | Seamless deployment for home servers |

---

## 📁 Repository Structure

```text
DeepPhotos/
├── frontend/                # SvelteKit frontend web application
│   ├── src/
│   │   ├── lib/             # Reusable UI components, stores, & API client
│   │   └── routes/          # SvelteKit pages (Timeline, Albums, Viewer, Search)
│   ├── static/              # Static public assets
│   ├── package.json
│   ├── svelte.config.js
│   └── vite.config.ts
│
├── backend/                 # Go REST API backend (in active development)
│   ├── cmd/server/          # Application main entry point
│   ├── internal/            # API handlers, services, database & storage drivers
│   └── go.mod
│
├── docker/                  # Infrastructure & container orchestration
│   └── docker-compose.yaml  # Service setup for MinIO, Go API, & Frontend
│
├── .env.template            # Template for required environment variables
├── build.sh                 # Build helper script
├── run.sh                   # Run helper script
├── dev.sh                   # Development startup script
└── README.md                # Project documentation
```

---

## 💾 Storage & Media Design

DeepPhotos cleanly decouples **structured metadata** from **raw media objects**.

### 📊 SQLite Metadata Schema (`photos.db`)

```text
photos
├── id (UUID / PRIMARY KEY)
├── filename
├── object_key        --> (MinIO reference for original file)
├── thumbnail_key     --> (MinIO reference for WebP thumbnail)
├── mime_type
├── size
├── width, height
├── taken_at, uploaded_at
├── latitude, longitude
├── hash              --> (SHA-256 for duplicate detection)
├── is_favorite
└── is_deleted
```

### 📦 MinIO Bucket Structure

```text
photos/
├── originals/
│   └── YYYY/MM/
│       ├── a8f92.jpg
│       └── b72ac.png
├── thumbnails/
│   └── YYYY/MM/
│       ├── a8f92.webp
│       └── b72ac.webp
└── videos/
    └── YYYY/MM/
```

---

## ⚡ Thumbnail Strategy & Media Pipeline

To keep the UI responsive on low-power servers and mobile browsers, the gallery renders lightweight WebP thumbnails rather than fetching multi-megabyte original files:

```text
                User opens gallery
                       │
                       ▼
                 Thumbnail image (.webp)
                       │
                       │ Click photo
                       ▼
                 Full-size viewer
                       │
                       ▼
                 Original image
```

### 🔄 Asynchronous Upload Architecture

```text
User  ──(Upload Photo)──>  SvelteKit UI  ──(POST /api/photos)──>  Go API Service
                                                                     │
                                                 ┌───────────────────┴───────────────────┐
                                                 ▼                                       ▼
                                           MinIO Storage                           SQLite Database
                                        (Save Original File)                     (Write Photo Record)
                                                 │
                                                 ▼
                                        Background Worker
                                   (Generate WebP Thumbnail)
                                                 │
                                                 ▼
                                           MinIO Storage
                                        (Save Thumbnail File)
```

---

## 🔐 Security & Access Control

* **Decoupled Credentials**: SvelteKit never handles MinIO credentials; media is securely proxied or served via presigned URLs generated by the Go backend.
* **Controlled Access**: Go API validates session tokens and user authorization before serving media or metadata.
* **Planned Features**: Multi-user authentication, private albums, link sharing with expiration dates, and rate limiting.

---

## 🚀 Getting Started

### Prerequisites

Ensure you have the following installed on your machine:

* **Node.js**: v18+ and `npm`
* **Go**: v1.21+ (for backend development)
* **Docker & Docker Compose**: for MinIO and container deployment

Verify your setup:

```bash
node --version
npm --version
go version
docker --version
docker compose version
```

### 1. Environment Setup

Copy `.env.template` to create your local `.env` file:

```bash
cp .env.template .env
```

### 2. Frontend Development

Navigate to the frontend directory, install dependencies, and start the Vite dev server:

```bash
cd frontend
npm install
npm run dev
```

The frontend client will be available at [http://localhost:5173](http://localhost:5173).

### 3. Backend Development

Navigate to the backend directory:

```bash
cd backend
go mod download
go run ./cmd/server
```

The REST API will run at [http://localhost:8080](http://localhost:8080).

### 4. MinIO Object Storage Setup

Start local storage services via Docker Compose:

```bash
docker compose -f docker/docker-compose.yaml up -d
```

Access points:
* **MinIO API**: `http://localhost:9000`
* **MinIO Web Console**: `http://localhost:9001`

---

## 🧭 Development Roadmap

### Phase 1 — Foundation (Current)
* [x] SvelteKit 5 project setup & layout
* [x] Infrastructure & Docker architecture definition
* [ ] Go REST API foundation
* [ ] SQLite schema & migration runner
* [ ] MinIO client integration

### Phase 2 — Core Photo Management
* [ ] Multi-file photo uploader with progress tracking
* [ ] Background worker for thumbnail generation (WebP)
* [ ] High-performance photo grid with lazy loading & virtual scroll
* [ ] Lightbox viewer with full EXIF metadata display
* [ ] Photo deletion (Trash) and restoration

### Phase 3 — Organization & Search
* [ ] Custom photo albums & tagging
* [ ] Favorites collection
* [ ] Timeline grouping by date/year
* [ ] Metadata and filename search

### Phase 4 — Optimization & Advanced Features
* [ ] SHA-256 duplicate image detection
* [ ] Video upload & streaming playback
* [ ] EXIF map visualization (GPS coordinates)
* [ ] Shared albums & temporary access links

---

## 🎯 Design Principles

* **Lightweight**: Minimal dependencies, low memory consumption. Ideal for Raspberry Pi or home NAS setup.
* **Fast & Responsive**: Thumbnail-first rendering, asynchronous background processing, and quick navigation.
* **Self-Hosted Privacy**: Full ownership of your data without reliance on cloud providers.
* **Simple Infrastructure**: Single Go binary, embedded SQLite metadata database, and MinIO object storage.

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).
