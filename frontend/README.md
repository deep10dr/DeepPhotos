# ⚡ DeepPhotos Frontend Web Application

The frontend web application for **DeepPhotos**, built on **SvelteKit** (Svelte 5 Runes + Vite + Tailwind CSS).

---

## 🌟 Key Application Features

* **Svelte 5 Runes State Architecture**: Managed via `$state()`, `$derived()`, and global store `src/lib/state.svelte.ts`.
* **Collapsible Sidebar**: Left navigation panel with Minimize / Maximize toggle (`Gallery`, `Memories`, `Albums`, `Locked Vault`, `Documents`, `Bin`, `Settings`).
* **Top Header & Sun/Moon Theme Switcher**: Minimal header with quick User Profile avatar button and Light/Dark mode switcher.
* **Timeline Photo Grid & Next/Prev Lightbox**: Photo gallery featuring full-screen viewer with Next (`ChevronRight`) / Prev (`ChevronLeft`) controls and keyboard shortcuts (`←`, `→`, `Esc`).
* **Admin User Management**: Admin panel to add and delete user accounts.
* **Security Audit Trail**: Timestamped login logs tracking IP addresses and client devices.

---

## 🛠️ Development & Build

```bash
# Install dependencies
npm install

# Start Vite dev server
npm run dev

# Run Svelte diagnostics check
npm run check

# Build production bundle
npm run build
```
