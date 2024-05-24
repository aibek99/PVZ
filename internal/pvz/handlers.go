package pvz

import (
	"context"

	"Homework-1/pkg/api/abstract"
	"Homework-1/pkg/api/pvz_v1"
)

// Handlers is
type Handlers interface {
	CreatePVZ(ctx context.Context, request *pvz_v1.PVZCreateRequest) (*abstract.MessageResponse, error)
	GetPVZByID(ctx context.Context, request *pvz_v1.PVZIDRequest) (*pvz_v1.PVZAllInfo, error)
	ListPVZ(ctx context.Context, request *abstract.Page) (*pvz_v1.ListResponse, error)
	UpdatePVZ(ctx context.Context, request *pvz_v1.UpdateRequest) (*abstract.MessageResponse, error)
	DeletePVZ(ctx context.Context, request *pvz_v1.PVZIDRequest) (*abstract.MessageResponse, error)
}
