package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"deepphotos/backend/internal/model"
)

type MemoryRepository struct {
	db *sql.DB
}

func NewMemoryRepository(db *sql.DB) *MemoryRepository {
	return &MemoryRepository{db: db}
}

func (r *MemoryRepository) Create(memory *model.Memory) error {
	query := `INSERT INTO memories (id, title, description, cover_photo_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now()
	_, err := r.db.Exec(query, memory.ID, memory.Title, memory.Description, memory.CoverPhotoID, now, now)
	return err
}

func (r *MemoryRepository) ListAll() ([]model.Memory, error) {
	query := `
	SELECT m.id, m.title, COALESCE(m.description, ''), COALESCE(m.cover_photo_id, ''), m.created_at, m.updated_at, COUNT(mp.photo_id) as items_count
	FROM memories m
	LEFT JOIN memory_photos mp ON m.id = mp.memory_id
	GROUP BY m.id
	ORDER BY m.created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []model.Memory
	for rows.Next() {
		var m model.Memory
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.CoverPhotoID, &m.CreatedAt, &m.UpdatedAt, &m.ItemsCount); err != nil {
			return nil, err
		}
		memories = append(memories, m)
	}
	return memories, nil
}

func (r *MemoryRepository) FindByID(id string) (*model.Memory, error) {
	query := `
	SELECT m.id, m.title, COALESCE(m.description, ''), COALESCE(m.cover_photo_id, ''), m.created_at, m.updated_at, COUNT(mp.photo_id) as items_count
	FROM memories m
	LEFT JOIN memory_photos mp ON m.id = mp.memory_id
	WHERE m.id = ?
	GROUP BY m.id
	`
	row := r.db.QueryRow(query, id)

	var m model.Memory
	err := row.Scan(&m.ID, &m.Title, &m.Description, &m.CoverPhotoID, &m.CreatedAt, &m.UpdatedAt, &m.ItemsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *MemoryRepository) GetPhotosForMemory(memoryID string) ([]model.Photo, error) {
	query := `
	SELECT p.id, p.filename, p.object_key, p.thumbnail_key, p.mime_type, p.file_type, p.size,
	       p.width, p.height, COALESCE(p.exif_model, ''), COALESCE(p.taken_at, ''), p.uploaded_at,
	       p.latitude, p.longitude, COALESCE(p.hash, ''), p.is_favorite, p.is_deleted, COALESCE(p.locked_folder_id, ''), p.title
	FROM photos p
	JOIN memory_photos mp ON p.id = mp.photo_id
	WHERE mp.memory_id = ? AND p.is_deleted = 0
	ORDER BY p.uploaded_at DESC
	`
	rows, err := r.db.Query(query, memoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []model.Photo
	for rows.Next() {
		var p model.Photo
		err := rows.Scan(
			&p.ID, &p.Filename, &p.ObjectKey, &p.ThumbnailKey, &p.MimeType, &p.FileType, &p.Size,
			&p.Width, &p.Height, &p.ExifModel, &p.TakenAt, &p.UploadedAt,
			&p.Latitude, &p.Longitude, &p.Hash, &p.IsFavorite, &p.IsDeleted, &p.LockedFolderID, &p.Title,
		)
		if err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	return photos, nil
}

func (r *MemoryRepository) Delete(id string) error {
	query := `DELETE FROM memories WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *MemoryRepository) AddPhotos(memoryID string, photoIDs []string) error {
	if len(photoIDs) == 0 {
		return nil
	}
	stmt, err := r.db.Prepare(`INSERT OR IGNORE INTO memory_photos (memory_id, photo_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, photoID := range photoIDs {
		if _, err := stmt.Exec(memoryID, photoID); err != nil {
			return err
		}
	}
	return nil
}

func (r *MemoryRepository) RemovePhotos(memoryID string, photoIDs []string) error {
	if len(photoIDs) == 0 {
		return nil
	}
	placeholders := make([]string, len(photoIDs))
	args := make([]interface{}, len(photoIDs)+1)
	args[0] = memoryID
	for i, id := range photoIDs {
		placeholders[i] = "?"
		args[i+1] = id
	}

	query := fmt.Sprintf(`DELETE FROM memory_photos WHERE memory_id = ? AND photo_id IN (%s)`, strings.Join(placeholders, ","))
	_, err := r.db.Exec(query, args...)
	return err
}
