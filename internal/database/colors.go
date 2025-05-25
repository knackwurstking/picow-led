package database

import (
	"database/sql"
	"log/slog"
	"slices"
)

type Colors struct {
	db *sql.DB
}

func NewColors(db *sql.DB) (*Colors, error) {
	r, err := db.Query(`SELECT name FROM sqlite_master WHERE type = "table" AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		return nil, err
	}

	tables := []string{}
	for r.Next() {
		var name string
		err := r.Scan(&name)
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

func (c *Colors) Close() {
	c.db.Close()
}
