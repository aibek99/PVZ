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

	"Homework-1/internal/kafka"
	abstractModel "Homework-1/internal/model/abstract"
	pvzModel "Homework-1/internal/model/pvz"
	"Homework-1/internal/pvz"
	"Homework-1/internal/pvz/mock"
	"Homework-1/pkg/api/abstract"
	"Homework-1/pkg/api/pvz_v1"
	"Homework-1/pkg/constants"
	"Homework-1/pkg/errlst"
)

// TestPVZHandler_CreatePVZ is
func TestPVZHandler_CreatePVZ(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)

	tests := []struct {
		description string
		requestBody *pvz_v1.PVZCreateRequest
		wantResp    *abstract.MessageResponse
		wantErr     error
		useCase     pvz.UseCase
		producer    func() kafka.Producer
	}{
		{
			description: "Request validation failed",

			requestBody: &pvz_v1.PVZCreateRequest{
				Pvz: &pvz_v1.PVZ{
					Name:    "test",
					Address: "kazan",
				},
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.InvalidArgument, "Failed to read request body: Key: 'Request.Contact' Error:Field validation for 'Contact' failed on the 'required' tag"),
			useCase:  mock.NewUseCaseMock(ctrl),
		},
		{
			description: "PVZ already exists", requestBody: &pvz_v1.PVZCreateRequest{
				Pvz: &pvz_v1.PVZ{
					Name:    "test",
					Address: "kazan",
					Contact: "+5654646546",
				},
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.AlreadyExists, "Failed to create: PVZ already exists in PVZ storage"),
			useCase: mock.NewUseCaseMock(ctrl).
				CreatePVZMock.
				When(minimock.AnyContext, pvzModel.Request{
					Name:    "test",
					Address: "kazan",
					Contact: "+5654646546",
				}).
				Then(errlst.ErrPVZAlreadyExists),
		},
		{
			description: "Internal server error",
			requestBody: &pvz_v1.PVZCreateRequest{
				Pvz: &pvz_v1.PVZ{
					Name:    "test",
					Address: "kazan",
					Contact: "+5654646546",
				},
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.Internal, "Failed to create: assert.AnError general error for testing"),
			useCase: mock.NewUseCaseMock(ctrl).
				CreatePVZMock.
				When(minimock.AnyContext, pvzModel.Request{
					Name:    "test",
					Address: "kazan",
					Contact: "+5654646546",
				}).
				Then(assert.AnError),
		},
		{
			description: constants.SuccessfullyCreatedPVZ,
			requestBody: &pvz_v1.PVZCreateRequest{
				Pvz: &pvz_v1.PVZ{
					Name:    "test",
					Address: "kazan",
					Contact: "+5654646546",
				},
			},
			wantResp: &abstract.MessageResponse{
				Message: "Successfully Created PVZ\n",
			},
			wantErr: nil,
			useCase: mock.NewUseCaseMock(ctrl).
				CreatePVZMock.
				When(minimock.AnyContext, pvzModel.Request{
					Name:    "test",
					Address: "kazan",
					Contact: "+5654646546",
				}).
				Then(nil),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			response, err := NewPVZHandler(tt.useCase).CreatePVZ(context.Background(), tt.requestBody)

			assert.Equal(t, tt.wantResp, response)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

// TestPVZHandler_DeletePVZByID is
func TestPVZHandler_DeletePVZByID(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)

	tests := []struct {
		description string
		requestID   *pvz_v1.PVZIDRequest
		wantResp    *abstract.MessageResponse
		wantErr     error
		useCase     pvz.UseCase
	}{
		{
			description: "Successfully Deleted PVZ",
			requestID: &pvz_v1.PVZIDRequest{
				PvzID: 1,
			},
			wantResp: &abstract.MessageResponse{
				Message: "Successfully Deleted PVZ\n",
			},
			wantErr: nil,
			useCase: mock.NewUseCaseMock(ctrl).
				DeletePVZByIDMock.
				When(minimock.AnyContext, 1).
				Then(nil),
		},
		{
			description: "PVZ not found",
			requestID: &pvz_v1.PVZIDRequest{
				PvzID: 1,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.NotFound, "Error: PVZ not found"),
			useCase: mock.NewUseCaseMock(ctrl).
				DeletePVZByIDMock.
				When(minimock.AnyContext, 1).
				Then(errlst.ErrPVZNotFound),
		},
		{
			description: "Internal server error",
			requestID: &pvz_v1.PVZIDRequest{
				PvzID: 1,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.Internal, "Failed to delete: assert.AnError general error for testing"),
			useCase: mock.NewUseCaseMock(ctrl).
				DeletePVZByIDMock.
				When(minimock.AnyContext, 1).
				Then(assert.AnError),
		},
		{
			description: "Unable to get pvzID",
			requestID: &pvz_v1.PVZIDRequest{
				PvzID: -1,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.InvalidArgument, "invalid request id"),
			useCase:  mock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			response, err := NewPVZHandler(tt.useCase).DeletePVZ(context.Background(), tt.requestID)

			assert.Equal(t, tt.wantResp, response)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

// TestPVZHandler_ListPVZ is
func TestPVZHandler_ListPVZ(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	fixedTime := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		description string
		requestBody *abstract.Page
		wantResp    *pvz_v1.ListResponse
		wantErr     error
		useCase     pvz.UseCase
	}{
		{
			description: "Successfully got list of PVZ",
			requestBody: &abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: &pvz_v1.ListResponse{
				PvzAllInfo: []*pvz_v1.PVZAllInfo{
					{
						ID: 1,
						Pvz: &pvz_v1.PVZ{
							Name:    "Sample PVZ",
							Address: "kazan",
							Contact: "+6546545654",
						},
						CreatedAt: timestamppb.New(fixedTime),
						UpdatedAt: timestamppb.New(fixedTime),
					},
				},
				Pagination: &abstract.Pagination{
					Page: &abstract.Page{
						CurrentPage:  1,
						ItemsPerPage: 10,
					},
					TotalItems: 1,
				},
			},
			wantErr: nil,
			useCase: mock.NewUseCaseMock(ctrl).ListPVZMock.
				When(minimock.AnyContext, abstractModel.Page{
					CurrentPage:  1,
					ItemsPerPage: 10,
				}).
				Then(
					abstractModel.PaginatedResponse[pvzModel.AllResponse]{
						Items: []pvzModel.AllResponse{
							{
								ID:        1,
								Name:      "Sample PVZ",
								Address:   "kazan",
								Contact:   "+6546545654",
								CreatedAt: fixedTime,
								UpdatedAt: fixedTime,
							},
						},
						CurrentPage:  1,
						ItemsPerPage: 10,
						TotalItems:   1,
					},
					nil,
				),
		},
		{
			description: "Empty pvz list",
			requestBody: &abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: &pvz_v1.ListResponse{
				PvzAllInfo: []*pvz_v1.PVZAllInfo{},
				Pagination: &abstract.Pagination{
					Page: &abstract.Page{
						CurrentPage:  0,
						ItemsPerPage: 0,
					},
					TotalItems: 0,
				},
			},
			wantErr: nil,
			useCase: mock.NewUseCaseMock(ctrl).
				ListPVZMock.
				When(minimock.AnyContext, abstractModel.Page{
					CurrentPage:  1,
					ItemsPerPage: 10,
				}).
				Then(
					abstractModel.PaginatedResponse[pvzModel.AllResponse]{},
					nil,
				),
		},
		{
			description: "Internal server error",
			requestBody: &abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.Internal, "Failed to get list of PVZ: assert.AnError general error for testing"),
			useCase: mock.NewUseCaseMock(ctrl).
				ListPVZMock.
				When(minimock.AnyContext, abstractModel.Page{
					CurrentPage:  1,
					ItemsPerPage: 10,
				}).
				Then(
					abstractModel.PaginatedResponse[pvzModel.AllResponse]{},
					assert.AnError,
				),
		},
		{
			description: "Request validation failed",
			requestBody: &abstract.Page{
				CurrentPage: 1,
			},
			wantResp: nil,
			wantErr: status.Errorf(
				codes.InvalidArgument,
				"Failed to read request body: Key: 'Page.ItemsPerPage' Error:Field validation for 'ItemsPerPage' failed on the 'required' tag",
			),
			useCase: mock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			response, err := NewPVZHandler(tt.useCase).ListPVZ(context.Background(), tt.requestBody)

			assert.Equal(t, tt.wantResp, response)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

// TestPVZHandler_GetPVZByID is
func TestPVZHandler_GetPVZByID(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	fixedTime := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		description string
		requestID   *pvz_v1.PVZIDRequest
		wantResp    *pvz_v1.PVZAllInfo
		wantErr     error
		useCase     pvz.UseCase
	}{
		{
			description: "Successfully Got a PVZ",
			requestID: &pvz_v1.PVZIDRequest{
				PvzID: 1,
			},
			wantResp: &pvz_v1.PVZAllInfo{
				ID: 1,
				Pvz: &pvz_v1.PVZ{
					Name:    "Sample PVZ",
					Address: "kazan",
					Contact: "+6546545654",
				},
				CreatedAt: timestamppb.New(fixedTime),
				UpdatedAt: timestamppb.New(fixedTime),
			},
			wantErr: nil,
			useCase: mock.NewUseCaseMock(ctrl).
				GetPVZMock.
				When(minimock.AnyContext, 1).
				Then(
					pvzModel.AllResponse{
						ID:        1,
						Name:      "Sample PVZ",
						Address:   "kazan",
						Contact:   "+6546545654",
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
					},
					nil,
				),
		},
		{
			description: "PVZ not found",
			requestID: &pvz_v1.PVZIDRequest{
				PvzID: 1,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.NotFound, "Error: PVZ not found"),
			useCase: mock.NewUseCaseMock(ctrl).
				GetPVZMock.
				When(minimock.AnyContext, 1).
				Then(
					pvzModel.AllResponse{},
					errlst.ErrPVZNotFound,
				),
		},
		{
			description: "Internal Server error",
			requestID: &pvz_v1.PVZIDRequest{
				PvzID: 1,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.Internal, "Failed to get: assert.AnError general error for testing"),
			useCase: mock.NewUseCaseMock(ctrl).
				GetPVZMock.
				When(minimock.AnyContext, 1).
				Then(
					pvzModel.AllResponse{},
					assert.AnError,
				),
		},
		{
			description: "Unable to get pvzID",
			requestID: &pvz_v1.PVZIDRequest{
				PvzID: -1,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.InvalidArgument, "invalid request id"),
			useCase:  mock.NewUseCaseMock(ctrl),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			response, err := NewPVZHandler(tt.useCase).GetPVZByID(context.Background(), tt.requestID)

			assert.Equal(t, tt.wantResp, response)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

// TestPVZHandler_UpdatePVZ is
func TestPVZHandler_UpdatePVZ(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)

	argument := pvz_v1.PVZ{
		Name:    "test",
		Address: "kazan",
		Contact: "+5654646546",
	}

	tests := []struct {
		description string
		requestBody *pvz_v1.UpdateRequest
		wantResp    *abstract.MessageResponse
		wantErr     error
		useCase     pvz.UseCase
	}{
		{
			description: "Request validation failed",
			requestBody: &pvz_v1.UpdateRequest{
				Pvz: &argument,
			},
			wantResp: nil,
			wantErr: status.Errorf(
				codes.InvalidArgument,
				"Failed to read request body: Key: 'UpdateRequest.ID' Error:Field validation for 'ID' failed on the 'required' tag",
			),
			useCase: mock.NewUseCaseMock(ctrl),
		},
		{
			description: "PVZ not found",
			requestBody: &pvz_v1.UpdateRequest{
				ID:  1,
				Pvz: &argument,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.NotFound, "Error: PVZ not found"),
			useCase: mock.NewUseCaseMock(ctrl).
				UpdatePVZMock.
				When(minimock.AnyContext, pvzModel.UpdateRequest{
					ID:      1,
					Name:    &argument.Name,
					Address: &argument.Address,
					Contact: &argument.Contact,
				}).
				Then(errlst.ErrPVZNotFound),
		},
		{
			description: "Internal server error",
			requestBody: &pvz_v1.UpdateRequest{
				ID:  1,
				Pvz: &argument,
			},
			wantResp: nil,
			wantErr:  status.Errorf(codes.Internal, "Failed to update: assert.AnError general error for testing"),
			useCase: mock.NewUseCaseMock(ctrl).
				UpdatePVZMock.
				When(minimock.AnyContext, pvzModel.UpdateRequest{
					ID:      1,
					Name:    &argument.Name,
					Address: &argument.Address,
					Contact: &argument.Contact,
				}).
				Then(assert.AnError),
		},
		{
			description: "Successfully Updated PVZ",
			requestBody: &pvz_v1.UpdateRequest{
				ID:  1,
				Pvz: &argument,
			},
			wantResp: &abstract.MessageResponse{Message: "Successfully Updated PVZ"},
			wantErr:  nil,
			useCase: mock.NewUseCaseMock(ctrl).
				UpdatePVZMock.
				When(minimock.AnyContext, pvzModel.UpdateRequest{
					ID:      1,
					Name:    &argument.Name,
					Address: &argument.Address,
					Contact: &argument.Contact,
				}).
				Then(nil),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			response, err := NewPVZHandler(tt.useCase).UpdatePVZ(context.Background(), tt.requestBody)

			assert.Equal(t, tt.wantResp, response)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
