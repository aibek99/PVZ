package order

import (
	"context"

	"Homework-1/internal/model/abstract"
	"Homework-1/internal/model/box"
	"Homework-1/internal/model/order"
)

// Repository is
type Repository interface {
	UpdateAcceptOrder(ctx context.Context, acceptOrderData order.RequestWithClientIDData) error
	CreateReceiveOrder(ctx context.Context, orderRequestData order.RequestData) error
	UpdateIssueOrder(ctx context.Context, orderID int64, clientID int64) error
	DeleteReturnOrder(ctx context.Context, orderID int64) error
	CountReturnedOrders(ctx context.Context) (int64, error)
	CountOrders(ctx context.Context) (int64, error)
	CountUniqueClients(ctx context.Context) (int64, error)
	ListOrders(ctx context.Context, orderPaginationData abstract.PageData) ([]order.AllResponseData, error)
	ListReturnedOrders(ctx context.Context, orderPaginationData abstract.PageData) ([]order.ReturnedData, error)
	ListUniqueClients(ctx context.Context, clientPaginationData abstract.PageData) ([]order.ListUniqueClientsData, error)
	BoxDataByOrderID(ctx context.Context, orderID int64) (box.DataWithOrder, error)
	GetClientID(ctx context.Context, orderID int64) (int64, error)
}
