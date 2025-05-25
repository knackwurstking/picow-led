package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Path string
}

func (db *DB) NewColors() (*Colors, error) {
	dataBase, err := sql.Open("sqlite3", db.Path)
	if err != nil {
		return nil, err
	}

	return NewColors(dataBase)
}

func NewDB(path string) *DB {
	return &DB{Path: path}
}
