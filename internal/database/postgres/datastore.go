package postgres

import (
	"context"
	"fmt"
	"log"
	"sync"

	"Homework-1/internal/box"
	boxRepository "Homework-1/internal/box/repository"
	"Homework-1/internal/connection"
	"Homework-1/internal/database"
	"Homework-1/internal/order"
	orderRepository "Homework-1/internal/order/repository"
	"Homework-1/internal/pvz"
	pvzRepository "Homework-1/internal/pvz/repository"
)

var _ database.Datastore = (*DataStore)(nil)

// DataStore is
type DataStore struct {
	db        connection.DB
	pvz       pvz.Repository
	pvzInit   sync.Once
	order     order.Repository
	orderInit sync.Once
	box       box.Repository
	boxInit   sync.Once
}

// PvzRepo is
func (d *DataStore) PvzRepo() pvz.Repository {
	d.pvzInit.Do(func() {
		d.pvz = pvzRepository.NewPVZPGRepository(d.db)
	})
	return d.pvz
}

// OrderRepo is
func (d *DataStore) OrderRepo() order.Repository {
	d.orderInit.Do(func() {
		d.order = orderRepository.NewOrdersPGRepository(d.db)
	})
	return d.order
}

// BoxRepo is
func (d *DataStore) BoxRepo() box.Repository {
	d.boxInit.Do(func() {
		d.box = boxRepository.NewBoxPGRepository(d.db)
	})
	return d.box
}

// NewDataStore is
func NewDataStore(db connection.DBops) database.Datastore {
	return &DataStore{
		db: db,
	}
}

// WithTransaction is
func (d *DataStore) WithTransaction(ctx context.Context, transactionFn database.Transaction) error {
	db, ok := d.db.(connection.DBops)
	if !ok {
		return fmt.Errorf("got error start of transaction")
	}
	tx, err := db.Begin(ctx, nil)
	if err != nil {
		return fmt.Errorf("db.Begin: %w", err)
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Printf("[postgres][WithTransaction] failed to rollback transaction: %v", err)
			}
		}
	}()

	transactionalDB := &DataStore{db: tx}
	if err = transactionFn(transactionalDB); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
