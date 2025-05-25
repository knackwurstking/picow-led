package database

import "database/sql"

type Colors struct {
	db *sql.DB
}

func NewColors(db *sql.DB) (*Colors, error) {
	// Insert default (colors) data (only if not exists)
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS colors (
            "id" INTEGER NOT NULL,
		    "r"  INTEGER NOT NULL,
		    "g"  INTEGER NOT NULL,
		    "b"  INTEGER NOT NULL,
		    PRIMARY KEY("id" AUTOINCREMENT)
		);
		AS SELECT 255 AS r, 255 AS g, 255 AS b
		AS SELECT 255 AS r, 0 AS g, 0 AS b
		AS SELECT 0 AS r, 255 AS g, 0 AS b
		AS SELECT 0 AS r, 0 AS g, 255 AS b
	`)
	if err != nil {
		return nil, err
	}

	return &Colors{
		db: db,
	}, nil
}

func (c *Colors) Close() {
	c.db.Close()
}
