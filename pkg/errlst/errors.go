package errlst

import "github.com/pkg/errors"

var (
	// ErrPVZAlreadyExists is
	ErrPVZAlreadyExists = errors.New("PVZ already exists in PVZ storage")
	// ErrPVZNotFound is
	ErrPVZNotFound = errors.New("PVZ not found")
	// ErrOrderAlreadyExists is
	ErrOrderAlreadyExists = errors.New("pq: duplicate key value violates unique constraint \"orders_pkey\"")
	// ErrOrderNotFound is
	ErrOrderNotFound = errors.New("Order not found")
	// ErrClientIDNotFound is
	ErrClientIDNotFound = errors.New("Not all client ids are same")
	// ErrBoxNotFound is
	ErrBoxNotFound = errors.New("Box not found")
	// ErrBoxAlreadyExists is
	ErrBoxAlreadyExists = errors.New("pq: duplicate key value violates unique constraint \"box_name_key\"")
	// ErrInvalidBoxLimit is
	ErrInvalidBoxLimit = errors.New("Exceeding box_v1 limit")
	// ErrNotFoundCache is
	ErrNotFoundCache = errors.New("Not found cache with key")
	// ErrInMemoryCacheNil is
	ErrInMemoryCacheNil = errors.New("InMemoryCache.Nil")
)
