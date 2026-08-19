# 📦 DeepPhotos Infrastructure & Database Setup

This directory contains the Docker Compose environment and infrastructure documentation for **DeepPhotos** self-hosted photo management system.

---

## 🗄️ MinIO Object Storage Setup

DeepPhotos uses self-hosted **MinIO** (S3-compatible object storage) for managing original media files and WebP thumbnails.

### Bucket Architecture (`deepphotos`)

```text
deepphotos/
├── originals/
│   └── YYYY/MM/
│       ├── img_a8f92b10.jpg
│       └── img_b92c41a0.png
├── thumbnails/
│   └── YYYY/MM/
│       ├── img_a8f92b10.webp
│       └── img_b92c41a0.webp
└── videos/
```

### Access Points & Default Ports
* **MinIO Storage S3 API**: `http://localhost:9000`
* **MinIO Web Console**: `http://localhost:9001`
* **Default Access Key**: `minioadmin`
* **Default Secret Key**: `minioadmin`

---

## 📊 SQLite Metadata Database (`photos.db`)

DeepPhotos uses an embedded, zero-configuration **SQLite 3** database located at `backend/data/photos.db`.

### Database Schema Tables

#### 1. `users`
| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | TEXT | PRIMARY KEY | User identifier (`usr_admin_1`) |
| `name` | TEXT | NOT NULL | User display name |
| `email` | TEXT | UNIQUE, NOT NULL | Account email address |
| `password` | TEXT | NOT NULL | Bcrypt hashed password |
| `role` | TEXT | DEFAULT 'Viewer' | Administrator, Editor, Viewer |
| `avatar` | TEXT | | Profile image URL |
| `status` | TEXT | DEFAULT 'Active' | Active or Inactive |
| `last_login` | TEXT | | Timestamp string of last active login |

#### 2. `photos`
| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | TEXT | PRIMARY KEY | Photo identifier (`img_a8f92b10`) |
| `filename` | TEXT | NOT NULL | Original upload filename |
| `object_key` | TEXT | NOT NULL | MinIO S3 object reference |
| `thumbnail_key` | TEXT | NOT NULL | MinIO WebP thumbnail reference |
| `mime_type` | TEXT | NOT NULL | `image/jpeg`, `image/png`, etc. |
| `size` | INTEGER | NOT NULL | File size in bytes |
| `width` / `height` | INTEGER | | Image pixel dimensions |
| `is_favorite` | INTEGER | DEFAULT 0 | 1 if favorited, 0 otherwise |
| `is_deleted` | INTEGER | DEFAULT 0 | 1 if moved to bin/trash, 0 otherwise |
| `title` | TEXT | | Editable photo title |

#### 3. `albums` & `album_photos`
* `albums`: Stores collection ID, name, description, cover photo URL.
* `album_photos`: Junction table linking `album_id` to `photo_id` with cascading foreign keys.

#### 4. `login_logs`
* Security audit trail storing `id`, `user_email`, `timestamp`, `ip`, `device`, and `status` (`Success` vs `Failed`).

---

## 🚀 Running Infrastructure Containers

```bash
docker compose -f docker/docker-compose.yaml up -d
```
