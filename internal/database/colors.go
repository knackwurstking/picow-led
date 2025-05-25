package database

import "database/sql"

type Colors struct {
	db *sql.DB
}

func NewColors(db *sql.DB) (*Colors, error) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS colors (
		);
	`)
	if err != nil {
		return nil, err
	}

	return &Colors{
		db: db,
	}, nil
}
