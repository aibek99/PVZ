// Package pvz provides the use case logic for managing PVZ (Pickup and Delivery Points).
// This package contains the business logic for creating, updating, and retrieving PVZ data.
// It interacts with the repository layer to perform database operations.
//
//go:generate minimock -g -i UseCase -o ./mock/usecase_mock.go -n UseCaseMock
package pvz

import (
	"context"

	"Homework-1/internal/model/abstract"
	pvzModel "Homework-1/internal/model/pvz"
)

// UseCase is
type UseCase interface {
	CreatePVZ(ctx context.Context, request pvzModel.Request) error
	GetPVZ(ctx context.Context, pvzID int64) (pvzModel.AllResponse, error)
	DeletePVZByID(ctx context.Context, pvzID int64) error
	UpdatePVZ(ctx context.Context, updatePVZRequest pvzModel.UpdateRequest) error
	ListPVZ(ctx context.Context, pvzPagination abstract.Page) (abstract.PaginatedResponse[pvzModel.AllResponse], error)
}
