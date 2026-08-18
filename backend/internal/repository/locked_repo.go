package repository

import (
	"database/sql"
	"errors"
	"time"

	"deepphotos/backend/internal/model"
)

type LockedRepository struct {
	db *sql.DB
}

func NewLockedRepository(db *sql.DB) *LockedRepository {
	return &LockedRepository{db: db}
}

func (r *LockedRepository) Create(folder *model.LockedFolder) error {
	query := `
	INSERT INTO locked_folders (id, user_id, name, description, passcode_hash, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	_, err := r.db.Exec(query, folder.ID, folder.UserID, folder.Name, folder.Description, folder.PasscodeHash, now)
	return err
}

func (r *LockedRepository) FindByID(id string) (*model.LockedFolder, error) {
	query := `
	SELECT lf.id, lf.user_id, COALESCE(u.name, 'Admin'), lf.name, COALESCE(lf.description, ''), lf.passcode_hash, lf.created_at,
	(SELECT COUNT(*) FROM photos p WHERE p.locked_folder_id = lf.id AND p.is_deleted = 0) as photos_count
	FROM locked_folders lf
	LEFT JOIN users u ON lf.user_id = u.id
	WHERE lf.id = ?
	`
	row := r.db.QueryRow(query, id)

	var f model.LockedFolder
	err := row.Scan(&f.ID, &f.UserID, &f.UserName, &f.Name, &f.Description, &f.PasscodeHash, &f.CreatedAt, &f.PhotosCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &f, nil
}

func (r *LockedRepository) ListAll() ([]model.LockedFolder, error) {
	query := `
	SELECT lf.id, lf.user_id, COALESCE(u.name, 'Admin'), lf.name, COALESCE(lf.description, ''), lf.passcode_hash, lf.created_at,
	(SELECT COUNT(*) FROM photos p WHERE p.locked_folder_id = lf.id AND p.is_deleted = 0) as photos_count
	FROM locked_folders lf
	LEFT JOIN users u ON lf.user_id = u.id
	ORDER BY lf.created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []model.LockedFolder
	for rows.Next() {
		var f model.LockedFolder
		if err := rows.Scan(&f.ID, &f.UserID, &f.UserName, &f.Name, &f.Description, &f.PasscodeHash, &f.CreatedAt, &f.PhotosCount); err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}
	return folders, nil
}

func (r *LockedRepository) Delete(id string) error {
	query := `DELETE FROM locked_folders WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
