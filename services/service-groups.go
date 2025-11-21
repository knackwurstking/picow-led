package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"

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
		devices TEXT NOT NULL,
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
		return nil, NewServiceError("iterate group rows", HandleSqlError(err))
	}

	return groups, nil
}

func (g *Groups) Add(group *models.Group) (models.GroupID, error) {
	if !group.Validate() {
		return 0, fmt.Errorf("%w: %v", ErrInvalidGroup, "group validation failed")
	}

	if err := g.validateDevices(group.Devices); err != nil {
		return 0, err
	}

	slog.Debug("Adding a new group", "name", group.Name)

	devices, err := json.Marshal(group.Devices)
	if err != nil {
		return 0, NewServiceError("marshal group devices", err)
	}

	query := `INSERT INTO groups (name, devices) VALUES (?, ?)`
	result, err := g.registry.db.Exec(query, group.Name, devices)
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

	if err := g.validateDevices(group.Devices); err != nil {
		return err
	}

	devices, err := json.Marshal(group.Devices)
	if err != nil {
		return NewServiceError("marshal group devices", err)
	}

	query := `UPDATE groups SET name = ?, devices = ? WHERE id = ?`
	_, err = g.registry.db.Exec(query, group.Name, devices, group.ID)
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

func (g *Groups) validateDevices(devices []models.DeviceID) error {
	if slices.Contains(devices, 0) {
		return ErrInvalidGroupDeviceID
	}
	// Check database if this device exists
	for _, device := range devices {
		if _, err := g.registry.Devices.Get(device); err != nil && IsNotFoundError(err) {
			return ErrInvalidGroupDeviceID
		} else if err != nil {
			return err
		}
	}
	return nil
}

func ScanGroup(scannable Scannable) (*models.Group, error) {
	group := &models.Group{}
	var devices string
	err := scannable.Scan(&group.ID, &group.Name, &devices, &group.CreatedAt)
	if err != nil {
		return nil, NewServiceError("scan group", err)
	}

	err = json.Unmarshal([]byte(devices), &group.Devices)
	if err != nil {
		return nil, NewServiceError("unmarshal group devices", err)
	}

	return group, nil
}
