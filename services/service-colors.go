package services

import (
	"log/slog"

	"github.com/knackwurstking/picow-led/models"
)

type Colors struct {
	registry *Registry
}

func NewColors(registry *Registry) *Colors {
	return &Colors{
		registry: registry,
	}
}

func (c *Colors) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS colors (
		id INTEGER PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		duty TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := c.registry.db.Exec(query)
	return err
}

func (c *Colors) Get(id models.ColorID) (*models.Color, error) {
	slog.Debug("Get color from database", "table", "colors", "id", id)

	query := `SELECT * FROM colors WHERE id = ?`
	color, err := ScanColor(c.registry.db.QueryRow(query, id))
	if err != nil {
		return nil, err
	}

	return color, nil
}

func (c *Colors) List() ([]*models.Color, error) {
	slog.Debug("List colors from database", "table", "colors")

	query := `SELECT * FROM colors`
	rows, err := c.registry.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var colors []*models.Color
	for rows.Next() {
		color, err := ScanColor(rows)
		if err != nil {
			return nil, err
		}
		colors = append(colors, color)
	}

	return colors, nil
}

func (c *Colors) Add(color *models.Color) (models.ColorID, error) {
	slog.Debug("Add color to database", "table", "colors", "id", color.ID)

	if !color.Validate() {
		return 0, ErrInvalidColor
	}

	query := `INSERT INTO colors (name, duty) VALUES (?, ?)`
	result, err := c.registry.db.Exec(query, color.Name, color.Duty)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return models.ColorID(id), nil
}

func (c *Colors) Update(color *models.Color) error {
	slog.Debug("Update color in database", "table", "colors", "id", color.ID)

	if !color.Validate() {
		return ErrInvalidColor
	}

	query := `UPDATE colors SET name = ?, duty = ? WHERE id = ?`
	_, err := c.registry.db.Exec(query, color.Name, color.Duty, color.ID)
	return err
}

func (c *Colors) Delete(id models.ColorID) error {
	slog.Debug("Delete color from database", "table", "colors", "id", id)

	query := `DELETE FROM colors WHERE id = ?`
	_, err := c.registry.db.Exec(query, id)
	return err
}

func ScanColor(scannable Scannable) (*models.Color, error) {
	color := &models.Color{}
	err := scannable.Scan(&color.ID, &color.Name, &color.Duty, &color.CreatedAt)
	if err != nil {
		return nil, err
	}

	return color, nil
}
