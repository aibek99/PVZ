package order

import (
	"context"

	"Homework-1/pkg/api/abstract"
	"Homework-1/pkg/api/order_v1"
)

// Handlers is
type Handlers interface {
	ReceiveOrder(ctx context.Context, request *order_v1.OrderCreateRequest) (*abstract.MessageResponse, error)
	IssueOrder(ctx context.Context, request *order_v1.IssueOrderRequest) (*abstract.MessageResponse, error)
	ReturnedOrders(ctx context.Context, request *abstract.Page) (*order_v1.ReturnedListResponse, error)
	AcceptOrder(ctx context.Context, request *order_v1.RequestWithClientID) (*abstract.MessageResponse, error)
	TurnInOrder(ctx context.Context, request *order_v1.OrderIDRequest) (*abstract.MessageResponse, error)
	OrderList(ctx context.Context, request *abstract.Page) (*order_v1.OrderListResponse, error)
	UniqueClientList(ctx context.Context, request *abstract.Page) (*order_v1.UniqueClientListResponse, error)
}
