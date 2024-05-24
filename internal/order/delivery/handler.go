package delivery

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	abstractModel "Homework-1/internal/model/abstract"
	orderModel "Homework-1/internal/model/order"
	orderInterface "Homework-1/internal/order"
	"Homework-1/pkg/api/abstract"
	"Homework-1/pkg/api/order_v1"
	"Homework-1/pkg/errlst"
	"Homework-1/pkg/reqvalidator"
	"Homework-1/pkg/tracing"
)

var (
	_ orderInterface.Handlers = (*OrderHandler)(nil)
)

// OrderHandler is
type OrderHandler struct {
	useCase orderInterface.UseCase
	order_v1.UnimplementedOrderServiceServer
}

// NewOrdersHandler is
func NewOrdersHandler(useCase orderInterface.UseCase) *OrderHandler {
	return &OrderHandler{
		useCase: useCase,
	}
}

// ReceiveOrder is
func (o *OrderHandler) ReceiveOrder(ctx context.Context, request *order_v1.OrderCreateRequest) (*abstract.MessageResponse, error) {
	log.Print("[order][delivery][ReceiveOrder]")
	tracer := otel.Tracer("[order][delivery]")
	ctx, span := tracer.Start(ctx, "[ReceiveOrder]")
	defer span.End()

	orderReq := orderModel.FromCreateGRPC(request)

	err := reqvalidator.ValidateRequest(orderReq)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Validation failed")
		return nil, status.Errorf(grpcCodes.InvalidArgument, fmt.Sprintf("reqvalidator.ReadRequest: %v", err))
	}

	err = o.useCase.CreateReceiveOrder(ctx, orderReq)
	if err != nil {
		if errors.Is(err, errlst.ErrOrderAlreadyExists) {
			tracing.EventErrorTracer(span, err, "order already exists")
			return nil, status.Errorf(grpcCodes.AlreadyExists, fmt.Sprintf("Failed to create: %v", err))
		}
		tracing.EventErrorTracer(span, err, "Internal server error")
		return nil, status.Errorf(grpcCodes.Internal, fmt.Sprintf("Failed to create: %v", err))
	}

	span.SetStatus(codes.Ok, "Successfully created a order")
	return &abstract.MessageResponse{Message: "Successfully Created Order"}, nil
}

// IssueOrder is
func (o *OrderHandler) IssueOrder(ctx context.Context, request *order_v1.IssueOrderRequest) (*abstract.MessageResponse, error) {
	log.Printf("[order][delivery][IssueOrder]")
	tracer := otel.Tracer("[order][delivery]")
	ctx, span := tracer.Start(ctx, "[IssueOrder]")
	defer span.End()
	issueReq := orderModel.FromIssueGRPC(request)

	err := reqvalidator.ValidateRequest(issueReq)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Validation failed")
		return nil, status.Errorf(grpcCodes.InvalidArgument, fmt.Sprintf("reqvalidator.ReadRequest: %v", err))
	}

	totalAdditionalCost, err := o.useCase.IssueOrders(ctx, issueReq)
	if err != nil {
		if errors.Is(err, errlst.ErrOrderNotFound) {
			tracing.EventErrorTracer(span, err, "order not found")
			return nil, status.Errorf(grpcCodes.NotFound, fmt.Sprintf("Failed to find: %v", err))
		}
		tracing.EventErrorTracer(span, err, "Internal server error")
		return nil, status.Errorf(grpcCodes.Internal, fmt.Sprintf("Failed to find: %v", err))
	}
	span.SetStatus(codes.Ok, "Successfully issued Orders")
	return &abstract.MessageResponse{Message: "Successfully Issued All OrderIDs and Total additional cost:" + strconv.FormatFloat(totalAdditionalCost, 'f', -1, 64)}, nil
}

// ReturnedOrders is
func (o *OrderHandler) ReturnedOrders(ctx context.Context, request *abstract.Page) (*order_v1.ReturnedListResponse, error) {
	log.Printf("[order][delivery][ReturnedOrders]")
	tracer := otel.Tracer("[order][delivery]")
	ctx, span := tracer.Start(ctx, "[ReturnedOrders]")
	defer span.End()
	page := abstractModel.PageFromGRPC(request)

	err := reqvalidator.ValidateRequest(page)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Validation failed")
		return nil, status.Errorf(grpcCodes.InvalidArgument, fmt.Sprintf("reqvalidator.ReadRequest: %v", err))
	}

	returnedList, err := o.useCase.ReturnedOrders(ctx, page)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Internal server error")
		return nil, status.Errorf(grpcCodes.Internal, fmt.Sprintf("Failed to get returned list: %v", err))
	}
	span.SetStatus(codes.Ok, "Successfully got list of returned orders")
	return orderModel.ReturnedListToGRPC(returnedList), nil
}

// AcceptOrder is
func (o *OrderHandler) AcceptOrder(ctx context.Context, request *order_v1.RequestWithClientID) (*abstract.MessageResponse, error) {
	log.Printf("[order][delivery][AcceptOrder]")
	tracer := otel.Tracer("[order][delivery]")
	ctx, span := tracer.Start(ctx, "[AcceptOrder]")
	defer span.End()

	acceptOrderRequest := orderModel.FromAcceptGRPC(request)
	err := reqvalidator.ValidateRequest(&acceptOrderRequest)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Validation failed")
		return nil, status.Errorf(grpcCodes.InvalidArgument, fmt.Sprintf("reqvalidator.ReadRequest: %v", err))
	}

	err = o.useCase.UpdateAcceptOrder(ctx, acceptOrderRequest)
	if err != nil {
		if errors.Is(err, errlst.ErrOrderNotFound) {
			tracing.EventErrorTracer(span, err, "order not found")
			return nil, status.Errorf(grpcCodes.NotFound, fmt.Sprintf("Failed to find order: %v", err))
		}
		tracing.EventErrorTracer(span, err, "Internal server error")
		return nil, status.Errorf(grpcCodes.Internal, fmt.Sprintf("Failed to accept order: %v", err))
	}

	span.SetStatus(codes.Ok, "Successfully accept Order")
	return &abstract.MessageResponse{Message: "Successfully Accepted Order ID with Client ID"}, nil
}

// TurnInOrder is
func (o *OrderHandler) TurnInOrder(ctx context.Context, request *order_v1.OrderIDRequest) (*abstract.MessageResponse, error) {
	log.Printf("[order][delivery][TurnInOrder]")
	tracer := otel.Tracer("[order][delivery]")
	ctx, span := tracer.Start(ctx, "[TurnInOrder]")
	defer span.End()

	if request.OrderID < 0 {
		tracing.EventErrorTracer(span, errors.New("invalid requst id"), "bad request")
		return nil, status.Errorf(grpcCodes.InvalidArgument, "invalid request id")
	}

	err := o.useCase.DeleteReturnedOrder(ctx, request.OrderID)
	if err != nil {
		if errors.Is(err, errlst.ErrOrderNotFound) {
			tracing.EventErrorTracer(span, err, "order not found")
			return nil, status.Errorf(grpcCodes.NotFound, fmt.Sprintf("Failed to find order: %v", err))
		}
		tracing.EventErrorTracer(span, err, "Internal server error")
		return nil, status.Errorf(grpcCodes.Internal, fmt.Sprintf("Failed to turn in order: %v", err))
	}
	span.SetStatus(codes.Ok, "Successfully turn in Order")
	return &abstract.MessageResponse{Message: "Successfully turn in with given order ID"}, nil
}

// OrderList is
func (o *OrderHandler) OrderList(ctx context.Context, request *abstract.Page) (*order_v1.OrderListResponse, error) {
	log.Printf("[order][delivery][OrderList]")
	tracer := otel.Tracer("[order][delivery]")
	ctx, span := tracer.Start(ctx, "[OrderList]")
	defer span.End()

	page := abstractModel.PageFromGRPC(request)

	err := reqvalidator.ValidateRequest(page)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Validation failed")
		return nil, status.Errorf(grpcCodes.InvalidArgument, fmt.Sprintf("reqvalidator.ReadRequest: %v", err))
	}

	orderList, err := o.useCase.OrderList(ctx, page)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Internal server error")
		return nil, status.Errorf(grpcCodes.Internal, fmt.Sprintf("Failed to get order list: %v", err))
	}

	span.SetStatus(codes.Ok, "Successfully got order list")
	return orderModel.ListToGRPC(orderList), nil
}

// UniqueClientList is
func (o *OrderHandler) UniqueClientList(ctx context.Context, request *abstract.Page) (*order_v1.UniqueClientListResponse, error) {
	log.Printf("[order][delivery][UniqueClientsList]")
	tracer := otel.Tracer("[order][delivery]")
	ctx, span := tracer.Start(ctx, "[UniqueClientList]")
	defer span.End()

	clientPagination := abstractModel.PageFromGRPC(request)

	err := reqvalidator.ValidateRequest(&clientPagination)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Validation failed")
		return nil, status.Errorf(grpcCodes.InvalidArgument, fmt.Sprintf("reqvalidator.ReadRequest: %v", err))
	}

	listUniqueClients, err := o.useCase.UniqueClientsList(ctx, clientPagination)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Internal server error")
		return nil, status.Errorf(grpcCodes.Internal, fmt.Sprintf("Failed to find in clients: %v", err))
	}

	span.SetStatus(codes.Ok, "Successfully got unique client list")

	return orderModel.UniqueClientListToGRPC(listUniqueClients), nil
}
