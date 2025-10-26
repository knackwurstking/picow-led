package services

type Scannable interface {
	Scan(dest ...any) error
}
