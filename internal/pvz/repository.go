package pvz

import (
	"context"

	"Homework-1/internal/model/abstract"
	"Homework-1/internal/model/pvz"
)

// Repository is
type Repository interface {
	CreatePVZ(ctx context.Context, pvzData pvz.Data) error
	CheckPVZ(ctx context.Context, pvzData pvz.Data) error
	GetPVZ(ctx context.Context, pvzID int64) (pvz.AllData, error)
	ListPVZ(
		ctx context.Context,
		pvzPaginationData abstract.PageData,
	) ([]pvz.AllData, error)
	CountOfPVZ(ctx context.Context) (int64, error)
	UpdatePVZ(ctx context.Context, updatePVZData pvz.UpdateData) error
	DeletePVZByID(ctx context.Context, pvzID int64) error
}
