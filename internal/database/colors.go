package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
)

type Color struct {
	ID int   `json:"id"`
	R  uint8 `json:"r"`
	G  uint8 `json:"g"`
	B  uint8 `json:"b"`
}

type Colors struct {
	db *sql.DB
}

func NewColors(db *sql.DB) (*Colors, error) {
	r, err := db.Query(`SELECT name FROM sqlite_master WHERE type = "table" AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	tables := []string{}
	var name string
	for r.Next() {
		err = r.Scan(&name)
		if err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}

	if !slices.Contains(tables, "colors") {
		slog.Debug("Create (sqlite) database table", "name", "colors")
		_, err = db.Exec(`
      		CREATE TABLE colors (
                "id" INTEGER NOT NULL,
      		    "r"  INTEGER NOT NULL,
      		    "g"  INTEGER NOT NULL,
      		    "b"  INTEGER NOT NULL,
      		    PRIMARY KEY("id" AUTOINCREMENT)
      		);
            INSERT INTO colors (r, g, b) VALUES (255, 255, 255), (255, 0, 0), (0, 255, 0), (0, 0, 255);
    	`)
		if err != nil {
			return nil, err
		}
	}

	return &Colors{
		db: db,
	}, nil
}

func (c *Colors) List() ([]Color, error) {
	colors := []Color{}

	r, err := c.db.Query(`SELECT id, r, g, b FROM colors`)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var color Color
	for r.Next() {
		err = r.Scan(&color.ID, &color.R, &color.G, &color.B)
		if err != nil {
			return nil, err
		}
		colors = append(colors, color)
	}

	return colors, nil
}

func (c *Colors) Get(id int) (Color, error) {
	query := fmt.Sprintf("SELECT r, g, b FROM colors WHERE id=%d", id)
	r, err := c.db.Query(query)
	if err != nil {
		return Color{}, err
	}
	defer r.Close()

	r.Next()
	color := Color{}
	err = r.Scan(&color.R, &color.G, &color.B)
	return color, err
}

func (c *Colors) Set(colors ...Color) error {
	err := c.DeleteAll()
	if err != nil {
		return err
	}

	query := ""
	for _, color := range colors {
		query += fmt.Sprintf(
			"INSERT INTO colors (r, g, b) VALUES (%d, %d, %d);\n",
			color.R, color.G, color.B,
		)
	}

	_, err = c.db.Exec(query)
	return err
}

func (c *Colors) Add(colors ...Color) error {
	query := ""
	for _, color := range colors {
		query += fmt.Sprintf(
			"INSERT INTO colors (r, g, b) VALUES (%d, %d, %d);\n",
			color.R, color.G, color.B,
		)
	}

	_, err := c.db.Exec(query)
	return err
}

func (c *Colors) Update(id int, color Color) error {
	query := fmt.Sprintf(
		"INSERT OR REPLACE INTO colors (id, r, g, b) VALUES (%d, %d, %d, %d);\n",
		id, color.R, color.G, color.B,
	)
	_, err := c.db.Exec(query)
	return err
}

func (c *Colors) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM colors WHERE id=%d;", id)
	_, err := c.db.Exec(query)
	return err
}

func (c *Colors) DeleteAll() error {
	_, err := c.db.Exec(`DELETE FROM colors`)
	return err
}

func (c *Colors) Close() {
	c.db.Close()
}
