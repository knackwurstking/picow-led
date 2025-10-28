package services

import (
	"encoding/json"
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
		setup TEXT NOT NULL,
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
		return nil, err
	}

	return group, nil
}

func (g *Groups) List() ([]*models.Group, error) {
	slog.Debug("List groups from database", "table", "groups")

	query := `SELECT * FROM groups`
	rows, err := g.registry.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]*models.Group, 0)
	for rows.Next() {
		group, err := ScanGroup(rows)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (g *Groups) Add(group *models.Group) (models.GroupID, error) {
	slog.Debug("Add group to database", "table", "groups", "group", group)

	setup, err := json.Marshal(group.Setup)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO groups (name, setup) VALUES (?, ?)`
	result, err := g.registry.db.Exec(query, group.Name, setup)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return models.GroupID(id), nil
}

func (g *Groups) Update(group *models.Group) error {
	slog.Debug("Update group in database", "table", "groups", "group", group)

	setup, err := json.Marshal(group.Setup)
	if err != nil {
		return err
	}

	query := `UPDATE groups SET name = ?, setup = ? WHERE id = ?`
	_, err = g.registry.db.Exec(query, group.Name, setup, group.ID)
	return err
}

func (g *Groups) Delete(id models.GroupID) error {
	slog.Debug("Delete group from database", "table", "groups", "id", id)

	query := `DELETE FROM groups WHERE id = ?`
	_, err := g.registry.db.Exec(query, id)
	return err
}

func ScanGroup(scannable Scannable) (*models.Group, error) {
	group := &models.Group{}
	var setup string
	err := scannable.Scan(&group.ID, &group.Name, &setup, &group.CreatedAt)
	err = json.Unmarshal([]byte(setup), &group.Setup)
	return group, err
}
