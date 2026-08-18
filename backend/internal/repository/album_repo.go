package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"deepphotos/backend/internal/model"
)

type AlbumRepository struct {
	db *sql.DB
}

func NewAlbumRepository(db *sql.DB) *AlbumRepository {
	return &AlbumRepository{db: db}
}

func (r *AlbumRepository) Create(album *model.Album) error {
	query := `INSERT INTO albums (id, name, description, cover_url, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now()
	_, err := r.db.Exec(query, album.ID, album.Name, album.Description, album.CoverURL, now, now)
	return err
}

func (r *AlbumRepository) ListAll() ([]model.Album, error) {
	query := `
	SELECT a.id, a.name, a.description, a.cover_url, a.created_at, a.updated_at, COUNT(ap.photo_id) as photos_count
	FROM albums a
	LEFT JOIN album_photos ap ON a.id = ap.album_id
	GROUP BY a.id
	ORDER BY a.created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []model.Album
	for rows.Next() {
		var a model.Album
		var cover sql.NullString
		var desc sql.NullString
		if err := rows.Scan(&a.ID, &a.Name, &desc, &cover, &a.CreatedAt, &a.UpdatedAt, &a.PhotosCount); err != nil {
			return nil, err
		}
		if desc.Valid {
			a.Description = desc.String
		}
		if cover.Valid {
			a.CoverURL = cover.String
		}
		albums = append(albums, a)
	}
	return albums, nil
}

func (r *AlbumRepository) FindByID(id string) (*model.Album, error) {
	query := `
	SELECT a.id, a.name, a.description, a.cover_url, a.created_at, a.updated_at, COUNT(ap.photo_id) as photos_count
	FROM albums a
	LEFT JOIN album_photos ap ON a.id = ap.album_id
	WHERE a.id = ?
	GROUP BY a.id
	`
	row := r.db.QueryRow(query, id)

	var a model.Album
	var cover sql.NullString
	var desc sql.NullString
	err := row.Scan(&a.ID, &a.Name, &desc, &cover, &a.CreatedAt, &a.UpdatedAt, &a.PhotosCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if desc.Valid {
		a.Description = desc.String
	}
	if cover.Valid {
		a.CoverURL = cover.String
	}
	return &a, nil
}

func (r *AlbumRepository) Update(album *model.Album) error {
	query := `UPDATE albums SET name = ?, description = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, album.Name, album.Description, time.Now(), album.ID)
	return err
}

func (r *AlbumRepository) Delete(id string) error {
	query := `DELETE FROM albums WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *AlbumRepository) AddPhotos(albumID string, photoIDs []string) error {
	if len(photoIDs) == 0 {
		return nil
	}
	stmt, err := r.db.Prepare(`INSERT OR IGNORE INTO album_photos (album_id, photo_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, photoID := range photoIDs {
		if _, err := stmt.Exec(albumID, photoID); err != nil {
			return err
		}
	}
	return nil
}

func (r *AlbumRepository) RemovePhotos(albumID string, photoIDs []string) error {
	if len(photoIDs) == 0 {
		return nil
	}
	placeholders := make([]string, len(photoIDs))
	args := make([]interface{}, len(photoIDs)+1)
	args[0] = albumID
	for i, id := range photoIDs {
		placeholders[i] = "?"
		args[i+1] = id
	}

	query := fmt.Sprintf(`DELETE FROM album_photos WHERE album_id = ? AND photo_id IN (%s)`, strings.Join(placeholders, ","))
	_, err := r.db.Exec(query, args...)
	return err
}
