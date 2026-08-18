# 🖼️ DeepPhotos

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

### 📸 Photo Management & Viewer
* **Timeline Grid**: View photos chronologically in a responsive, clean grid layout with expressive emoji empty states (🖼️✨, 📄✨, 🔒✨).
* **Premium Theater Mode Lightbox**: High-resolution, edge-to-edge viewer with backdrop blur inspired by native apps (Google/Apple Photos). Includes floating controls and keyboard shortcuts (`←`, `→`, `Esc`).
* **Universal & External Drag & Drop**: Drag and drop local files, or drag images directly from external websites into the gallery to instantly upload them to MinIO!
* **Favorites & Filtering**: Quickly filter media by favorites, videos, or all media.

### 📁 Albums, Vault & Organization
* **Albums**: Create and manage custom photo collections with real-time SQLite persistence. Images inside albums open in the full-screen Lightbox viewer!
* **Locked Vault & Passcode Folders**: 2-Step passcode creation workflow to generate secure locked folders.
* **Move to Vault**: Seamlessly move any photo directly into a locked folder from the Lightbox viewer by providing the target folder's passcode.
* **Admin Oversight for Locked Folders**: Administrators can view locked folder metadata (creator, timestamp, count) and manage them.
* **Documents & Scans**: Dedicated document storage for scans and PDFs.
* **Bin / Trash**: Deleted media recovery and storage purging.

### 👤 User Management & Security Audit
* **Admin User Control**: Add new user accounts, update user roles (*Administrator*, *Editor*, *Viewer*), or remove accounts.
* **Admin Password Reset**: Reset any user's password directly from the User Management panel.
* **Login History Audit Logs**: Detailed audit log tracking authentication attempts, IP addresses, client devices, and timestamps.
* **Light & Dark Theme Switcher**: Toggle seamlessly between Light Sky/Cloud mode and Minimal Charcoal Dark mode.

---

## 🚀 Quick Start & How to Run

DeepPhotos provides automated scripts for both Mac/Linux and Windows to easily build and run the application.

### Mac & Linux
All Mac and Linux shell scripts are located in the `scripts/mac-linux/` folder.

- **Start Development Server**: `./scripts/mac-linux/dev.sh`
- **Build Local Application**: `./scripts/mac-linux/build.sh`
- **Build Docker Containers**: `./scripts/mac-linux/docker_build.sh`
- **Run the Application**: `./scripts/mac-linux/run.sh` (Starts Docker if available, otherwise runs standalone)
- **Build Local and Run**: `./scripts/mac-linux/build_and_run.sh`

### Windows
All Windows batch scripts are located in the `scripts/windows/` folder. You can double-click them in File Explorer or run them from the command prompt.

- **Start Development Server**: `scripts\windows\dev.bat`
- **Build Local Application**: `scripts\windows\build.bat`
- **Build Docker Containers**: `scripts\windows\docker_build.bat`
- **Run the Application**: `scripts\windows\run.bat`
- **Build Local and Run**: `scripts\windows\build_and_run.bat`

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

## 📜 License

This project is licensed under the [Apache License 2.0](LICENSE).
