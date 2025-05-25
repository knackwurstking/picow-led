package database

import "database/sql"

type Colors struct {
	db *sql.DB
}

func NewColors(db *sql.DB) (*Colors, error) {
	// TODO: Insert default (colors) data (only if not exists)
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS colors (
            "id" INTEGER NOT NULL,
		    "r"  INTEGER NOT NULL,
		    "g"  INTEGER NOT NULL,
		    "b"  INTEGER NOT NULL,
		    PRIMARY KEY("id" AUTOINCREMENT)
		);
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
