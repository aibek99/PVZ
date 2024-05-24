package database

import (
	"context"

	"Homework-1/internal/box"
	"Homework-1/internal/order"
	"Homework-1/internal/pvz"
)

// Transaction is
type Transaction func(db Datastore) error

// Datastore is
type Datastore interface {
	WithTransaction(ctx context.Context, transaction Transaction) error
	PvzRepo() pvz.Repository
	OrderRepo() order.Repository
	BoxRepo() box.Repository
}
