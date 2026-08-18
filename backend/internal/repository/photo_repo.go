package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"deepphotos/backend/internal/model"
)

type PhotoRepository struct {
	db *sql.DB
}

func NewPhotoRepository(db *sql.DB) *PhotoRepository {
	return &PhotoRepository{db: db}
}

func (r *PhotoRepository) Create(photo *model.Photo) error {
	query := `
	INSERT INTO photos (id, filename, object_key, thumbnail_key, mime_type, file_type, size, width, height, exif_model, taken_at, uploaded_at, latitude, longitude, hash, is_favorite, is_deleted, locked_folder_id, title)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	fav := 0
	if photo.IsFavorite {
		fav = 1
	}
	del := 0
	if photo.IsDeleted {
		del = 1
	}
	fileType := photo.FileType
	if fileType == "" {
		fileType = "image"
	}
	now := time.Now()
	_, err := r.db.Exec(query, photo.ID, photo.Filename, photo.ObjectKey, photo.ThumbnailKey, photo.MimeType, fileType, photo.Size, photo.Width, photo.Height, photo.ExifModel, photo.TakenAt, now, photo.Latitude, photo.Longitude, photo.Hash, fav, del, photo.LockedFolderID, photo.Title)
	return err
}

func (r *PhotoRepository) FindByID(id string) (*model.Photo, error) {
	query := `SELECT id, filename, object_key, thumbnail_key, mime_type, COALESCE(file_type, 'image'), size, width, height, COALESCE(exif_model, ''), taken_at, uploaded_at, latitude, longitude, hash, is_favorite, is_deleted, COALESCE(locked_folder_id, ''), title FROM photos WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var p model.Photo
	var fav, del int
	err := row.Scan(&p.ID, &p.Filename, &p.ObjectKey, &p.ThumbnailKey, &p.MimeType, &p.FileType, &p.Size, &p.Width, &p.Height, &p.ExifModel, &p.TakenAt, &p.UploadedAt, &p.Latitude, &p.Longitude, &p.Hash, &fav, &del, &p.LockedFolderID, &p.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	p.IsFavorite = fav == 1
	p.IsDeleted = del == 1
	return &p, nil
}

func (r *PhotoRepository) List(favoriteOnly bool, deletedOnly bool, searchQuery string, fileType string, lockedFolderID string) ([]model.Photo, error) {
	var conditions []string
	var args []interface{}

	if favoriteOnly {
		conditions = append(conditions, "is_favorite = 1")
	}

	if deletedOnly {
		conditions = append(conditions, "is_deleted = 1")
	} else {
		conditions = append(conditions, "is_deleted = 0")
	}

	if fileType != "" {
		conditions = append(conditions, "file_type = ?")
		args = append(args, fileType)
	}

	if lockedFolderID != "" {
		conditions = append(conditions, "locked_folder_id = ?")
		args = append(args, lockedFolderID)
	} else if !deletedOnly && fileType == "" {
		// By default in main gallery, hide items placed inside locked folders unless specifically queried
		conditions = append(conditions, "(locked_folder_id IS NULL OR locked_folder_id = '')")
	}

	if searchQuery != "" {
		conditions = append(conditions, "(title LIKE ? OR filename LIKE ?)")
		term := "%" + searchQuery + "%"
		args = append(args, term, term)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
	SELECT id, filename, object_key, thumbnail_key, mime_type, COALESCE(file_type, 'image'), size, width, height, COALESCE(exif_model, ''), taken_at, uploaded_at, latitude, longitude, hash, is_favorite, is_deleted, COALESCE(locked_folder_id, ''), title 
	FROM photos %s ORDER BY uploaded_at DESC`, whereClause)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []model.Photo
	for rows.Next() {
		var p model.Photo
		var fav, del int
		if err := rows.Scan(&p.ID, &p.Filename, &p.ObjectKey, &p.ThumbnailKey, &p.MimeType, &p.FileType, &p.Size, &p.Width, &p.Height, &p.ExifModel, &p.TakenAt, &p.UploadedAt, &p.Latitude, &p.Longitude, &p.Hash, &fav, &del, &p.LockedFolderID, &p.Title); err != nil {
			return nil, err
		}
		p.IsFavorite = fav == 1
		p.IsDeleted = del == 1
		photos = append(photos, p)
	}
	return photos, nil
}

func (r *PhotoRepository) Update(photo *model.Photo) error {
	query := `UPDATE photos SET title = ?, is_favorite = ?, is_deleted = ?, locked_folder_id = ? WHERE id = ?`
	fav := 0
	if photo.IsFavorite {
		fav = 1
	}
	del := 0
	if photo.IsDeleted {
		del = 1
	}
	_, err := r.db.Exec(query, photo.Title, fav, del, photo.LockedFolderID, photo.ID)
	return err
}

func (r *PhotoRepository) DeleteSingle(id string) (*model.Photo, error) {
	photo, err := r.FindByID(id)
	if err != nil || photo == nil {
		return nil, err
	}

	query := `DELETE FROM photos WHERE id = ?`
	_, err = r.db.Exec(query, id)
	return photo, err
}

func (r *PhotoRepository) DeleteBatch(ids []string) ([]model.Photo, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	queryFetch := fmt.Sprintf(`SELECT id, filename, object_key, thumbnail_key, mime_type, COALESCE(file_type, 'image'), size, width, height, COALESCE(exif_model, ''), taken_at, uploaded_at, latitude, longitude, hash, is_favorite, is_deleted, COALESCE(locked_folder_id, ''), title FROM photos WHERE id IN (%s)`, strings.Join(placeholders, ","))
	rows, err := r.db.Query(queryFetch, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []model.Photo
	for rows.Next() {
		var p model.Photo
		var fav, del int
		if err := rows.Scan(&p.ID, &p.Filename, &p.ObjectKey, &p.ThumbnailKey, &p.MimeType, &p.FileType, &p.Size, &p.Width, &p.Height, &p.ExifModel, &p.TakenAt, &p.UploadedAt, &p.Latitude, &p.Longitude, &p.Hash, &fav, &del, &p.LockedFolderID, &p.Title); err != nil {
			return nil, err
		}
		p.IsFavorite = fav == 1
		p.IsDeleted = del == 1
		photos = append(photos, p)
	}

	queryDelete := fmt.Sprintf(`DELETE FROM photos WHERE id IN (%s)`, strings.Join(placeholders, ","))
	_, err = r.db.Exec(queryDelete, args...)
	return photos, err
}

func (r *PhotoRepository) RestoreBatch(ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`UPDATE photos SET is_deleted = 0 WHERE id IN (%s)`, strings.Join(placeholders, ","))
	_, err := r.db.Exec(query, args...)
	return err
}
