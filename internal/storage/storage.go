package storage

import "fmt"

var (
	ErrOrderNotFound    = fmt.Errorf("not found")
	ErrAlreadyExists    = fmt.Errorf("already exists")
	ErrInternal         = fmt.Errorf("internal error")
	ErrCanNotConnect    = fmt.Errorf("can not connect")
	ErrConntctionString = fmt.Errorf("connection string invalid")
	ErrNotResponding    = fmt.Errorf("database is not responding")
)
