package pvz

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	abstractModel "Homework-1/internal/model/abstract"
	"Homework-1/pkg/api/pvz_v1"
)

// Request is
type Request struct {
	Name    string `json:"name" validate:"required,min=4,max=100"`
	Address string `json:"address" validate:"required,min=2"`
	Contact string `json:"contact" validate:"required,phone"`
}

// UpdateRequest is
type UpdateRequest struct {
	ID      int64   `json:"id" validate:"required"`
	Name    *string `json:"name,omitempty"`
	Address *string `json:"address,omitempty"`
	Contact *string `json:"contact,omitempty"`
}

// UpdateData is
type UpdateData struct {
	ID      int64   `db:"id"`
	Name    *string `db:"name"`
	Address *string `db:"address"`
	Contact *string `db:"contact"`
}

// Data is
type Data struct {
	Name    string `db:"name"`
	Address string `db:"address"`
	Contact string `db:"contact"`
}

// AllResponse is
type AllResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Contact   string    `json:"contact"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// AllData is
type AllData struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Address   string    `db:"address"`
	Contact   string    `db:"contact"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// ToStorage is
func (p *Request) ToStorage() Data {
	return Data{
		Name:    p.Name,
		Address: p.Address,
		Contact: p.Contact,
	}
}

// ToPVZServer is
func (p *AllData) ToPVZServer() AllResponse {
	return AllResponse{
		ID:        p.ID,
		Name:      p.Name,
		Address:   p.Address,
		Contact:   p.Contact,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// ToStorage is
func (p *UpdateRequest) ToStorage() UpdateData {
	return UpdateData{
		ID:      p.ID,
		Name:    p.Name,
		Address: p.Address,
		Contact: p.Contact,
	}
}

// FromCreateGRPC is
func FromCreateGRPC(pvz *pvz_v1.PVZ) Request {
	return Request{
		Name:    pvz.Name,
		Address: pvz.Address,
		Contact: pvz.Contact,
	}
}

// InfoToGRPC is
func InfoToGRPC(allResponse AllResponse) *pvz_v1.PVZAllInfo {
	return &pvz_v1.PVZAllInfo{
		ID: allResponse.ID,
		Pvz: &pvz_v1.PVZ{
			Name:    allResponse.Name,
			Address: allResponse.Address,
			Contact: allResponse.Contact,
		},
		CreatedAt: timestamppb.New(allResponse.CreatedAt),
		UpdatedAt: timestamppb.New(allResponse.UpdatedAt),
	}
}

// AllInfoToGRPC is
func AllInfoToGRPC(allResponse []AllResponse) []*pvz_v1.PVZAllInfo {
	response := make([]*pvz_v1.PVZAllInfo, len(allResponse))
	for index, value := range allResponse {
		response[index] = InfoToGRPC(value)
	}
	return response
}

// ListToGRPC is
func ListToGRPC(response abstractModel.PaginatedResponse[AllResponse]) *pvz_v1.ListResponse {
	return &pvz_v1.ListResponse{
		PvzAllInfo: AllInfoToGRPC(response.Items),
		Pagination: abstractModel.PaginationToGRPC(
			abstractModel.Page{
				CurrentPage:  response.CurrentPage,
				ItemsPerPage: response.ItemsPerPage,
			},
			response.TotalItems,
		),
	}
}

// FromUpdateGRPC is
func FromUpdateGRPC(pvz *pvz_v1.UpdateRequest) UpdateRequest {
	return UpdateRequest{
		ID:      pvz.ID,
		Name:    &pvz.Pvz.Name,
		Address: &pvz.Pvz.Address,
		Contact: &pvz.Pvz.Contact,
	}
}
