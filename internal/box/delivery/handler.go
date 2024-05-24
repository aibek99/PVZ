package delivery

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	boxUseCase "Homework-1/internal/box"
	abstractModel "Homework-1/internal/model/abstract"
	boxModel "Homework-1/internal/model/box"
	"Homework-1/pkg/api/abstract"
	"Homework-1/pkg/api/box_v1"
	"Homework-1/pkg/errlst"
	"Homework-1/pkg/reqvalidator"
	"Homework-1/pkg/tracing"
)

var (
	_ boxUseCase.Handlers = (*BoxHandler)(nil)
)

// BoxHandler is
type BoxHandler struct {
	useCase boxUseCase.UseCase
	box_v1.UnimplementedBoxServiceServer
}

// NewBoxHandler is
func NewBoxHandler(useCase boxUseCase.UseCase) *BoxHandler {
	return &BoxHandler{
		useCase: useCase,
	}
}

// CreateBox is
func (b *BoxHandler) CreateBox(ctx context.Context, request *box_v1.BoxCreateRequest) (*abstract.MessageResponse, error) {
	log.Printf("[box_v1][delivery][CreateBox]")
	tracer := otel.Tracer("[box_v1][delivery]")
	ctx, span := tracer.Start(ctx, "[CreateBox]")
	defer span.End()

	boxReq := boxModel.FromGRPC(request.Box)

	err := reqvalidator.ValidateRequest(boxReq)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Validation failed")

		return nil, status.Errorf(grpcCodes.InvalidArgument, fmt.Sprintf("reqvalidator.ValidateRequest %v", err))
	}

	id, err := b.useCase.CreateBox(ctx, boxReq)
	if err != nil {
		if strings.Contains(err.Error(), errlst.ErrBoxAlreadyExists.Error()) {
			tracing.EventErrorTracer(span, err, "Box already exists")

			return nil, status.Errorf(grpcCodes.AlreadyExists, "Failed to create: %v", err)
		}
		tracing.EventErrorTracer(span, err, "Internal server error")

		return nil, status.Errorf(grpcCodes.Internal, "Failed to create: %v", err)
	}

	span.SetStatus(codes.Ok, "Box created successfully")
	return &abstract.MessageResponse{Message: strconv.FormatInt(id, 10)}, nil
}

// DeleteBox is
func (b *BoxHandler) DeleteBox(ctx context.Context, request *box_v1.BoxIDRequest) (*abstract.MessageResponse, error) {
	log.Printf("[box_v1][delivery][DeleteBoxByID]")
	tracer := otel.Tracer("[box_v1][delivery]")
	ctx, span := tracer.Start(ctx, "[DeleteBox]")
	defer span.End()

	if request.BoxID < 0 {
		tracing.EventErrorTracer(span, errors.New("invalid request id"), "bad request")
		return nil, status.Errorf(grpcCodes.InvalidArgument, "invalid request id")
	}

	err := b.useCase.DeleteBoxByID(ctx, request.BoxID)
	if err != nil {
		if errors.Is(err, errlst.ErrBoxNotFound) {
			tracing.EventErrorTracer(span, err, "Box not found")
			return nil, status.Errorf(grpcCodes.NotFound, "Error: %v", err)
		}

		tracing.EventErrorTracer(span, err, "Internal server error")
		return nil, status.Errorf(grpcCodes.Internal, "Failed to delete: %v", err)
	}

	span.SetStatus(codes.Ok, "Successfully deleted box by ID")
	return &abstract.MessageResponse{Message: "Successfully Deleted Box\n"}, nil
}

// ListBoxes is
func (b *BoxHandler) ListBoxes(ctx context.Context, request *abstract.Page) (*box_v1.BoxListResponse, error) {
	log.Printf("[box_v1][handler][ListBoxes]")
	tracer := otel.Tracer("[box_v1][delivery]")
	ctx, span := tracer.Start(ctx, "[ListBoxes]")
	defer span.End()
	page := abstractModel.PageFromGRPC(request)

	err := reqvalidator.ValidateRequest(page)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Validation failed")
		return nil, status.Errorf(grpcCodes.InvalidArgument, fmt.Sprintf("reqvalidator.ValidateRequest %v", err))
	}

	listOfBox, err := b.useCase.ListBoxes(ctx, page)
	if err != nil {
		tracing.EventErrorTracer(span, err, "unable to get list of Box")
		return nil, status.Errorf(grpcCodes.Internal, "Failed to get list of Box: %v", err)
	}

	span.SetStatus(codes.Ok, "Successfully received box list")
	return boxModel.ListToGRPC(listOfBox), nil
}

// GetBoxByID is
func (b *BoxHandler) GetBoxByID(ctx context.Context, request *box_v1.BoxIDRequest) (*box_v1.BoxAllInfo, error) {
	log.Printf("[box_v1][delivery][GetBoxByID]")
	tracer := otel.Tracer("[box_v1][delivery]")
	ctx, span := tracer.Start(ctx, "[GetBoxByID]")
	defer span.End()

	if request == nil || request.BoxID < 0 {
		tracing.EventErrorTracer(span, errors.New("invalid request id"), "bad request")
		return nil, status.Errorf(grpcCodes.InvalidArgument, "invalid request id")
	}

	boxResponse, err := b.useCase.GetBox(ctx, request.BoxID)
	if err != nil {
		if errors.Is(err, errlst.ErrBoxNotFound) {
			tracing.EventErrorTracer(span, err, "Box not found")
			return nil, status.Errorf(grpcCodes.NotFound, fmt.Sprintf("Error: %v", err))
		}
		tracing.EventErrorTracer(span, err, "internal server error")
		return nil, status.Errorf(grpcCodes.Internal, fmt.Sprintf("Failed to get: %v", err))
	}

	span.SetStatus(codes.Ok, "Successfully received box info by ID")
	return boxModel.InfoToGRPC(boxResponse), nil
}
