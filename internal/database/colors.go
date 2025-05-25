package database

import (
	"database/sql"
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

func (c *Colors) List() []Color {
	colors := []Color{}

	r, err := c.db.Query(`SELECT (id, r, g, b) FROM colors`)
	var color Color
	for r.Next() {
		err = r.Scan(&color.ID, &color.R, &color.G, &color.B)
		if err != nil {
			panic(err)
		}
		colors = append(colors, color)
	}

	return colors
}

func (c *Colors) Close() {
	c.db.Close()
}
