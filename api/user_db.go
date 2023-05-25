package api

import "fmt"

func getUserByToken(token string) (*user, error) {
	q := `SELECT id, name, created_at, updated_at FROM users WHERE token = ?`
	Debugf(q)
	var u user
	if err := db.QueryRow(q, token).Scan(&u.id, &u.name, &u.createdAt, &u.updatedAt); err != nil {
		return nil, fmt.Errorf("row.Scan failed: %w", err)
	}

	return &u, nil
}
