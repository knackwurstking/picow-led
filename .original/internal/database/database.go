package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Devices *Devices
	Colors  *Colors

	Path string
}

func (db *DB) NewDevices() (*Devices, error) {
	dataBase, err := sql.Open("sqlite3", db.Path)
	if err != nil {
		return nil, err
	}

	if db.Devices != nil {
		db.Devices.Close()
	}

	db.Devices, err = NewDevices(dataBase)
	return db.Devices, err
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

func (db *DB) Close() {
	if db.Colors != nil {
		db.Colors.Close()
	}
}

func NewDB(path string) *DB {
	db := &DB{Path: path}

	_, err := db.NewDevices()
	if err != nil {
		panic(err)
	}

	_, err = db.NewColors()
	if err != nil {
		panic(err)
	}

	return db
}
