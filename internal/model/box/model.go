package box

import (
	"time"

	abstractModel "Homework-1/internal/model/abstract"
	"Homework-1/pkg/api/box_v1"
)

// Request is
type Request struct {
	Name    string  `json:"name" validate:"required"`
	Cost    float64 `json:"cost" validate:"required"`
	IsCheck bool    `json:"isCheck"`
	Weight  float64 `json:"weight" validate:"required"`
}

// Data is
type Data struct {
	Name    string  `db:"name"`
	Cost    float64 `db:"cost"`
	IsCheck bool    `db:"is_check"`
	Weight  float64 `db:"weight"`
}

// DataWithOrder is
type DataWithOrder struct {
	Cost        float64 `db:"cost"`
	IsCheck     bool    `db:"is_check"`
	Weight      float64 `db:"weight"`
	OrderWeight float64 `db:"order_weight"`
}

// ToStorage is
func (b *Request) ToStorage() Data {
	return Data{
		Name:    b.Name,
		Cost:    b.Cost,
		IsCheck: b.IsCheck,
		Weight:  b.Weight,
	}
}

// AllData is
type AllData struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Cost      float64   `db:"cost"`
	IsCheck   bool      `db:"is_check"`
	Weight    float64   `db:"weight"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// AllResponse is
type AllResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Cost      float64   `json:"cost"`
	IsCheck   bool      `json:"isCheck"`
	Weight    float64   `json:"weight"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToServer is
func (b *AllData) ToServer() AllResponse {
	return AllResponse{
		ID:        b.ID,
		Name:      b.Name,
		Cost:      b.Cost,
		IsCheck:   b.IsCheck,
		Weight:    b.Weight,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

// FromGRPC is
func FromGRPC(boxGRPC *box_v1.Box) Request {
	return Request{
		Name:    boxGRPC.Name,
		Cost:    boxGRPC.Cost,
		IsCheck: boxGRPC.IsCheck,
		Weight:  boxGRPC.Weight,
	}
}

// InfoToGRPC is
func InfoToGRPC(allResponse AllResponse) *box_v1.BoxAllInfo {
	return &box_v1.BoxAllInfo{
		ID: allResponse.ID,
		Box: &box_v1.Box{
			Name:    allResponse.Name,
			Cost:    allResponse.Cost,
			IsCheck: allResponse.IsCheck,
			Weight:  allResponse.Weight,
		},
		CreatedAt: abstractModel.SafeTimestamp(&allResponse.CreatedAt),
		UpdatedAt: abstractModel.SafeTimestamp(&allResponse.UpdatedAt),
	}
}

// AllInfoToGRPC is
func AllInfoToGRPC(allResponse []AllResponse) []*box_v1.BoxAllInfo {
	var response []*box_v1.BoxAllInfo
	for _, value := range allResponse {
		response = append(response, InfoToGRPC(value))
	}

	return response
}

// ListToGRPC is
func ListToGRPC(response abstractModel.PaginatedResponse[AllResponse]) *box_v1.BoxListResponse {
	return &box_v1.BoxListResponse{
		BoxAllInfo: AllInfoToGRPC(response.Items),
		Pagination: abstractModel.PaginationToGRPC(
			abstractModel.Page{
				CurrentPage:  response.CurrentPage,
				ItemsPerPage: response.ItemsPerPage,
			},
			response.TotalItems,
		)}
}
