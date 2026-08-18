package repository

import (
	"database/sql"
	"time"

	"deepphotos/backend/internal/model"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) LogLogin(logEntry *model.LoginLog) error {
	query := `
	INSERT INTO login_logs (id, user_email, timestamp, ip, device, status, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	_, err := r.db.Exec(query, logEntry.ID, logEntry.User, logEntry.Timestamp, logEntry.IP, logEntry.Device, logEntry.Status, now)
	return err
}

func (r *AuditRepository) ListLogs(limit int) ([]model.LoginLog, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `SELECT id, user_email, timestamp, ip, device, status, created_at FROM login_logs ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.LoginLog
	for rows.Next() {
		var l model.LoginLog
		if err := rows.Scan(&l.ID, &l.User, &l.Timestamp, &l.IP, &l.Device, &l.Status, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
