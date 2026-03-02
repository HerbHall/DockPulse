package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// SaveCheck inserts a new image check record into the database.
func (s *Store) SaveCheck(ctx context.Context, check *ImageCheck) error {
	result, err := s.db.ExecContext(ctx,
		`INSERT INTO image_checks
			(container_name, container_id, image_ref, local_digest, remote_digest, status, checked_at, registry)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		check.ContainerName,
		check.ContainerID,
		check.ImageRef,
		check.LocalDigest,
		check.RemoteDigest,
		string(check.Status),
		check.CheckedAt,
		check.Registry,
	)
	if err != nil {
		return fmt.Errorf("insert image check: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	check.ID = id

	return nil
}

// GetLatestChecks returns the most recent check for each unique container.
func (s *Store) GetLatestChecks(ctx context.Context) ([]ImageCheck, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, container_name, container_id, image_ref, local_digest,
			remote_digest, status, checked_at, registry
		FROM image_checks
		WHERE id IN (
			SELECT MAX(id) FROM image_checks GROUP BY container_id
		)
		ORDER BY checked_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("query latest checks: %w", err)
	}
	defer func() { _ = rows.Close() }()

	return scanChecks(rows)
}

// GetCheckByContainerID returns the most recent check for the given container.
// Returns nil if no check exists for that container.
func (s *Store) GetCheckByContainerID(ctx context.Context, containerID string) (*ImageCheck, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, container_name, container_id, image_ref, local_digest,
			remote_digest, status, checked_at, registry
		FROM image_checks
		WHERE container_id = ?
		ORDER BY checked_at DESC
		LIMIT 1`, containerID)

	var c ImageCheck
	err := row.Scan(
		&c.ID, &c.ContainerName, &c.ContainerID, &c.ImageRef,
		&c.LocalDigest, &c.RemoteDigest, &c.Status, &c.CheckedAt, &c.Registry,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan check by container id: %w", err)
	}

	return &c, nil
}

// SetPreference stores a key/value preference, replacing any existing value.
func (s *Store) SetPreference(ctx context.Context, key, value string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT OR REPLACE INTO preferences (key, value) VALUES (?, ?)`,
		key, value)
	if err != nil {
		return fmt.Errorf("set preference %q: %w", key, err)
	}
	return nil
}

// GetPreference retrieves the value for the given preference key.
// Returns an empty string if the key does not exist.
func (s *Store) GetPreference(ctx context.Context, key string) (string, error) {
	var value string
	err := s.db.QueryRowContext(ctx,
		`SELECT value FROM preferences WHERE key = ?`, key).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get preference %q: %w", key, err)
	}
	return value, nil
}

func scanChecks(rows *sql.Rows) ([]ImageCheck, error) {
	var checks []ImageCheck
	for rows.Next() {
		var c ImageCheck
		if err := rows.Scan(
			&c.ID, &c.ContainerName, &c.ContainerID, &c.ImageRef,
			&c.LocalDigest, &c.RemoteDigest, &c.Status, &c.CheckedAt, &c.Registry,
		); err != nil {
			return nil, fmt.Errorf("scan image check row: %w", err)
		}
		checks = append(checks, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate image check rows: %w", err)
	}
	return checks, nil
}
