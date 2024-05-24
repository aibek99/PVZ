package delivery

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	abstractModel "Homework-1/internal/model/abstract"
	pvzModel "Homework-1/internal/model/pvz"
	PVZUseCase "Homework-1/internal/pvz"
	"Homework-1/pkg/api/abstract"
	"Homework-1/pkg/api/pvz_v1"
	"Homework-1/pkg/errlst"
	"Homework-1/pkg/reqvalidator"
	"Homework-1/pkg/tracing"
)

var _ PVZUseCase.Handlers = &PVZHandler{}

// PVZHandler is
type PVZHandler struct {
	useCase PVZUseCase.UseCase
	pvz_v1.UnimplementedPVZServiceServer
}

// NewPVZHandler is
func NewPVZHandler(useCase PVZUseCase.UseCase) *PVZHandler {
	return &PVZHandler{
		useCase: useCase,
	}
}

// CreatePVZ is
func (p *PVZHandler) CreatePVZ(ctx context.Context, request *pvz_v1.PVZCreateRequest) (*abstract.MessageResponse, error) {
	log.Printf("[pvz][delivery][CreatePVZ]")
	tracer := otel.Tracer("[pvz_v1][delivery]")
	ctx, span := tracer.Start(ctx, "[CreatePVZ]")
	defer span.End()

	pvzReq := pvzModel.FromCreateGRPC(request.Pvz)

	err := reqvalidator.ValidateRequest(pvzReq)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Validation failed")
		return nil, status.Errorf(grpcCodes.InvalidArgument, fmt.Sprintf("Failed to read request body: %v", err))
	}

	err = p.useCase.CreatePVZ(ctx, pvzReq)
	if err != nil {
		if errors.Is(err, errlst.ErrPVZAlreadyExists) {
			tracing.EventErrorTracer(span, err, "PVZ already exists")

			return nil, status.Errorf(grpcCodes.AlreadyExists, fmt.Sprintf("Failed to create: %v", err))
		}
		tracing.EventErrorTracer(span, err, "Internal server error")

		return nil, status.Errorf(grpcCodes.Internal, fmt.Sprintf("Failed to create: %v", err))
	}

	span.SetStatus(codes.Ok, "PVZ created successfully")
	return &abstract.MessageResponse{Message: "Successfully Created PVZ\n"}, nil
}

// GetPVZByID is
func (p *PVZHandler) GetPVZByID(ctx context.Context, request *pvz_v1.PVZIDRequest) (*pvz_v1.PVZAllInfo, error) {
	log.Printf("[pvz][delivery][GetPVZByID]")
	tracer := otel.Tracer("[pvz_v1][delivery]")
	ctx, span := tracer.Start(ctx, "[GetPVZByID]")
	defer span.End()

	if request.PvzID < 0 {
		tracing.EventErrorTracer(span, errors.New("invalid request id"), "bad request")
		return nil, status.Errorf(grpcCodes.InvalidArgument, "invalid request id")
	}

	pvzResponse, err := p.useCase.GetPVZ(ctx, request.PvzID)
	if err != nil {
		if errors.Is(err, errlst.ErrPVZNotFound) {
			tracing.EventErrorTracer(span, err, "PVZ not found")
			return nil, status.Errorf(grpcCodes.NotFound, fmt.Sprintf("Error: %v", err))
		}

		tracing.EventErrorTracer(span, err, "internal server error")
		return nil, status.Errorf(grpcCodes.Internal, fmt.Sprintf("Failed to get: %v", err))
	}
	span.SetStatus(codes.Ok, "Successfully received pvz info by ID")
	return pvzModel.InfoToGRPC(pvzResponse), nil
}

// ListPVZ is
func (p *PVZHandler) ListPVZ(ctx context.Context, request *abstract.Page) (*pvz_v1.ListResponse, error) {
	log.Printf("[pvz][delivery][ListPVZ]")
	tracer := otel.Tracer("[pvz_v1][delivery]")
	ctx, span := tracer.Start(ctx, "[ListPVZ]")
	defer span.End()
	page := abstractModel.PageFromGRPC(request)

	err := reqvalidator.ValidateRequest(page)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Validation failed")
		return nil, status.Errorf(grpcCodes.InvalidArgument, fmt.Sprintf("Failed to read request body: %v", err))
	}

	listPVZ, err := p.useCase.ListPVZ(ctx, page)
	if err != nil {
		tracing.EventErrorTracer(span, err, "unable to get list of PVZ")
		return nil, status.Errorf(grpcCodes.Internal, fmt.Sprintf("Failed to get list of PVZ: %v", err))
	}

	span.SetStatus(codes.Ok, "Successfully received pvz list")
	return pvzModel.ListToGRPC(listPVZ), nil
}

// UpdatePVZ is
func (p *PVZHandler) UpdatePVZ(ctx context.Context, request *pvz_v1.UpdateRequest) (*abstract.MessageResponse, error) {
	log.Printf("[pvz][delivery][UpdatePVZ]")
	tracer := otel.Tracer("[pvz_v1][delivery]")
	ctx, span := tracer.Start(ctx, "[UpdatePVZ]")
	defer span.End()

	updatePVZRequest := pvzModel.FromUpdateGRPC(request)

	err := reqvalidator.ValidateRequest(updatePVZRequest)
	if err != nil {
		tracing.EventErrorTracer(span, err, "Validation failed")
		return nil, status.Errorf(grpcCodes.InvalidArgument, fmt.Sprintf("Failed to read request body: %v", err))
	}

	err = p.useCase.UpdatePVZ(ctx, updatePVZRequest)
	if err != nil {
		if errors.Is(err, errlst.ErrPVZNotFound) {
			tracing.EventErrorTracer(span, err, "PVZ not found")
			return nil, status.Errorf(grpcCodes.NotFound, fmt.Sprintf("Error: %v", err))
		}

		tracing.EventErrorTracer(span, err, "internal server error")
		return nil, status.Errorf(grpcCodes.Internal, fmt.Sprintf("Failed to update: %v", err))
	}

	span.SetStatus(codes.Ok, "Successfully updated PVZ")
	return &abstract.MessageResponse{Message: "Successfully Updated PVZ"}, nil
}

// DeletePVZ is
func (p *PVZHandler) DeletePVZ(ctx context.Context, request *pvz_v1.PVZIDRequest) (*abstract.MessageResponse, error) {
	log.Printf("[pvz][delivery][DeletePVZByID]")
	tracer := otel.Tracer("[pvz_v1][delivery]")
	ctx, span := tracer.Start(ctx, "[DeletePVZ]")
	defer span.End()

	if request.PvzID < 0 {
		tracing.EventErrorTracer(span, errors.New("invalid request id"), "bad request")
		return nil, status.Errorf(grpcCodes.InvalidArgument, "invalid request id")
	}

	err := p.useCase.DeletePVZByID(ctx, request.PvzID)
	if err != nil {
		if errors.Is(err, errlst.ErrPVZNotFound) {
			tracing.EventErrorTracer(span, err, "PVZ not found")
			return nil, status.Errorf(grpcCodes.NotFound, fmt.Sprintf("Error: %v", err))
		}
		tracing.EventErrorTracer(span, err, "Internal server error")
		return nil, status.Errorf(grpcCodes.Internal, fmt.Sprintf("Failed to delete: %v", err))
	}

	span.SetStatus(codes.Ok, "Successfully deleted PVZ by ID")
	return &abstract.MessageResponse{Message: "Successfully Deleted PVZ\n"}, nil
}
