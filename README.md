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
* **Timeline Grid**: View photos chronologically in a responsive, clean grid layout.
* **Next & Previous Lightbox Viewer**: High-resolution viewer with **Previous (⬅️)** and **Next (➡️)** navigation controls, image counter, and keyboard shortcuts (`←`, `→`, `Esc`).
* **Universal Drag & Drop Ingestion**: Drag and drop photos, videos, or documents anywhere on screen to upload.
* **Favorites & Filtering**: Quickly filter media by favorites, videos, or all media.

### 📁 Albums, Vault & Organization
* **Albums**: Create and manage custom photo collections with real-time SQLite persistence.
* **Locked Vault & Passcode Folders**: 2-Step passcode creation workflow (Passcode & Confirm Passcode validation) to generate secure locked folders.
* **Admin Oversight for Locked Folders**: Administrators can view locked folder metadata (creator, timestamp, count) and manage them.
* **Documents & Scans**: Dedicated document storage for scans and PDFs.
* **Bin / Trash**: Deleted media recovery and storage purging.

### 👤 User Management & Security Audit
* **Admin User Control**: Add new user accounts, update user roles (*Administrator*, *Editor*, *Viewer*), or remove accounts.
* **Admin Password Reset**: Reset any user's password directly from the User Management panel.
* **Login History Audit Logs**: Detailed audit log tracking authentication attempts, IP addresses, client devices, and timestamps.
* **Light & Dark Theme Switcher**: Toggle seamlessly between Light Sky/Cloud mode and Minimal Charcoal Dark mode.

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
