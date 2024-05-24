// Package order ..
//
//go:generate minimock -g -i UseCase -o ./mock/usecase_mock.go -n UseCaseMock
package order

import (
	"context"

	"Homework-1/internal/model/abstract"
	orderModel "Homework-1/internal/model/order"
)

// UseCase is
type UseCase interface {
	CreateReceiveOrder(ctx context.Context, request orderModel.Request) error
	IssueOrders(ctx context.Context, request orderModel.RequestOrderIDs) (float64, error)
	ReturnedOrders(ctx context.Context, request abstract.Page) (abstract.PaginatedResponse[orderModel.ReturnedResponse], error)
	UpdateAcceptOrder(ctx context.Context, request orderModel.RequestWithClientID) error
	DeleteReturnedOrder(ctx context.Context, orderID int64) error
	OrderList(ctx context.Context, request abstract.Page) (abstract.PaginatedResponse[orderModel.AllResponse], error)
	UniqueClientsList(ctx context.Context, request abstract.Page) (abstract.PaginatedResponse[orderModel.ListUniqueClients], error)
}
