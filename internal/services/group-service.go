package services

import (
	"database/sql"
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
		return fmt.Errorf("create groups table: %w", err)
	}

	return nil
}

func (g *GroupService) Get(groupID models.ID) (*models.Group, error) {
	query := `SELECT * FROM groups WHERE id = ?`
	group, err := ScanGroup(g.registry.db.QueryRow(query, groupID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("get group by ID: %w", err)
	}

	return group, nil
}

func (g *GroupService) List() ([]*models.Group, error) {
	query := `SELECT * FROM groups`
	rows, err := g.registry.db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.Group{}, nil
		}
		return nil, fmt.Errorf("list groups: %w", err)
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		group, err := ScanGroup(rows)
		if err != nil {
			return nil, fmt.Errorf("scan group row: %w", err)
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate group rows: %w", err)
	}

	return groups, nil
}

func (g *GroupService) Add(group *models.Group) (models.ID, error) {
	if err := group.Validate(); err != nil {
		return 0, ErrorValidation
	}

	if err := g.validateDevices(group.Devices); err != nil {
		return 0, err
	}

	devices, err := json.Marshal(group.Devices)
	if err != nil {
		return 0, fmt.Errorf("marshal group devices: %w", err)
	}

	query := `INSERT INTO groups (name, devices) VALUES (?, ?)`
	result, err := g.registry.db.Exec(query, group.Name, devices)
	if err != nil {
		return 0, fmt.Errorf("insert group: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last inserted group ID: %w", err)
	}

	return models.ID(id), nil
}

func (g *GroupService) Update(group *models.Group) error {
	if err := group.Validate(); err != nil {
		return ErrorValidation
	}

	if err := g.validateDevices(group.Devices); err != nil {
		return err
	}

	devices, err := json.Marshal(group.Devices)
	if err != nil {
		return fmt.Errorf("marshal group devices: %w", err)
	}
	query := `UPDATE groups SET name = ?, devices = ? WHERE id = ?`
	if r, err := g.registry.db.Exec(query, group.Name, devices, group.ID); err != nil {
		if i, _ := r.RowsAffected(); i == 0 {
			return ErrorNotFound
		}
		return fmt.Errorf("update group: %w", err)
	}

	return nil
}

func (g *GroupService) Delete(groupID models.ID) error {
	query := `DELETE FROM groups WHERE id = ?`
	if r, err := g.registry.db.Exec(query, groupID); err != nil {
		if i, _ := r.RowsAffected(); i == 0 {
			return nil
		}
		return fmt.Errorf("delete group: %w", err)
	}

	return nil
}

func (g *GroupService) validateDevices(devices []models.ID) error {
	if slices.Contains(devices, 0) {
		return ErrorValidation
	}
	// Check database if this device exists
	for _, id := range devices {
		if _, err := g.registry.Device.Get(id); err != nil {
			if err == ErrorNotFound {
				return ErrorNotFound
			} else {
				return fmt.Errorf("check device existence for group: %w", err)
			}
		}
	}
	return nil
}

func ScanGroup(scannable Scannable) (*models.Group, error) {
	group := &models.Group{}
	var devices string
	err := scannable.Scan(&group.ID, &group.Name, &devices)
	if err != nil {
		return nil, fmt.Errorf("scan group: %w", err)
	}

	err = json.Unmarshal([]byte(devices), &group.Devices)
	if err != nil {
		return nil, fmt.Errorf("unmarshal group devices: %w", err)
	}

	return group, nil
}

var _ Service = (*GroupService)(nil)
