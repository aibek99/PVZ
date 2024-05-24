package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"Homework-1/internal/connection"
	"Homework-1/internal/model/abstract"
	"Homework-1/internal/model/box"
	"Homework-1/internal/model/order"
	orderInterface "Homework-1/internal/order"
	"Homework-1/pkg/errlst"
)

const offsetLimitQuery = "OFFSET $1 LIMIT $2"

var (
	_ orderInterface.Repository = (*OrdersRepository)(nil)
)

// OrdersRepository is
type OrdersRepository struct {
	psqlDB connection.DB
}

// NewOrdersPGRepository is
func NewOrdersPGRepository(psqlDB connection.DB) *OrdersRepository {
	return &OrdersRepository{
		psqlDB: psqlDB,
	}
}

// CreateReceiveOrder is
func (o *OrdersRepository) CreateReceiveOrder(ctx context.Context, orderData order.RequestData) error {
	log.Printf("[order][repository][CreateReceiveOrder]")

	expiresAt := time.Now().AddDate(0, 0, orderData.ExpireTimeDuration)
	_, err := o.psqlDB.Execute(ctx,
		"INSERT INTO orders(order_id, client_id, expires_at, weight, box_id) VALUES ($1,$2,$3,$4,$5);",
		orderData.OrderID,
		orderData.ClientID,
		expiresAt,
		orderData.Weight,
		orderData.BoxID,
	)
	if err != nil {
		if errors.Is(err, errlst.ErrOrderAlreadyExists) {
			return errlst.ErrOrderAlreadyExists
		}

		return fmt.Errorf("tx.ExecContext: %w", err)
	}

	return nil
}

// GetClientID is
func (o *OrdersRepository) GetClientID(ctx context.Context, orderID int64) (int64, error) {
	log.Printf("[order][repository][GetClientID]")
	var clientID int64

	err := o.psqlDB.Get(ctx, &clientID, "SELECT client_id FROM orders WHERE order_id = $1 AND returned_at IS NULL AND issued_at IS NULL", orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, errlst.ErrOrderNotFound
		}

		return -1, err
	}

	return clientID, nil
}

// BoxDataByOrderID is
func (o *OrdersRepository) BoxDataByOrderID(ctx context.Context, orderID int64) (box.DataWithOrder, error) {
	log.Println("[order][repository][BoxDataByOrderID]")
	var boxData box.DataWithOrder

	err := o.psqlDB.Get(ctx, &boxData, "SELECT cost, is_check, box.weight as weight, o.weight as order_weight FROM box join orders o on box.id = o.box_id WHERE o.order_id = $1", orderID)
	if err != nil {
		return box.DataWithOrder{}, errlst.ErrBoxNotFound
	}

	return boxData, nil
}

// UpdateIssueOrder is
func (o *OrdersRepository) UpdateIssueOrder(ctx context.Context, orderID int64, clientID int64) error {
	log.Printf("[order][repository][UpdateIssueOrder]")

	result, err := o.psqlDB.Execute(ctx, "UPDATE orders SET issued_at = $1, updated_at = $2 WHERE order_id=$3 AND client_id = $4 AND returned_at IS NULL AND issued_at IS NULL", time.Now(), time.Now(), orderID, clientID)
	if err != nil {
		return fmt.Errorf("tx.ExecContext: %w", err)
	}

	var rows int64

	rows, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}

	if rows == 0 {
		return errlst.ErrClientIDNotFound
	}

	return nil
}

// CountReturnedOrders is
func (o *OrdersRepository) CountReturnedOrders(ctx context.Context) (int64, error) {
	var totalCount int64

	err := o.psqlDB.Get(
		ctx,
		&totalCount,
		"SELECT COUNT(order_id) FROM orders WHERE returned_at IS NOT NULL",
	)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

// ListReturnedOrders is
func (o *OrdersRepository) ListReturnedOrders(ctx context.Context, orderPaginationData abstract.PageData) ([]order.ReturnedData, error) {
	log.Printf("[order][repository][ListReturnedOrders]")
	offset := (orderPaginationData.CurrentPage - 1) * orderPaginationData.ItemsPerPage

	var orderReturnedListData []order.ReturnedData

	err := o.psqlDB.Select(
		ctx,
		&orderReturnedListData,
		"SELECT order_id, client_id, returned_at FROM orders WHERE returned_at IS NOT NULL "+
			offsetLimitQuery,
		offset,
		orderPaginationData.ItemsPerPage,
	)
	if err != nil {
		return []order.ReturnedData{}, err
	}

	return orderReturnedListData, nil
}

// UpdateAcceptOrder is
func (o *OrdersRepository) UpdateAcceptOrder(ctx context.Context, acceptOrderData order.RequestWithClientIDData) error {
	log.Printf("[order][repository][UpdateAcceptOrder]")

	result, err := o.psqlDB.Execute(ctx, "UPDATE orders SET accepted_at = $1, updated_at = $2 WHERE order_id=$3 AND client_id = $4 AND returned_at IS NULL AND accepted_at IS NULL", time.Now(), time.Now(), acceptOrderData.OrderID, acceptOrderData.ClientID)
	if err != nil {
		return fmt.Errorf("o.psqlDB.ExecContext: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}

	if rows == 0 {
		return errlst.ErrOrderNotFound
	}

	return nil
}

// DeleteReturnOrder is
func (o *OrdersRepository) DeleteReturnOrder(ctx context.Context, orderID int64) error {
	log.Printf("[order][repository][DeleteReturnOrder]")

	result, err := o.psqlDB.Execute(ctx, "UPDATE orders SET returned_at = $1, updated_at = $2  WHERE order_id=$3 AND returned_at IS NULL", time.Now(), time.Now(), orderID)
	if err != nil {
		return fmt.Errorf("o.psqlDB.ExecContext: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}

	if rows == 0 {
		return errlst.ErrOrderNotFound
	}

	return nil
}

// CountOrders is
func (o *OrdersRepository) CountOrders(ctx context.Context) (int64, error) {
	var totalCount int64

	err := o.psqlDB.Get(
		ctx,
		&totalCount,
		"SELECT COUNT(order_id) FROM orders WHERE returned_at IS NULL",
	)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

// ListOrders is
func (o *OrdersRepository) ListOrders(ctx context.Context, orderPaginationData abstract.PageData) ([]order.AllResponseData, error) {
	log.Printf("[order][repository][DeleteReturnOrder]")
	offset := (orderPaginationData.CurrentPage - 1) * orderPaginationData.ItemsPerPage

	var orderListData []order.AllResponseData

	err := o.psqlDB.Select(
		ctx,
		&orderListData,
		"SELECT order_id, box_id, client_id, accepted_at, issued_at, expires_at, created_at, updated_at FROM orders WHERE returned_at IS NULL ORDER BY created_at DESC "+
			"OFFSET $1 LIMIT $2",
		offset,
		orderPaginationData.ItemsPerPage,
	)
	if err != nil {
		return []order.AllResponseData{}, err
	}

	return orderListData, nil
}

// CountUniqueClients is
func (o *OrdersRepository) CountUniqueClients(ctx context.Context) (int64, error) {
	var totalCount int64

	err := o.psqlDB.Get(
		ctx,
		&totalCount,
		"SELECT COUNT(DISTINCT client_id) FROM orders",
	)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

// ListUniqueClients is
func (o *OrdersRepository) ListUniqueClients(ctx context.Context, clientPaginationData abstract.PageData) ([]order.ListUniqueClientsData, error) {
	log.Printf("[order][repository][UniqueClientsList]")
	offset := (clientPaginationData.CurrentPage - 1) * clientPaginationData.ItemsPerPage

	var clientListData []order.ListUniqueClientsData

	err := o.psqlDB.Select(
		ctx,
		&clientListData,
		"SELECT DISTINCT ON (client_id) client_id FROM orders OFFSET $1 LIMIT $2",
		offset,
		clientPaginationData.ItemsPerPage,
	)
	if err != nil {
		return []order.ListUniqueClientsData{}, err
	}

	return clientListData, nil
}
