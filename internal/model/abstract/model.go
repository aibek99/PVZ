package abstract

import (
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	redisModel "Homework-1/internal/model/redis"
	"Homework-1/pkg/api/abstract"
)

// PageData is
type PageData struct {
	CurrentPage  int64 `db:"current_page"`
	ItemsPerPage int64 `db:"items_per_page"`
}

// Page is
type Page struct {
	CurrentPage  int64 `json:"currentPage" validate:"required,gte=0"`
	ItemsPerPage int64 `json:"itemsPerPage" validate:"required,gte=0,lt=101"`
}

// PaginatedResponse is
type PaginatedResponse[T any] struct {
	Items        []T   `json:"items"`
	CurrentPage  int64 `json:"currentPage"`
	ItemsPerPage int64 `json:"itemsPerPage"`
	TotalItems   int64 `json:"totalItems"`
}

// PaginatedResponseData is
type PaginatedResponseData[T any] struct {
	Items        []T   `db:"items"`
	CurrentPage  int64 `db:"current_page"`
	ItemsPerPage int64 `db:"items_per_page"`
	TotalItems   int64 `db:"total_items"`
}

// ToStorage is
func (p *Page) ToStorage() PageData {
	return PageData{
		CurrentPage:  p.CurrentPage,
		ItemsPerPage: p.ItemsPerPage,
	}
}

// PageFromGRPC is
func PageFromGRPC(pageGRPC *abstract.Page) Page {
	return Page{
		CurrentPage:  pageGRPC.CurrentPage,
		ItemsPerPage: pageGRPC.ItemsPerPage,
	}
}

// PaginationToGRPC is
func PaginationToGRPC(page Page, totalItems int64) *abstract.Pagination {
	return &abstract.Pagination{
		Page: &abstract.Page{
			CurrentPage:  page.CurrentPage,
			ItemsPerPage: page.ItemsPerPage,
		},
		TotalItems: totalItems,
	}
}

// CacheArgument is
type CacheArgument struct {
	ObjectType string
	ObjectID   int64
}

// ToCacheStorage is
func (c *CacheArgument) ToCacheStorage() redisModel.CacheKey {
	return redisModel.CacheKey{
		ObjectType: c.ObjectType,
		ID:         strconv.Itoa(int(c.ObjectID)),
	}
}

// SafeTimestamp is
func SafeTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
