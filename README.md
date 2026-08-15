# 🖼️ DeepPhotos

A lightweight, fast, self-hosted photo management application inspired by Google Photos.

DeepPhotos is designed to provide a fast and simple way to upload, organize, browse, and manage photos while keeping the application lightweight enough to run on a personal server, home lab, or low-resource machine.

---

## 🚧 Project Status

> **Current Phase: 🏗️ Foundation & Svelte 5 UI Active**
>
> DeepPhotos is actively under development. The frontend is fully built on **SvelteKit** (Svelte 5 Runes + Vite + Tailwind CSS), and the backend is being architected with **Go**, **SQLite**, and **MinIO**.

---

## ✨ Features & Capabilities

### 📸 Photo Management & Viewer
* **Timeline Grid**: View photos chronologically in a responsive, clean grid layout.
* **Next & Previous Lightbox Viewer**: High-resolution viewer with **Previous (⬅️)** and **Next (➡️)** navigation controls, image counter, and keyboard shortcuts (`←`, `→`, `Esc`).
* **Favorites & Filtering**: Quickly filter photos by favorites or all media.
* **Original Retrieval**: High-res original photo downloads.

### 📁 Albums, Vault & Organization
* **Albums**: Create and manage custom photo collections.
* **Memories**: Rediscover special highlights and travel memories.
* **Locked Vault**: Passcode-protected (AES-256 encrypted) folder for sensitive items.
* **Documents & Scans**: Dedicated document storage for scans and PDFs.
* **Bin / Trash**: Deleted media recovery and storage purging.

### 👤 User Management & Security Audit
* **Top-Right Profile Avatar**: Compact header avatar linking to user profile & settings.
* **Admin User Control**: Add new user accounts with assigned roles (*Administrator*, *Editor*, *Viewer*) or remove accounts.
* **Login History Audit Logs**: Detailed audit log tracking authentication attempts, IP addresses, client devices, and timestamps.
* **Light & Dark Theme Switcher**: Toggle seamlessly between Light Sky/Cloud mode and Minimal Charcoal Dark mode.

### 🖼️ Optimized Media Pipeline
* **Asynchronous WebP Thumbnails**: Background worker generates WebP thumbnails to keep the gallery UI fast.
* **Original Storage**: High-res originals stored safely in MinIO object storage.
* **Decoupled Storage**: Metadata stored in SQLite (`photos.db`), files stored in MinIO.

---

## 🏗️ Architecture

```text
                    ┌──────────────────────┐
                    │  SvelteKit Frontend  │
                    │      (Svelte 5)      │
                    │                      │
                    │  Gallery │  Albums   │
                    │  Viewer  │  Profile  │
                    └──────────┬───────────┘
                               │
                               │ REST API
                               ▼
                    ┌──────────────────────┐
                    │       Go Backend     │
                    │                      │
                    │ Authentication       │
                    │ Photo Management     │
                    │ User & Audit Logs    │
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
| **Frontend** | [Svelte 5](https://svelte.dev/) / [SvelteKit](https://svelte.dev/docs/kit) | Modern, fast reactive UI framework using Svelte Runes |
| **Styling & Tooling** | TypeScript, Vite, Tailwind CSS | Type-safe front-end tooling with Light & Dark mode support |
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
│   │   ├── lib/             # Reusable UI components, state store, & Sidebar/Header
│   │   │   ├── components/  # Sidebar, TopHeader, Shadcn UI primitives
│   │   │   └── state.svelte.ts # Svelte 5 global state (Theme, Auth, Users, Audit logs)
│   │   └── routes/          # SvelteKit pages
│   │       ├── +page.svelte           # Minimal right-aligned login page
│   │       └── (app)/                 # Main authenticated layout
│   │           ├── gallery/           # Photo grid with Next/Prev Lightbox
│   │           ├── memories/          # Photo memories & highlights
│   │           ├── albums/            # Photo collection albums
│   │           ├── locked/            # AES-256 passcode vault
│   │           ├── documents/         # Document scans & PDFs
│   │           ├── bin/               # Deleted photo trash
│   │           └── profile/           # User details, Admin User Management & Audit Logs
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
├── LICENSE                  # Apache License 2.0
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
                 Next / Prev Lightbox viewer
                       │
                       ▼
                 Original image
```

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

### Phase 1 — Foundation & UI (Completed / Active)
* [x] SvelteKit 5 project setup & Svelte Runes architecture
* [x] Minimal right-side login page with pre-filled dev credentials
* [x] Collapsible sidebar (Minimize / Maximize toggle)
* [x] TopHeader with compact user profile avatar & Sun/Moon theme switcher
* [x] Next & Previous Lightbox photo viewer with keyboard shortcuts
* [x] Admin User Management (Add & Delete Users)
* [x] Login Activity Audit Logs
* [x] Infrastructure & Docker architecture definition

### Phase 2 — Core Photo Management & Backend Integration
* [ ] Multi-file photo uploader with progress tracking
* [ ] Background worker for thumbnail generation (WebP)
* [ ] High-performance photo grid connected to Go REST API
* [ ] EXIF metadata parsing & display

### Phase 3 — Organization & Advanced Features
* [ ] Custom photo albums & tagging
* [ ] Timeline grouping by date/year
* [ ] SHA-256 duplicate image detection
* [ ] Shared albums & temporary access links

---

## 🎯 Design Principles

* **Lightweight**: Minimal dependencies, low memory consumption. Ideal for Raspberry Pi or home NAS setup.
* **Fast & Responsive**: Thumbnail-first rendering, asynchronous background processing, and quick navigation.
* **Self-Hosted Privacy**: Full ownership of your data without reliance on cloud providers.
* **Simple Infrastructure**: Single Go binary, embedded SQLite metadata database, and MinIO object storage.

---

## 📜 License

This project is licensed under the [Apache License 2.0](LICENSE).
