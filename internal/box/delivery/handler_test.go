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

	"Homework-1/internal/box"
	boxMock "Homework-1/internal/box/mock"
	"Homework-1/internal/kafka"
	abstractModel "Homework-1/internal/model/abstract"
	boxModel "Homework-1/internal/model/box"
	"Homework-1/pkg/api/abstract"
	"Homework-1/pkg/api/box_v1"
	"Homework-1/pkg/errlst"
)

// TestBoxHandler_CreateBox is
func TestBoxHandler_CreateBox(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)

	tests := []*struct {
		description string
		requestBody box_v1.BoxCreateRequest
		wantResp    *abstract.MessageResponse
		wantErr     error
		useCase     box.UseCase
	}{
		{
			description: "Successfully Created Box",
			requestBody: box_v1.BoxCreateRequest{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12,
					IsCheck: true,
					Weight:  12.1,
				},
			},
			wantResp: &abstract.MessageResponse{Message: "1"},
			wantErr:  nil,
			useCase: boxMock.NewUseCaseMock(ctrl).
				CreateBoxMock.
				When(
					minimock.AnyContext,
					boxModel.Request{
						Name:    "test",
						Cost:    12,
						IsCheck: true,
						Weight:  12.1,
					}).Then(1, nil),
		},
		{
			description: "Box already exists",
			requestBody: box_v1.BoxCreateRequest{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12,
					IsCheck: true,
					Weight:  12.1,
				},
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.AlreadyExists, "Failed to create: pq: duplicate key value violates unique constraint \"box_name_key\""),
			useCase: boxMock.NewUseCaseMock(ctrl).
				CreateBoxMock.
				When(
					minimock.AnyContext,
					boxModel.Request{
						Name:    "test",
						Cost:    12,
						IsCheck: true,
						Weight:  12.1,
					}).Then(-1, errlst.ErrBoxAlreadyExists),
		},
		{
			description: "Internal server error",
			requestBody: box_v1.BoxCreateRequest{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12,
					IsCheck: true,
					Weight:  12.1,
				},
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.Internal, "Failed to create: assert.AnError general error for testing"),
			useCase: boxMock.NewUseCaseMock(ctrl).CreateBoxMock.When(
				minimock.AnyContext,
				boxModel.Request{
					Name:    "test",
					Cost:    12,
					IsCheck: true,
					Weight:  12.1,
				}).Then(-1, assert.AnError),
		},
		{
			description: "Request validation failed",
			requestBody: box_v1.BoxCreateRequest{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12,
					IsCheck: true,
				},
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.InvalidArgument, "reqvalidator.ValidateRequest Key: 'Request.Weight' Error:Field validation for 'Weight' failed on the 'required' tag"),
			useCase:  boxMock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			response, err := NewBoxHandler(tt.useCase).CreateBox(context.Background(), &tt.requestBody)

			assert.Equal(t, tt.wantResp, response)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

// TestBoxHandler_DeleteBox is
func TestBoxHandler_DeleteBox(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)
	tests := []*struct {
		description string
		requestID   box_v1.BoxIDRequest
		wantResp    *abstract.MessageResponse
		wantErr     error
		useCase     box.UseCase
	}{
		{
			description: "Successfully Deleted Box",
			requestID:   box_v1.BoxIDRequest{BoxID: 1},
			wantResp:    &abstract.MessageResponse{Message: "Successfully Deleted Box\n"},
			wantErr:     nil,
			useCase:     boxMock.NewUseCaseMock(ctrl).DeleteBoxByIDMock.When(minimock.AnyContext, 1).Then(nil),
		},
		{
			description: "Box not found",
			requestID:   box_v1.BoxIDRequest{BoxID: 1},
			wantResp:    nil,
			wantErr:     status.Errorf(codes.NotFound, "Error: Box not found"),
			useCase: boxMock.NewUseCaseMock(ctrl).DeleteBoxByIDMock.
				When(minimock.AnyContext, 1).
				Then(errlst.ErrBoxNotFound),
		},
		{
			description: "Internal server error",
			requestID:   box_v1.BoxIDRequest{BoxID: 1},
			wantResp:    nil,
			wantErr:     status.Errorf(codes.Internal, "Failed to delete: assert.AnError general error for testing"),
			useCase:     boxMock.NewUseCaseMock(ctrl).DeleteBoxByIDMock.When(minimock.AnyContext, 1).Then(assert.AnError),
		},
		{
			description: "Unable to parse boxID",
			requestID:   box_v1.BoxIDRequest{BoxID: -1},
			wantResp:    nil,
			wantErr:     status.Errorf(codes.InvalidArgument, "invalid request id"),
			useCase:     boxMock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			response, err := NewBoxHandler(tt.useCase).DeleteBox(context.Background(), &tt.requestID)

			assert.Equal(t, tt.wantResp, response)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

// TestBoxHandler_ListBox is
func TestBoxHandler_ListBox(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	fixedTime := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
	tests := []*struct {
		description string
		requestBody abstract.Page
		wantResp    *box_v1.BoxListResponse
		wantErr     error
		useCase     box.UseCase
		producer    func() kafka.Producer
	}{
		{
			description: "Successfully got list of box",
			requestBody: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: &box_v1.BoxListResponse{BoxAllInfo: []*box_v1.BoxAllInfo{
				{
					ID: 1,
					Box: &box_v1.Box{
						Name:    "Sample Box",
						Cost:    100.0,
						IsCheck: true,
						Weight:  10.1,
					},
					CreatedAt: timestamppb.New(fixedTime),
					UpdatedAt: timestamppb.New(fixedTime),
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
			useCase: boxMock.NewUseCaseMock(ctrl).ListBoxesMock.
				When(
					minimock.AnyContext,
					abstractModel.Page{
						CurrentPage:  1,
						ItemsPerPage: 10,
					}).
				Then(
					abstractModel.PaginatedResponse[boxModel.AllResponse]{
						Items: []boxModel.AllResponse{
							{
								ID:        1,
								Name:      "Sample Box",
								Cost:      100.0,
								IsCheck:   true,
								Weight:    10.1,
								CreatedAt: fixedTime,
								UpdatedAt: fixedTime,
							},
						},
						CurrentPage:  1,
						ItemsPerPage: 1,
						TotalItems:   1,
					}, nil),
		},
		{
			description: "Empty box list",
			requestBody: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: &box_v1.BoxListResponse{BoxAllInfo: nil,
				Pagination: &abstract.Pagination{
					Page: &abstract.Page{
						CurrentPage:  0,
						ItemsPerPage: 0,
					},
					TotalItems: 0,
				},
			},
			wantErr: nil,
			useCase: boxMock.NewUseCaseMock(ctrl).ListBoxesMock.When(minimock.AnyContext, abstractModel.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			}).Then(abstractModel.PaginatedResponse[boxModel.AllResponse]{}, nil),
		},
		{
			description: "Internal server error",
			requestBody: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.Internal, "Failed to get list of Box: assert.AnError general error for testing"),
			useCase: boxMock.NewUseCaseMock(ctrl).ListBoxesMock.
				When(
					minimock.AnyContext, abstractModel.Page{
						CurrentPage:  1,
						ItemsPerPage: 10,
					}).Then(
				abstractModel.PaginatedResponse[boxModel.AllResponse]{},
				assert.AnError),
		},
		{
			description: "Request validation failed",
			requestBody: abstract.Page{
				CurrentPage: 1,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.InvalidArgument, "reqvalidator.ValidateRequest Key: 'Page.ItemsPerPage' Error:Field validation for 'ItemsPerPage' failed on the 'required' tag"),
			useCase:  boxMock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			response, err := NewBoxHandler(tt.useCase).ListBoxes(context.Background(), &tt.requestBody)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantResp, response)
		})
	}
}

// TestBoxHandler_GetBoxByID is
func TestBoxHandler_GetBoxByID(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	fixedTime := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
	tests := []*struct {
		description string
		requestID   box_v1.BoxIDRequest
		wantResp    *box_v1.BoxAllInfo
		wantErr     error
		useCase     box.UseCase
	}{
		{
			description: "Successfully got BoxData",
			requestID:   box_v1.BoxIDRequest{BoxID: 1},
			wantResp: &box_v1.BoxAllInfo{
				ID: 1,
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12.1,
					IsCheck: true,
					Weight:  10.1,
				},
				CreatedAt: timestamppb.New(fixedTime),
				UpdatedAt: timestamppb.New(fixedTime),
			},
			wantErr: nil,
			useCase: boxMock.NewUseCaseMock(ctrl).GetBoxMock.
				When(minimock.AnyContext, 1).
				Then(
					boxModel.AllResponse{
						ID:        1,
						Name:      "test",
						Cost:      12.1,
						IsCheck:   true,
						Weight:    10.1,
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
					}, nil,
				),
		},
		{
			description: "Box not found",
			requestID:   box_v1.BoxIDRequest{BoxID: 1},
			wantResp:    nil,
			wantErr:     status.Errorf(codes.NotFound, "Error: Box not found"),
			useCase: boxMock.NewUseCaseMock(ctrl).GetBoxMock.
				When(minimock.AnyContext, 1).
				Then(boxModel.AllResponse{}, errlst.ErrBoxNotFound),
		},
		{
			description: "Internal server error",
			requestID:   box_v1.BoxIDRequest{BoxID: 1},
			wantResp:    nil,
			wantErr:     status.Errorf(codes.Internal, "Failed to get: assert.AnError general error for testing"),
			useCase: boxMock.NewUseCaseMock(ctrl).GetBoxMock.
				When(minimock.AnyContext, 1).
				Then(boxModel.AllResponse{}, assert.AnError),
		},
		{
			description: "Unable to parse boxID",
			requestID:   box_v1.BoxIDRequest{BoxID: -1},
			wantResp:    nil,
			wantErr:     status.Errorf(codes.InvalidArgument, "invalid request id"),
			useCase:     boxMock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			response, err := NewBoxHandler(tt.useCase).GetBoxByID(context.Background(), &tt.requestID)

			assert.Equal(t, tt.wantResp, response)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
