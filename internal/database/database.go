package database

// TODO: Using sqlite3? (store: colors)
type DB struct{}

func NewDB(path string) *DB {
	return &DB{}
}
