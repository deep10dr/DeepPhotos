# ⚡ DeepPhotos Frontend Web Application

The web frontend for **DeepPhotos**, built on **SvelteKit** using **Svelte 5 Runes**, **Vite**, **Tailwind CSS v4**, and **Lucide Icons**.

---

## 🧭 Table of Contents
1. [Architecture & Svelte 5 Runes](#-architecture--svelte-5-runes)
2. [Directory & Component Layout](#-directory--component-layout)
3. [Page-to-API Endpoint Mapping Matrix](#-page-to-api-endpoint-mapping-matrix)
4. [Reusable UI Components API](#-reusable-ui-components-api)
5. [State Management & Cache Invalidation](#-state-management--cache-invalidation)
6. [Theme Engine (Light/Dark Mode)](#-theme-engine-lightdark-mode)
7. [Development & Build Commands](#-development--build-commands)

---

## 🚀 Architecture & Svelte 5 Runes

The frontend uses **Svelte 5 Runes** for fine-grained reactivity:
* **`$state()`**: Reactive local and global application state declarations.
* **`$derived()`**: Pure reactive computed values (e.g. filtered photo lists, active selected photo).
* **`$effect()`**: Side-effect execution (e.g. theme switching, keyboard shortcuts).
* **`apiFetch` Wrapper (`src/lib/api.ts`)**: Central fetch handler automatically attaching JWT `Authorization` headers, caching GET responses in TTL memory, and converting URLs with auth tokens (`?token=<jwt>`).

---

## 📁 Directory & Component Layout

```text
frontend/src/
├── lib/
│   ├── api.ts                      <-- Central fetch client with JWT token & TTL caching
│   ├── state.svelte.ts             <-- Svelte 5 global application state store
│   ├── notify.svelte.ts            <-- Toast notifications & confirmDialog state store
│   └── components/
│       ├── MediaViewer.svelte      <-- Centralized full-screen photo/video Lightbox component
│       ├── NotificationToast.svelte <-- Toast notifications & modal confirmation container
│       ├── DropZone.svelte         <-- Page-level drag-and-drop file upload container
│       ├── UploadButton.svelte     <-- Variant-driven file upload button component
│       └── TopHeader.svelte        <-- Global navigation header component
└── routes/
    └── (app)/
        ├── gallery/                <-- Main gallery media timeline page
        ├── memories/               <-- Custom user memory collections & highlights page
        ├── albums/                 <-- Photo albums management page
        ├── vault/                  <-- Passcode-protected locked folders page
        ├── documents/              <-- PDF & document scanner page
        ├── bin/                    <-- Trash recovery & permanent purge page
        └── profile/                <-- Profile settings & Admin user management page
```

---

## 🗺️ Page-to-API Endpoint Mapping Matrix

The following table documents exactly which backend REST API endpoint is called by each frontend page and component:

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

## 🧩 Reusable UI Components API

### 1. `<MediaViewer>` (`src/lib/components/MediaViewer.svelte`)
Full-screen Lightbox component for high-res photo viewing and HTML5 video streaming:
* **Props**: `photos`, `selectedIndex`, `onclose`, `onselect`, `ontogglefavorite`, `ondelete`, `onmovetovault`.
* **Hotkeys**: `Escape` (close), `ArrowLeft` (previous), `ArrowRight` (next), `Spacebar` (play/pause video).
* **Badge Indicator**: Scaling indicator with tooltip (`"Image is resized because of the screen size"`).

### 2. `<NotificationToast>` (`src/lib/components/NotificationToast.svelte`)
Global toast notification manager and modal confirmation dialog handler (`confirmDialog.ask(...)`).

---

## ⚡ State Management & Cache Invalidation

Mutations (uploads, updates, deletes) automatically clear client-side TTL memory caches via `invalidateCache('/api/media')`:

```ts
import { invalidateCache } from '$lib/api';

// After successful upload / update
invalidateCache('/api/media');
appState.refreshPhotos();
```

---

## 🎨 Theme Engine (Light/Dark Mode)

Tailwind CSS dynamic `dark:` class variant engine controlled via `appState.theme`:
* Reads saved theme preference from `localStorage`.
* Dynamically toggles `dark` class on `document.documentElement`.
* Adapts all text, card backgrounds, modal borders, and lightboxes seamlessly.

---

## 🛠️ Development & Build Commands

```bash
# Install NPM packages
npm install

# Start Vite development server
npm run dev

# Run Svelte & TypeScript diagnostic checks
npm run check

# Build production bundle
npm run build

# Preview built production app
npm run preview
```
