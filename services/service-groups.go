package services

import (
	"fmt"
	"log/slog"

	"github.com/knackwurstking/picow-led/models"
)

type Groups struct {
	registry *Registry
}

func NewGroups(registry *Registry) *Groups {
	return &Groups{
		registry: registry,
	}
}

func (g *Groups) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS groups (
		id INTEGER PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := g.registry.db.Exec(query)
	if err != nil {
		return NewServiceError("create groups table", err)
	}

	return nil
}

func (g *Groups) Get(id models.GroupID) (*models.Group, error) {
	slog.Debug("Get group with ID", "id", id)

	query := `SELECT * FROM groups WHERE id = ?`
	group, err := ScanGroup(g.registry.db.QueryRow(query, id))
	if err != nil {
		return nil, NewServiceError("get group by ID", HandleSqlError(err))
	}

	return group, nil
}

func (g *Groups) List() ([]*models.Group, error) {
	slog.Debug("Get all groups")

	query := `SELECT * FROM groups`
	rows, err := g.registry.db.Query(query)
	if err != nil {
		return nil, NewServiceError("list groups", HandleSqlError(err))
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Warn("failed to close groups rows", "error", err)
		}
	}()

	var groups []*models.Group
	for rows.Next() {
		group, err := ScanGroup(rows)
		if err != nil {
			return nil, NewServiceError("scan group from rows", HandleSqlError(err))
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, NewServiceError("iterate group rows", err)
	}

	return groups, nil
}

func (g *Groups) Add(group *models.Group) (models.GroupID, error) {
	if !group.Validate() {
		return 0, fmt.Errorf("%w: %v", ErrInvalidGroup, "group validation failed")
	}

	slog.Debug("Adding a new group", "name", group.Name)

	query := `INSERT INTO groups (name) VALUES (?)`
	result, err := g.registry.db.Exec(query, group.Name)
	if err != nil {
		return 0, NewServiceError("add group", HandleSqlError(err))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, NewServiceError("get last inserted group ID", HandleSqlError(err))
	}

	return models.GroupID(id), nil
}

func (g *Groups) Update(group *models.Group) error {
	slog.Debug("Updating group", "id", group.ID, "name", group.Name)

	if !group.Validate() {
		return fmt.Errorf("%w: %v", ErrInvalidGroup, "group validation failed")
	}

	query := `UPDATE groups SET name = ? WHERE id = ?`
	_, err := g.registry.db.Exec(query, group.Name, group.ID)
	if err != nil {
		return NewServiceError("update group", HandleSqlError(err))
	}

	return nil
}

func (g *Groups) Delete(id models.GroupID) error {
	slog.Debug("Deleting group", "id", id)

	query := `DELETE FROM groups WHERE id = ?`
	_, err := g.registry.db.Exec(query, id)
	if err != nil {
		return NewServiceError("delete group", HandleSqlError(err))
	}

	return nil
}

func ScanGroup(scannable Scannable) (*models.Group, error) {
	group := &models.Group{}
	err := scannable.Scan(&group.ID, &group.Name, &group.CreatedAt)
	if err != nil {
		return nil, NewServiceError("scan group", err)
	}

	return group, nil
}
