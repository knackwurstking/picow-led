package services

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/knackwurstking/picow-led/pkg/models"
)

type GroupService struct {
	registry *Registry
}

func NewGroupService(r *Registry) *GroupService {
	return &GroupService{
		registry: r,
	}
}

func (g *GroupService) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS groups (
		id INTEGER NOT NULL,
		name TEXT NOT NULL,
		devices TEXT NOT NULL,

		PRIMARY KEY ("id" AUTOINCREMENT)
	);`

	if _, err := g.registry.db.Exec(query); err != nil {
		return NewServiceError("create groups table", err)
	}

	return nil
}

func (g *GroupService) Get(groupID models.ID) (*models.Group, error) {
	query := `SELECT * FROM groups WHERE id = ?`
	group, err := ScanGroup(g.registry.db.QueryRow(query, groupID))
	if err != nil {
		return nil, NewServiceError("get group by ID", HandleSqlError(err))
	}

	return group, nil
}

func (g *GroupService) List() ([]*models.Group, error) {
	query := `SELECT * FROM groups`
	rows, err := g.registry.db.Query(query)
	if err != nil {
		return nil, NewServiceError("list groups", HandleSqlError(err))
	}
	defer rows.Close()

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

func (g *GroupService) Add(group *models.Group) (models.ID, error) {
	if group.Validate() != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidGroup, "group validation failed")
	}

	if err := g.validateDevices(group.Devices); err != nil {
		return 0, err
	}

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

	return models.ID(id), nil
}

func (g *GroupService) Update(group *models.Group) error {
	if group.Validate() != nil {
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

func (g *GroupService) Delete(groupID models.ID) error {
	query := `DELETE FROM groups WHERE id = ?`
	_, err := g.registry.db.Exec(query, groupID)
	if err != nil {
		return NewServiceError("delete group", HandleSqlError(err))
	}

	return nil
}

func (g *GroupService) validateDevices(devices []models.ID) error {
	if slices.Contains(devices, 0) {
		return ErrInvalidGroupDeviceID
	}
	// Check database if this device exists
	for _, id := range devices {
		if _, err := g.registry.Device.Get(id); err != nil && IsNotFoundError(err) {
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
	err := scannable.Scan(&group.ID, &group.Name, &devices)
	if err != nil {
		return nil, NewServiceError("scan group", err)
	}

	err = json.Unmarshal([]byte(devices), &group.Devices)
	if err != nil {
		return nil, NewServiceError("unmarshal group devices", err)
	}

	return group, nil
}

var _ Service = (*GroupService)(nil)
