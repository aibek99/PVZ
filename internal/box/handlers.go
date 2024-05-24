package box

import (
	"context"

	"Homework-1/pkg/api/abstract"
	"Homework-1/pkg/api/box_v1"
)

// Handlers is
type Handlers interface {
	CreateBox(context.Context, *box_v1.BoxCreateRequest) (*abstract.MessageResponse, error)
	DeleteBox(context.Context, *box_v1.BoxIDRequest) (*abstract.MessageResponse, error)
	ListBoxes(context.Context, *abstract.Page) (*box_v1.BoxListResponse, error)
	GetBoxByID(context.Context, *box_v1.BoxIDRequest) (*box_v1.BoxAllInfo, error)
}
