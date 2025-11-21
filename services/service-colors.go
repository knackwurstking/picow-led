package services

import (
	"fmt"
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
	if err != nil {
		return NewServiceError("create colors table", err)
	}

	return nil
}

func (c *Colors) Get(id models.ColorID) (*models.Color, error) {
	slog.Debug("Get color with ID", "id", id)

	query := `SELECT * FROM colors WHERE id = ?`
	color, err := ScanColor(c.registry.db.QueryRow(query, id))
	if err != nil {
		return nil, NewServiceError("get color by ID", HandleSqlError(err))
	}

	return color, nil
}

func (c *Colors) List() ([]*models.Color, error) {
	slog.Debug("Get all colors")

	query := `SELECT * FROM colors`
	rows, err := c.registry.db.Query(query)
	if err != nil {
		return nil, NewServiceError("list colors", HandleSqlError(err))
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Warn("failed to close colors rows", "error", err)
		}
	}()

	var colors []*models.Color
	for rows.Next() {
		color, err := ScanColor(rows)
		if err != nil {
			return nil, NewServiceError("scan color from rows", HandleSqlError(err))
		}
		colors = append(colors, color)
	}

	if err := rows.Err(); err != nil {
		return nil, NewServiceError("iterate color rows", err)
	}

	return colors, nil
}

func (c *Colors) Add(color *models.Color) (models.ColorID, error) {
	if !color.Validate() {
		return 0, fmt.Errorf("%w: %v", ErrInvalidColor, "color validation failed")
	}

	slog.Debug("Adding a new color", "name", color.Name, "duty", color.Duty)

	query := `INSERT INTO colors (name, duty) VALUES (?, ?)`
	result, err := c.registry.db.Exec(query, color.Name, color.Duty)
	if err != nil {
		return 0, NewServiceError("add color", HandleSqlError(err))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, NewServiceError("get last inserted color ID", HandleSqlError(err))
	}

	return models.ColorID(id), nil
}

func (c *Colors) Update(color *models.Color) error {
	slog.Debug("Updating color", "id", color.ID, "name", color.Name, "duty", color.Duty)

	if !color.Validate() {
		return fmt.Errorf("%w: %v", ErrInvalidColor, "color validation failed")
	}

	query := `UPDATE colors SET name = ?, duty = ? WHERE id = ?`
	_, err := c.registry.db.Exec(query, color.Name, color.Duty, color.ID)
	if err != nil {
		return NewServiceError("update color", HandleSqlError(err))
	}

	return nil
}

func (c *Colors) Delete(id models.ColorID) error {
	slog.Debug("Deleting color", "id", id)

	query := `DELETE FROM colors WHERE id = ?`
	_, err := c.registry.db.Exec(query, id)
	if err != nil {
		return NewServiceError("delete color", HandleSqlError(err))
	}

	return nil
}

func ScanColor(scannable Scannable) (*models.Color, error) {
	color := &models.Color{}
	err := scannable.Scan(&color.ID, &color.Name, &color.Duty, &color.CreatedAt)
	if err != nil {
		return nil, NewServiceError("scan color", err)
	}

	return color, nil
}
