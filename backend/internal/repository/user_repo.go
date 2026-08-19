package repository

import (
	"database/sql"
	"errors"
	"time"

	"deepphotos/backend/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	query := `
	INSERT INTO users (id, name, email, password, role, avatar, status, last_login, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	_, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.Password, user.Role, user.Avatar, user.Status, user.LastLogin, now, now)
	return err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	query := `SELECT id, name, email, password, role, avatar, status, last_login, created_at, updated_at FROM users WHERE email = ?`
	row := r.db.QueryRow(query, email)

	var u model.User
	var lastLogin sql.NullString
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.Avatar, &u.Status, &lastLogin, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if lastLogin.Valid {
		u.LastLogin = lastLogin.String
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	query := `SELECT id, name, email, password, role, avatar, status, last_login, created_at, updated_at FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var u model.User
	var lastLogin sql.NullString
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.Avatar, &u.Status, &lastLogin, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if lastLogin.Valid {
		u.LastLogin = lastLogin.String
	}
	return &u, nil
}

func (r *UserRepository) ListAll() ([]model.User, error) {
	query := `SELECT id, name, email, role, avatar, status, last_login, created_at, updated_at FROM users ORDER BY created_at ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var lastLogin sql.NullString
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Avatar, &u.Status, &lastLogin, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		if lastLogin.Valid {
			u.LastLogin = lastLogin.String
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Update(user *model.User) error {
	query := `
	UPDATE users SET name = ?, email = ?, role = ?, status = ?, updated_at = ? WHERE id = ?
	`
	_, err := r.db.Exec(query, user.Name, user.Email, user.Role, user.Status, time.Now(), user.ID)
	return err
}

func (r *UserRepository) UpdatePassword(id string, hashedPassword string) error {
	query := `UPDATE users SET password = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, hashedPassword, time.Now(), id)
	return err
}

func (r *UserRepository) UpdateRole(id string, role string) error {
	query := `UPDATE users SET role = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, role, time.Now(), id)
	return err
}

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *UserRepository) UpdateLastLogin(email, timestamp string) error {
	query := `UPDATE users SET last_login = ? WHERE email = ?`
	_, err := r.db.Exec(query, timestamp, email)
	return err
}
