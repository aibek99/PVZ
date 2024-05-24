package delivery

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	abstractModel "Homework-1/internal/model/abstract"
	orderModel "Homework-1/internal/model/order"
	"Homework-1/internal/order"
	orderMock "Homework-1/internal/order/mock"
	"Homework-1/pkg/api/abstract"
	"Homework-1/pkg/api/order_v1"
	"Homework-1/pkg/errlst"
)

// TestOrderHandler_ReceiveOrder is
func TestOrderHandler_ReceiveOrder(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)

	tests := []*struct {
		description string
		requestBody order_v1.OrderCreateRequest
		wantResp    *abstract.MessageResponse
		wantErr     error
		useCase     order.UseCase
	}{
		{
			description: "Successfully created order",
			requestBody: order_v1.OrderCreateRequest{
				Order: &order_v1.Order{
					OrderID:  1,
					ClientID: 2,
					Weight:   9,
					BoxID:    3,
				},
				ExpireTimeDuration: 30,
			},
			wantResp: &abstract.MessageResponse{Message: "Successfully Created Order"},
			wantErr:  nil,
			useCase: orderMock.NewUseCaseMock(ctrl).CreateReceiveOrderMock.
				When(
					minimock.AnyContext,
					orderModel.Request{
						ExpireTimeDuration: 30,
						OrderID:            1,
						ClientID:           2,
						Weight:             9,
						BoxID:              3,
					}).Then(nil),
		},
		{
			description: "Order already exists",
			requestBody: order_v1.OrderCreateRequest{
				Order: &order_v1.Order{
					OrderID:  1,
					ClientID: 2,
					Weight:   9,
					BoxID:    3,
				},
				ExpireTimeDuration: 30,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.AlreadyExists, "Failed to create: pq: duplicate key value violates unique constraint \"orders_pkey\""),
			useCase: orderMock.NewUseCaseMock(ctrl).CreateReceiveOrderMock.
				When(
					minimock.AnyContext,
					orderModel.Request{
						ExpireTimeDuration: 30,
						OrderID:            1,
						ClientID:           2,
						Weight:             9,
						BoxID:              3,
					}).Then(errlst.ErrOrderAlreadyExists),
		},
		{
			description: "Internal Server error",
			requestBody: order_v1.OrderCreateRequest{
				Order: &order_v1.Order{
					OrderID:  1,
					ClientID: 2,
					Weight:   9,
					BoxID:    3,
				},
				ExpireTimeDuration: 30,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.Internal, "Failed to create: assert.AnError general error for testing"),
			useCase: orderMock.NewUseCaseMock(ctrl).CreateReceiveOrderMock.
				When(
					minimock.AnyContext,
					orderModel.Request{
						ExpireTimeDuration: 30,
						OrderID:            1,
						ClientID:           2,
						Weight:             9,
						BoxID:              3,
					}).Then(assert.AnError),
		},
		{
			description: "Request validation failed",
			requestBody: order_v1.OrderCreateRequest{
				Order: &order_v1.Order{
					OrderID:  1,
					ClientID: 2,
					Weight:   9,
					BoxID:    3,
				},
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.InvalidArgument, "reqvalidator.ReadRequest: Key: 'Request.ExpireTimeDuration' Error:Field validation for 'ExpireTimeDuration' failed on the 'required' tag"),
			useCase:  orderMock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			response, err := NewOrdersHandler(tt.useCase).ReceiveOrder(context.Background(), &tt.requestBody)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantResp, response)
		})
	}
}

// TestOrderHandler_AcceptOrder is
func TestOrderHandler_AcceptOrder(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)

	tests := []*struct {
		description string
		requestBody order_v1.RequestWithClientID
		wantResp    *abstract.MessageResponse
		wantErr     error
		useCase     order.UseCase
	}{
		{
			description: "Successfully accepted order ID with Client ID",
			requestBody: order_v1.RequestWithClientID{
				OrderID:  1,
				ClientID: 2,
			},
			wantResp: &abstract.MessageResponse{Message: "Successfully Accepted Order ID with Client ID"},
			wantErr:  nil,
			useCase: orderMock.NewUseCaseMock(ctrl).UpdateAcceptOrderMock.
				When(
					minimock.AnyContext,
					orderModel.RequestWithClientID{
						OrderID:  1,
						ClientID: 2,
					}).Then(nil),
		},
		{
			description: "Order not found",
			requestBody: order_v1.RequestWithClientID{
				OrderID:  1,
				ClientID: 2,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.NotFound, "Failed to find order: Order not found"),
			useCase: orderMock.NewUseCaseMock(ctrl).UpdateAcceptOrderMock.
				When(
					minimock.AnyContext,
					orderModel.RequestWithClientID{
						OrderID:  1,
						ClientID: 2,
					}).Then(errlst.ErrOrderNotFound),
		},
		{
			description: "Internal Server error",
			requestBody: order_v1.RequestWithClientID{
				OrderID:  1,
				ClientID: 2,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.Internal, "Failed to accept order: assert.AnError general error for testing"),
			useCase: orderMock.NewUseCaseMock(ctrl).UpdateAcceptOrderMock.
				When(
					minimock.AnyContext,
					orderModel.RequestWithClientID{
						OrderID:  1,
						ClientID: 2,
					}).Then(assert.AnError),
		},
		{
			description: "Request validation failed",
			requestBody: order_v1.RequestWithClientID{
				OrderID: 1,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.InvalidArgument, "reqvalidator.ReadRequest: Key: 'RequestWithClientID.ClientID' Error:Field validation for 'ClientID' failed on the 'required' tag"),
			useCase:  orderMock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			response, err := NewOrdersHandler(tt.useCase).AcceptOrder(context.Background(), &tt.requestBody)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantResp, response)
		})
	}
}

// TestOrderHandler_IssueOrder is
func TestOrderHandler_IssueOrder(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)

	tests := []*struct {
		description string
		requestBody order_v1.IssueOrderRequest
		wantResp    *abstract.MessageResponse
		wantErr     error
		useCase     order.UseCase
	}{
		{
			description: "Successfully accepted order ID with Client ID",
			requestBody: order_v1.IssueOrderRequest{
				OrderIDRequest: []*order_v1.OrderIDRequest{
					{OrderID: 1},
					{OrderID: 2},
				},
			},
			wantResp: &abstract.MessageResponse{Message: "Successfully Issued All OrderIDs and Total additional cost:10"},
			wantErr:  nil,
			useCase: orderMock.NewUseCaseMock(ctrl).IssueOrdersMock.
				When(
					minimock.AnyContext, orderModel.RequestOrderIDs{
						OrderIDs: []int64{1, 2},
					}).
				Then(10, nil),
		},
		{
			description: "Order not found",
			requestBody: order_v1.IssueOrderRequest{
				OrderIDRequest: []*order_v1.OrderIDRequest{
					{OrderID: 1},
					{OrderID: 2},
				},
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.NotFound, "Failed to find: Order not found"),
			useCase: orderMock.NewUseCaseMock(ctrl).IssueOrdersMock.
				When(
					minimock.AnyContext, orderModel.RequestOrderIDs{
						OrderIDs: []int64{1, 2},
					}).
				Then(10, errlst.ErrOrderNotFound),
		},
		{
			description: "Internal Server error",
			requestBody: order_v1.IssueOrderRequest{
				OrderIDRequest: []*order_v1.OrderIDRequest{
					{OrderID: 1},
					{OrderID: 2},
				},
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.Internal, "Failed to find: assert.AnError general error for testing"),
			useCase: orderMock.NewUseCaseMock(ctrl).IssueOrdersMock.
				When(
					minimock.AnyContext, orderModel.RequestOrderIDs{
						OrderIDs: []int64{1, 2},
					}).
				Then(10, assert.AnError),
		},
		{
			description: "Request validation failed",
			requestBody: order_v1.IssueOrderRequest{
				OrderIDRequest: []*order_v1.OrderIDRequest{},
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.InvalidArgument, "reqvalidator.ReadRequest: Key: 'RequestOrderIDs.OrderIDs' Error:Field validation for 'OrderIDs' failed on the 'min' tag"),
			useCase:  orderMock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			response, err := NewOrdersHandler(tt.useCase).IssueOrder(context.Background(), &tt.requestBody)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantResp, response)
		})
	}
}

// TestOrderHandler_OrderList is
func TestOrderHandler_OrderList(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)

	ctrl := minimock.NewController(t)

	tests := []*struct {
		description string
		requestBody abstract.Page
		wantResp    *order_v1.OrderListResponse
		wantErr     error
		useCase     order.UseCase
	}{
		{
			description: "Successfully got list of order",
			requestBody: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: &order_v1.OrderListResponse{
				OrderAllInfo: []*order_v1.OrderAllInfo{
					{
						Order: &order_v1.Order{
							OrderID:  1,
							ClientID: 1,
							Weight:   0,
							BoxID:    0,
						},
						CreatedAt:  timestamppb.New(fixedTime),
						UpdatedAt:  timestamppb.New(fixedTime),
						AcceptedAt: nil,
						IssuedAt:   nil,
						ExpiresAt:  timestamppb.New(fixedTime),
					},
				},
				Pagination: &abstract.Pagination{
					Page: &abstract.Page{
						CurrentPage:  1,
						ItemsPerPage: 1,
					},
					TotalItems: 1,
				},
			},
			wantErr: nil,
			useCase: orderMock.NewUseCaseMock(ctrl).OrderListMock.
				When(
					minimock.AnyContext,
					abstractModel.Page{
						CurrentPage:  1,
						ItemsPerPage: 10,
					}).
				Then(abstractModel.PaginatedResponse[orderModel.AllResponse]{
					Items: []orderModel.AllResponse{
						{
							OrderID:    1,
							ClientID:   1,
							AcceptedAt: nil,
							IssuedAt:   nil,
							ExpiresAt:  &fixedTime,
							Weight:     0,
							BoxID:      0,
							CreatedAt:  fixedTime,
							UpdatedAt:  fixedTime,
						},
					},
					CurrentPage:  1,
					ItemsPerPage: 1,
					TotalItems:   1,
				}, nil),
		},
		{
			description: "Empty order list",
			requestBody: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: &order_v1.OrderListResponse{
				OrderAllInfo: []*order_v1.OrderAllInfo{},
				Pagination: &abstract.Pagination{
					Page: &abstract.Page{
						CurrentPage:  0,
						ItemsPerPage: 0,
					},
					TotalItems: 0,
				},
			},
			wantErr: nil,
			useCase: orderMock.NewUseCaseMock(ctrl).OrderListMock.When(minimock.AnyContext, abstractModel.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			}).Then(abstractModel.PaginatedResponse[orderModel.AllResponse]{}, nil),
		},
		{
			description: "Internal server error",
			requestBody: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.Internal, "Failed to get order list: assert.AnError general error for testing"),
			useCase: orderMock.NewUseCaseMock(ctrl).OrderListMock.When(minimock.AnyContext, abstractModel.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			}).Then(abstractModel.PaginatedResponse[orderModel.AllResponse]{}, assert.AnError),
		},
		{
			description: "Request validation failed",
			requestBody: abstract.Page{
				CurrentPage: 1,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.InvalidArgument, "reqvalidator.ReadRequest: Key: 'Page.ItemsPerPage' Error:Field validation for 'ItemsPerPage' failed on the 'required' tag"),
			useCase:  orderMock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			response, err := NewOrdersHandler(tt.useCase).OrderList(context.Background(), &tt.requestBody)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantResp, response)
		})
	}
}

// TestOrderHandler_ReturnedOrders is
func TestOrderHandler_ReturnedOrders(t *testing.T) {
	t.Parallel()
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)

	ctrl := minimock.NewController(t)

	tests := []*struct {
		description string
		requestBody abstract.Page
		wantResp    *order_v1.ReturnedListResponse
		wantErr     error
		useCase     order.UseCase
	}{
		{
			description: "Successfully got list of returned orders",
			requestBody: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: &order_v1.ReturnedListResponse{ReturnedResponse: []*order_v1.ReturnedResponse{
				{
					OrderID:    1,
					ClientID:   1,
					ReturnedAt: timestamppb.New(fixedTime),
				},
			},
				Pagination: &abstract.Pagination{
					Page: &abstract.Page{
						CurrentPage:  1,
						ItemsPerPage: 1,
					},
					TotalItems: 1,
				}},
			wantErr: nil,
			useCase: orderMock.NewUseCaseMock(ctrl).ReturnedOrdersMock.When(
				minimock.AnyContext,
				abstractModel.Page{
					CurrentPage:  1,
					ItemsPerPage: 10,
				}).
				Then(abstractModel.PaginatedResponse[orderModel.ReturnedResponse]{
					Items: []orderModel.ReturnedResponse{
						{
							OrderID:    1,
							ClientID:   1,
							ReturnedAt: fixedTime,
						},
					},
					CurrentPage:  1,
					ItemsPerPage: 1,
					TotalItems:   1,
				}, nil),
		},
		{
			description: "Empty returned order list",
			requestBody: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: &order_v1.ReturnedListResponse{ReturnedResponse: []*order_v1.ReturnedResponse{},
				Pagination: &abstract.Pagination{
					Page: &abstract.Page{
						CurrentPage:  0,
						ItemsPerPage: 0,
					},
					TotalItems: 0,
				},
			},
			wantErr: nil,
			useCase: orderMock.NewUseCaseMock(ctrl).ReturnedOrdersMock.
				When(
					minimock.AnyContext,
					abstractModel.Page{
						CurrentPage:  1,
						ItemsPerPage: 10,
					}).
				Then(abstractModel.PaginatedResponse[orderModel.ReturnedResponse]{}, nil),
		},
		{
			description: "Internal server error",
			requestBody: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.Internal, "Failed to get returned list: assert.AnError general error for testing"),
			useCase: orderMock.NewUseCaseMock(ctrl).ReturnedOrdersMock.
				When(
					minimock.AnyContext,
					abstractModel.Page{
						CurrentPage:  1,
						ItemsPerPage: 10,
					}).
				Then(abstractModel.PaginatedResponse[orderModel.ReturnedResponse]{}, assert.AnError),
		},
		{
			description: "Request validation failed",
			requestBody: abstract.Page{
				CurrentPage: 1,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.InvalidArgument, "reqvalidator.ReadRequest: Key: 'Page.ItemsPerPage' Error:Field validation for 'ItemsPerPage' failed on the 'required' tag"),
			useCase:  orderMock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			response, err := NewOrdersHandler(tt.useCase).ReturnedOrders(context.Background(), &tt.requestBody)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantResp, response)
		})
	}
}

// TestOrderHandler_TurnInOrder is
func TestOrderHandler_TurnInOrder(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	var mockOrderID int64 = 1

	tests := []*struct {
		description string
		requestID   order_v1.OrderIDRequest
		wantResp    *abstract.MessageResponse
		wantErr     error
		useCase     order.UseCase
	}{
		{
			description: "Successfully turn in with given order ID",
			requestID:   order_v1.OrderIDRequest{OrderID: 1},
			wantResp:    &abstract.MessageResponse{Message: "Successfully turn in with given order ID"},
			wantErr:     nil,
			useCase:     orderMock.NewUseCaseMock(ctrl).DeleteReturnedOrderMock.When(minimock.AnyContext, mockOrderID).Then(nil),
		},
		{
			description: "Order not found",
			requestID:   order_v1.OrderIDRequest{OrderID: 1},
			wantResp:    nil,
			wantErr:     status.Errorf(codes.NotFound, "Failed to find order: Order not found"),
			useCase:     orderMock.NewUseCaseMock(ctrl).DeleteReturnedOrderMock.When(minimock.AnyContext, mockOrderID).Then(errlst.ErrOrderNotFound),
		},
		{
			description: "Internal server error",
			requestID:   order_v1.OrderIDRequest{OrderID: 1},
			wantResp:    nil,
			wantErr:     status.Errorf(codes.Internal, "Failed to turn in order: assert.AnError general error for testing"),
			useCase:     orderMock.NewUseCaseMock(ctrl).DeleteReturnedOrderMock.When(minimock.AnyContext, mockOrderID).Then(assert.AnError),
		},
		{
			description: "Unable to parse orderID",
			requestID:   order_v1.OrderIDRequest{OrderID: -1},
			wantResp:    nil,
			wantErr:     status.Errorf(codes.InvalidArgument, "invalid request id"),
			useCase:     orderMock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			response, err := NewOrdersHandler(tt.useCase).TurnInOrder(context.Background(), &tt.requestID)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantResp, response)
		})
	}
}

// TestOrderHandler_UniqueClientList is
func TestOrderHandler_UniqueClientList(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)

	tests := []*struct {
		description string
		requestBody abstract.Page
		wantResp    *order_v1.UniqueClientListResponse
		wantErr     error
		useCase     order.UseCase
	}{
		{
			description: "Successfully got list of unique clients",
			requestBody: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: &order_v1.UniqueClientListResponse{
				ClientIDs: []int64{1, 2},
				Pagination: &abstract.Pagination{
					Page: &abstract.Page{
						CurrentPage:  1,
						ItemsPerPage: 2,
					},
					TotalItems: 2,
				},
			},
			wantErr: nil,
			useCase: orderMock.NewUseCaseMock(ctrl).UniqueClientsListMock.
				When(
					minimock.AnyContext,
					abstractModel.Page{
						CurrentPage:  1,
						ItemsPerPage: 10,
					}).
				Then(
					abstractModel.PaginatedResponse[orderModel.ListUniqueClients]{
						Items: []orderModel.ListUniqueClients{
							{
								ClientID: 1,
							},
							{
								ClientID: 2,
							},
						},
						CurrentPage:  1,
						ItemsPerPage: 2,
						TotalItems:   2,
					}, nil),
		},
		{
			description: "Empty unique client list",
			requestBody: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: &order_v1.UniqueClientListResponse{ClientIDs: []int64{},
				Pagination: &abstract.Pagination{
					Page: &abstract.Page{
						CurrentPage:  0,
						ItemsPerPage: 0,
					},
					TotalItems: 0,
				},
			},
			wantErr: nil,
			useCase: orderMock.NewUseCaseMock(ctrl).UniqueClientsListMock.
				When(
					minimock.AnyContext,
					abstractModel.Page{
						CurrentPage:  1,
						ItemsPerPage: 10,
					}).
				Then(
					abstractModel.PaginatedResponse[orderModel.ListUniqueClients]{}, nil),
		},
		{
			description: "Internal server error",
			requestBody: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.Internal, "Failed to find in clients: assert.AnError general error for testing"),
			useCase: orderMock.NewUseCaseMock(ctrl).UniqueClientsListMock.
				When(
					minimock.AnyContext,
					abstractModel.Page{
						CurrentPage:  1,
						ItemsPerPage: 10,
					}).
				Then(
					abstractModel.PaginatedResponse[orderModel.ListUniqueClients]{}, assert.AnError),
		},
		{
			description: "Request validation failed",
			requestBody: abstract.Page{
				CurrentPage: 1,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.InvalidArgument, "reqvalidator.ReadRequest: Key: 'Page.ItemsPerPage' Error:Field validation for 'ItemsPerPage' failed on the 'required' tag"),
			useCase:  orderMock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			response, err := NewOrdersHandler(tt.useCase).UniqueClientList(context.Background(), &tt.requestBody)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantResp, response)
		})
	}
}
