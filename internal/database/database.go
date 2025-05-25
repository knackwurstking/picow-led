package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Colors *Colors

	Path string
}

func (db *DB) NewColors() (*Colors, error) {
	dataBase, err := sql.Open("sqlite3", db.Path)
	if err != nil {
		return nil, err
	}

	if db.Colors != nil {
		db.Colors.Close()
	}

	db.Colors, err = NewColors(dataBase)
	return db.Colors, err
}

func NewDB(path string) *DB {
	return &DB{Path: path}
}
