// Package pvz ...
//
//go:generate minimock -g -i Storage -o ./mock/storage_mock.go -n StorageMock
package pvz

import (
	"Homework-1/internal/model/cli"
)

// Storage is
type Storage interface {
	Close() error
	Save() error
	Create(value *cli.PVZ) error
	Find(key string) (cli.PVZ, error)
	ReadAll() error
}
