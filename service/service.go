package service

type Scannable interface {
	Scan(dest ...any) error
}
