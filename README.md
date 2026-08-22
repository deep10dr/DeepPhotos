## DeepPhotos

A lightweight, fast, self-hosted photo management application inspired by Google Photos.

DeepPhotos is designed to provide a fast and simple way to upload, organize, browse, and manage photos while keeping the application lightweight enough to run on a personal server, home lab, or low-resource machine.

---

## 🔒 Private Hierarchical MinIO Object Key Layout

Internal objects are securely stored in MinIO using a 3-level **UUIDv4** folder hierarchy:

```text
deepphotos/
├── image/
│   └── {uuid4_1}/
│       └── {uuid4_2}/
│           └── {uuid4_3}/
│               └── photo.jpg
│
├── video/
│   └── {uuid4_1}/
│       └── {uuid4_2}/
│           └── {uuid4_3}/
│               └── video.mp4
│
├── document/
│   └── {uuid4_1}/
│       └── {uuid4_2}/
│           └── {uuid4_3}/
│               └── document.pdf
│
└── lockedfolder/
    └── {uuid4_1}/
        └── {uuid4_2}/
            └── {uuid4_3}/
                └── secret.png
```

---

## ✨ Features & Capabilities

### 📸 Photo Management & Theater Viewer
* **Timeline Grid**: View photos chronologically in a responsive, clean grid layout with expressive emoji empty states.
* **Premium Theater Mode Lightbox**: High-resolution, edge-to-edge viewer with backdrop blur. Includes high z-index (`z-[9999]`) layering, automatic aspect-ratio fit to prevent image cut-off, and a **Resized to Fit Screen** indicator.
* **Reusable `<DropZone>` Component**: Drag & drop local files or web image URLs anywhere onto the page to upload them to MinIO storage.
* **Reusable `<UploadButton>` Component**: Multi-variant upload component (`primary`, `secondary`, `iconOnly`) for toolbars and page headers.
* **Favorites & Filtering**: Quickly filter media by favorites, videos, or all media.

### 📁 Albums, Vault & Organization
* **Albums**: Create and manage custom photo collections with real-time SQLite persistence. Images inside albums open in the full-screen Lightbox viewer!
* **Locked Vault & Passcode Folders**: Passcode creation workflow to generate secure locked folders.
* **Move to Vault**: Move any photo directly into a locked folder from the Lightbox viewer or grid card overlay.
* **Admin Oversight for Locked Folders**: Administrators can view locked folder metadata (creator, timestamp, count) and manage them.
* **Documents & Scans**: Dedicated storage for documents, scans, and PDFs.
* **Bin / Trash**: Soft-delete recovery and permanent trash purging.

### 👤 User Management & Security Audit
* **Admin User Control**: Add new user accounts, update user roles (*Administrator*, *Editor*, *Viewer*), or remove accounts.
* **Admin Password Reset**: Reset any user's password directly from the User Management panel.
* **Login History Audit Logs**: Detailed audit log tracking authentication attempts, IP addresses, client devices, and timestamps.
* **JWT Authentication**: Full Bearer token authentication with query parameter (`?token=...`) support for browser media tags.

---

## ⚡ Development & Debugging

### 🐛 Delve (`dlv`) Debugger for Go Backend
The Go REST API backend supports step-by-step debugging via [Delve](https://github.com/go-delve/delve):

```bash
# Mac / Linux — Start dev environment with Delve debugger on port 4000
./scripts/mac-linux/dev.sh --debug

# Windows — Start dev environment with Delve debugger on port 4000
scripts\windows\dev.bat --debug
```

#### VS Code & Cursor IDE 1-Click Debugging (`.vscode/launch.json`)
- **`🐛 Debug Go Backend (Direct)`**: Press `F5` in VS Code / Cursor to build and debug the Go backend directly with breakpoints.
- **`🔌 Attach to Remote Delve Debugger (Port 4000)`**: 1-click attach to the running Delve debug server on port 4000.

#### Docker Development Debugging
`backend/Dockerfile.dev` provides a Docker environment with Delve pre-installed listening on headless port `4000`.

---

## 🚀 Quick Start & How to Run

DeepPhotos provides automated scripts for both Mac/Linux and Windows to easily build and run the application.

### Mac & Linux
All Mac and Linux shell scripts are located in the `scripts/mac-linux/` folder.

- **Start Dev Environment**: `./scripts/mac-linux/dev.sh` *(Add `--debug` for Delve debugger)*
- **Start Container Services**: `./scripts/mac-linux/docker_run.sh` *(Commands: `up`, `stop`, `restart`)*
- **Build Docker Images**: `./scripts/mac-linux/docker_build.sh`
- **Build Local Application**: `./scripts/mac-linux/build.sh`

### Windows
All Windows batch scripts are located in the `scripts/windows/` folder.

- **Start Dev Environment**: `scripts\windows\dev.bat` *(Add `--debug` for Delve debugger)*
- **Start Container Services**: `scripts\windows\docker_run.bat` *(Commands: `up`, `stop`, `restart`)*
- **Build Docker Images**: `scripts\windows\docker_build.bat`
- **Build Local Application**: `scripts\windows\build.bat`

---

## 🛠️ Tech Stack & Optimizations

| Layer | Technology | Description |
|---|---|---|
| **Frontend** | [Svelte 5](https://svelte.dev/) / [SvelteKit](https://svelte.dev/docs/kit) | Modern, fast reactive UI framework using Svelte Runes |
| **Styling & Tooling** | TypeScript, Vite, Tailwind CSS | Type-safe front-end tooling with Light & Dark mode support |
| **Backend** | [Go (Golang)](https://golang.org/) | High-concurrency REST API microservice |
| **Backend Debugger** | [Delve (dlv)](https://github.com/go-delve/delve) | Headless debugger listening on port 4000 |
| **Database** | [SQLite](https://www.sqlite.org/) | Embedded zero-config metadata store |
| **Object Storage** | [MinIO](https://min.io/) | S3-compatible self-hosted object storage |
| **Memory Optimization** | Multi-stage Docker (`npm ci --omit=dev`) | Ultra-low RAM footprint (~120MB in production) |
| **API Caching** | Central `apiFetch` with TTL | Request deduplication & automatic cache invalidation |

---

## 📜 License

This project is licensed under the [Apache License 2.0](LICENSE).
