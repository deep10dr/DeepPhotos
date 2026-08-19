# ⚙️ DeepPhotos Backend Microservice

A high-performance, lightweight REST API built in **Go (Golang)** with embedded **SQLite 3** (`photos.db`) metadata storage, self-hosted **MinIO S3** object storage, and JWT authentication.

---

## 🏗️ Architecture & Storage Strategy

DeepPhotos decouples **structured metadata** from **heavy binary media files**:

```text
 ┌──────────────────────┐
 │  SvelteKit Frontend  │
 └──────────┬───────────┘
            │
            │ REST API
            ▼
 ┌──────────────────────┐
 │    Go Backend API    │
 └──────┬────────┬──────┘
        │        │
  Metadata      Files
        ▼        ▼
  ┌──────────┐ ┌──────────┐
  │  SQLite  │ │  MinIO   │
  │photos.db │ │  Bucket  │
  └──────────┘ └──────────┘
```

* **SQLite (`data/photos.db`)**: Embedded zero-CGO database storing `users`, `photos`, `albums`, `album_photos`, and `login_logs`.
* **MinIO (`deepphotos` bucket)**: S3-compatible object storage storing original uploaded photos (`originals/`) and WebP thumbnails (`thumbnails/`).

---

## 🛠️ Complete REST API Documentation

### 1. System Health (`/api/health`)

#### `GET /api/health`
Checks backend microservice status, SQLite database connection, and MinIO object storage availability.

**Response `200 OK`**:
```json
{
  "status": "OK",
  "database": "SQLite 3 Connected (WAL Mode)",
  "storage": "MinIO S3 Connected",
  "timestamp": "2026-08-16T00:10:00Z"
}
```

---

### 2. Authentication API (`/api/auth`)

#### `POST /api/auth/login`
Authenticates user credentials and generates a signed JWT bearer token.

**Request Body**:
```json
{
  "email": "admin@deepphotos.local",
  "password": "deepphotos2026"
}
```

**Response `200 OK`**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "usr_admin_1",
    "name": "Deepak (Admin)",
    "email": "admin@deepphotos.local",
    "role": "Administrator",
    "avatar": "https://images.unsplash.com/...",
    "status": "Active",
    "last_login": "Just now"
  }
}
```

---

### 3. Photos API (`/api/photos`)

#### `GET /api/photos`
List timeline photos with optional query parameters.

* **Query Parameters**:
  - `favorite=true` (Filter favorites only)
  - `deleted=true` (Filter bin/trash items only)
  - `search=query` (Search by title or filename)

**Response `200 OK`**:
```json
[
  {
    "id": "img_a8f92b10",
    "filename": "sunset.jpg",
    "object_key": "originals/img_a8f92b10.jpg",
    "thumbnail_key": "thumbnails/img_a8f92b10.webp",
    "mime_type": "image/jpeg",
    "size": 4200100,
    "width": 1920,
    "height": 1080,
    "taken_at": "Recently",
    "is_favorite": true,
    "is_deleted": false,
    "title": "sunset",
    "url": "/api/photos/img_a8f92b10/file",
    "thumbnail_url": "/api/photos/img_a8f92b10/thumbnail"
  }
]
```

#### `POST /api/photos`
Upload single or multiple image files.
* **Content-Type**: `multipart/form-data`
* **Form Field**: `files` (array of file streams)

**Response `201 Created`**:
```json
[
  {
    "id": "img_b92c41a0",
    "filename": "coastal_waves.png",
    "object_key": "originals/img_b92c41a0.png",
    "thumbnail_key": "thumbnails/img_b92c41a0.webp",
    "mime_type": "image/png",
    "size": 3120400,
    "title": "coastal_waves",
    "url": "/api/photos/img_b92c41a0/file",
    "thumbnail_url": "/api/photos/img_b92c41a0/thumbnail"
  }
]
```

#### `GET /api/photos/{id}`
Fetch single photo metadata.

#### `PUT /api/photos/{id}`
Update photo title, favorite status, or move to bin.

**Request Body**:
```json
{
  "title": "Alpine Peak Sunrise Updated",
  "is_favorite": true,
  "is_deleted": false
}
```

#### `DELETE /api/photos/{id}`
Permanently delete single photo from SQLite database and MinIO storage.

**Response `200 OK`**:
```json
{
  "message": "Photo deleted successfully"
}
```

#### `POST /api/photos/batch-delete`
Permanently delete multiple photos in batch.

**Request Body**:
```json
{
  "ids": ["img_a8f92b10", "img_b92c41a0"]
}
```

#### `POST /api/photos/batch-restore`
Restore multiple deleted photos from bin back to active timeline.

**Request Body**:
```json
{
  "ids": ["img_a8f92b10"]
}
```

#### `GET /api/photos/{id}/file`
Stream original full-resolution media file from MinIO.

#### `GET /api/photos/{id}/thumbnail`
Stream WebP thumbnail file from MinIO.

---

### 4. Albums API (`/api/albums`)

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/albums` | List all albums |
| `POST` | `/api/albums` | Create a new album (`{"name": "Summer 2026", "description": "Vacation"}`) |
| `GET` | `/api/albums/{id}` | Get album detail and photos count |
| `PUT` | `/api/albums/{id}` | Update album name or description |
| `DELETE` | `/api/albums/{id}` | Delete album |
| `POST` | `/api/albums/{id}/photos` | Add photos to album (`{"photo_ids": ["img_1", "img_2"]}`) |

---

### 5. Users API (`/api/users`)

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/users` | List all registered user accounts |
| `POST` | `/api/users` | Create user account (`{"name": "Sarah", "email": "sarah@local", "role": "Editor"}`) |
| `PUT` | `/api/users/{id}` | Update user name, email, role, or status |
| `DELETE` | `/api/users/{id}` | Delete user account |

---

### 6. Audit API (`/api/audit-logs`)

#### `GET /api/audit-logs`
Get recent authentication audit logs (`?limit=50`).

---

## ⚡ Environment Configuration

| Variable | Default Value | Description |
|---|---|---|
| `PORT` | `8080` | HTTP REST API listening port |
| `DB_PATH` | `data/photos.db` | SQLite database file location |
| `MINIO_ENDPOINT` | `localhost:9000` | MinIO storage endpoint |
| `MINIO_ACCESS_KEY` | `minioadmin` | MinIO access key |
| `MINIO_SECRET_KEY` | `minioadmin` | MinIO secret key |
| `MINIO_BUCKET` | `deepphotos` | MinIO storage bucket name |
| `JWT_SECRET` | `deepphotos-secret-key-2026` | HMAC signing secret for JWT tokens |

---

## 🚀 Running Locally

```bash
# Compile and run
go run ./cmd/server
```

Or build binary:

```bash
go build -o server ./cmd/server
./server
```
