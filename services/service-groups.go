package services

import (
	"encoding/json"
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
	return err
}

func (g *Groups) Get(id models.GroupID) (*models.Group, error) {
	slog.Debug("Get group from database", "table", "groups", "id", id)

	query := `SELECT * FROM groups WHERE id = ?`
	group, err := ScanGroup(g.registry.db.QueryRow(query, id))
	if err != nil {
		return nil, HandleSqlError(err)
	}

	return group, nil
}

func (g *Groups) List() ([]*models.Group, error) {
	slog.Debug("List groups from database", "table", "groups")

	query := `SELECT * FROM groups`
	rows, err := g.registry.db.Query(query)
	if err != nil {
		return nil, HandleSqlError(err)
	}
	defer rows.Close()

	groups := make([]*models.Group, 0)
	for rows.Next() {
		group, err := ScanGroup(rows)
		if err != nil {
			return nil, HandleSqlError(err)
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (g *Groups) Add(group *models.Group) (models.GroupID, error) {
	slog.Debug("Add group to database", "table", "groups", "group", group)

	if err := g.validateDevices(group.Devices); err != nil {
		return 0, err
	}

	devices, err := json.Marshal(group.Devices)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO groups (name, devices) VALUES (?, ?)`
	result, err := g.registry.db.Exec(query, group.Name, devices)
	if err != nil {
		return 0, HandleSqlError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, HandleSqlError(err)
	}

	return models.GroupID(id), nil
}

func (g *Groups) Update(group *models.Group) error {
	slog.Debug("Update group in database", "table", "groups", "group", group)

	if err := g.validateDevices(group.Devices); err != nil {
		return err
	}

	devices, err := json.Marshal(group.Devices)
	if err != nil {
		return err
	}

	query := `UPDATE groups SET name = ?, devices = ? WHERE id = ?`
	_, err = g.registry.db.Exec(query, group.Name, devices, group.ID)
	return HandleSqlError(err)
}

func (g *Groups) Delete(id models.GroupID) error {
	slog.Debug("Delete group from database", "table", "groups", "id", id)

	query := `DELETE FROM groups WHERE id = ?`
	_, err := g.registry.db.Exec(query, id)
	return HandleSqlError(err)
}

func (g *Groups) validateDevices(devices []models.DeviceID) error {
	if slices.Contains(devices, 0) {
		return ErrInvalidDeviceID
	}
	// Check database if this device exists
	for _, device := range devices {
		if _, err := g.registry.Devices.Get(device); err != nil && err == ErrNotFound {
			return ErrInvalidDeviceID
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
	err = json.Unmarshal([]byte(devices), &group.Devices)
	return group, err
}
