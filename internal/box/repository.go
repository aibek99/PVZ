package box

import (
	"context"

	"Homework-1/internal/model/abstract"
	"Homework-1/internal/model/box"
)

// Repository is
type Repository interface {
	CreateBox(ctx context.Context, box box.Data) (int64, error)
	DeleteBoxByID(ctx context.Context, id int64) error
	ListBoxes(ctx context.Context, boxPagination abstract.PageData) ([]box.AllData, error)
	CountBoxes(ctx context.Context) (int64, error)
	GetBox(ctx context.Context, id int64) (box.AllData, error)
}
