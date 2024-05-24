// Package box ...
//
//go:generate minimock -g -i UseCase -o ./mock/usecase_mock.go -n UseCaseMock
package box

import (
	"context"

	"Homework-1/internal/model/abstract"
	boxModel "Homework-1/internal/model/box"
)

// UseCase is
type UseCase interface {
	CreateBox(ctx context.Context, request boxModel.Request) (int64, error)
	DeleteBoxByID(ctx context.Context, boxID int64) error
	ListBoxes(ctx context.Context, boxPagination abstract.Page) (abstract.PaginatedResponse[boxModel.AllResponse], error)
	GetBox(ctx context.Context, boxID int64) (boxModel.AllResponse, error)
}
