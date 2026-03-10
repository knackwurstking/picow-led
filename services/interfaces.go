package services

type Service interface {
	CreateTable() error
}

type Scannable interface {
	Scan(dest ...any) error
}
