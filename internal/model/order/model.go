package order

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	abstractModel "Homework-1/internal/model/abstract"
	"Homework-1/pkg/api/order_v1"
)

// Request is - receive order -expire in days, -orderid, -userid
type Request struct {
	ExpireTimeDuration int     `json:"expireTimeDuration" validate:"required"`
	OrderID            int64   `json:"orderID" validate:"required"`
	ClientID           int64   `json:"clientID" validate:"required"`
	Weight             float64 `json:"weight" validate:"required"`
	BoxID              int64   `json:"boxID" validate:"required"`
}

// RequestData is
type RequestData struct {
	ExpireTimeDuration int     `db:"expire_time_duration"`
	OrderID            int64   `db:"order_ID" `
	ClientID           int64   `db:"client_ID"`
	Weight             float64 `db:"weight"`
	BoxID              int64   `db:"box_id"`
}

// ToStorage is
func (o *Request) ToStorage() RequestData {
	return RequestData{
		ExpireTimeDuration: o.ExpireTimeDuration,
		OrderID:            o.OrderID,
		ClientID:           o.ClientID,
		Weight:             o.Weight,
		BoxID:              o.BoxID,
	}
}

// RequestOrderIDs is
type RequestOrderIDs struct {
	OrderIDs []int64 `json:"orderIDs" validate:"required,min=1"`
}

// WithBoxID is
type WithBoxID struct {
	OrderID int64 `json:"orderID" validate:"required"`
	BoxID   int64 `json:"boxID" validate:"required"`
}

// RequestOrderIDsData is
type RequestOrderIDsData struct {
	OrderIDs []int64 `db:"orderIDs"`
}

// ToStorage is
func (o *RequestOrderIDs) ToStorage() RequestOrderIDsData {
	orderIDs := make([]int64, len(o.OrderIDs))
	_ = copy(orderIDs, o.OrderIDs)
	return RequestOrderIDsData{OrderIDs: orderIDs}
}

// ReturnedRequestOrder is
type ReturnedRequestOrder struct {
	RequestOrderIDs
	abstractModel.Page
}

// ReturnedRequestOrderData is
type ReturnedRequestOrderData struct {
	RequestOrderIDsData
	abstractModel.PageData
}

// ToStorage is
func (r *ReturnedRequestOrder) ToStorage() ReturnedRequestOrderData {
	return ReturnedRequestOrderData{
		RequestOrderIDsData: r.RequestOrderIDs.ToStorage(),
		PageData:            r.Page.ToStorage(),
	}
}

// ReturnedData is
type ReturnedData struct {
	OrderID    int64     `db:"order_id"`
	ClientID   int64     `db:"client_id"`
	ReturnedAt time.Time `db:"returned_at"`
}

// ReturnedResponse is
type ReturnedResponse struct {
	OrderID    int64     `json:"orderID"`
	ClientID   int64     `json:"clientID"`
	ReturnedAt time.Time `json:"returnedAt"`
}

// ToServer is
func (r *ReturnedData) ToServer() ReturnedResponse {
	return ReturnedResponse{
		OrderID:    r.OrderID,
		ClientID:   r.ClientID,
		ReturnedAt: r.ReturnedAt,
	}
}

// RequestWithClientID is
type RequestWithClientID struct {
	OrderID  int64 `json:"orderID" validate:"required"`
	ClientID int64 `json:"clientID" validate:"required"`
}

// RequestWithClientIDData is
type RequestWithClientIDData struct {
	OrderID  int64 `db:"order_id"`
	ClientID int64 `db:"client_id"`
}

// ToStorage is
func (o *RequestWithClientID) ToStorage() RequestWithClientIDData {
	return RequestWithClientIDData{
		OrderID:  o.OrderID,
		ClientID: o.ClientID,
	}
}

// AllResponse is
type AllResponse struct {
	OrderID    int64      `json:"orderID"`
	ClientID   int64      `json:"clientID"`
	AcceptedAt *time.Time `json:"acceptedAt"`
	IssuedAt   *time.Time `json:"issuedAt"`
	ExpiresAt  *time.Time `json:"expiresAt"`
	Weight     float64    `json:"weight"`
	BoxID      int64      `json:"boxID"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// AllResponseData is
type AllResponseData struct {
	OrderID    int64      `db:"order_id"`
	Weight     float64    `db:"weight"`
	ClientID   int64      `db:"client_id"`
	BoxID      int64      `db:"box_id"`
	AcceptedAt *time.Time `db:"accepted_at"`
	IssuedAt   *time.Time `db:"issued_at"`
	ExpiresAt  *time.Time `db:"expires_at"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

// ToServer is
func (o *AllResponseData) ToServer() AllResponse {
	return AllResponse{
		OrderID:    o.OrderID,
		ClientID:   o.ClientID,
		AcceptedAt: o.AcceptedAt,
		IssuedAt:   o.IssuedAt,
		ExpiresAt:  o.ExpiresAt,
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
		Weight:     o.Weight,
		BoxID:      o.BoxID,
	}
}

// ListUniqueClients is
type ListUniqueClients struct {
	ClientID int64 `json:"clientID"`
}

// ListUniqueClientsData is
type ListUniqueClientsData struct {
	ClientIDs int64 `db:"client_id"`
}

// ToServer is
func (l *ListUniqueClientsData) ToServer() ListUniqueClients {
	return ListUniqueClients{ClientID: l.ClientIDs}
}

// FromCreateGRPC is
func FromCreateGRPC(request *order_v1.OrderCreateRequest) Request {
	return Request{
		ExpireTimeDuration: int(request.ExpireTimeDuration),
		OrderID:            request.Order.OrderID,
		ClientID:           request.Order.ClientID,
		Weight:             request.Order.Weight,
		BoxID:              request.Order.BoxID,
	}
}

// FromIssueGRPC is
func FromIssueGRPC(request *order_v1.IssueOrderRequest) RequestOrderIDs {
	requestOrderIDs := make([]int64, len(request.OrderIDRequest))
	for index, orderID := range request.OrderIDRequest {
		requestOrderIDs[index] = orderID.OrderID
	}
	return RequestOrderIDs{OrderIDs: requestOrderIDs}
}

// UniqueClientInfoToGRPC is
func UniqueClientInfoToGRPC(uniqueClientResponse []ListUniqueClients) []int64 {
	response := make([]int64, len(uniqueClientResponse))
	for index, value := range uniqueClientResponse {
		response[index] = value.ClientID
	}

	return response
}

// UniqueClientListToGRPC is
func UniqueClientListToGRPC(response abstractModel.PaginatedResponse[ListUniqueClients]) *order_v1.UniqueClientListResponse {
	return &order_v1.UniqueClientListResponse{
		ClientIDs: UniqueClientInfoToGRPC(response.Items),
		Pagination: abstractModel.PaginationToGRPC(
			abstractModel.Page{
				CurrentPage:  response.CurrentPage,
				ItemsPerPage: response.ItemsPerPage,
			},
			response.TotalItems,
		),
	}
}

// ReturnedInfoToGRPC is
func ReturnedInfoToGRPC(returnedResponse []ReturnedResponse) []*order_v1.ReturnedResponse {
	response := make([]*order_v1.ReturnedResponse, len(returnedResponse))
	for index, value := range returnedResponse {
		response[index] = &order_v1.ReturnedResponse{
			OrderID:    value.OrderID,
			ClientID:   value.ClientID,
			ReturnedAt: timestamppb.New(value.ReturnedAt),
		}
	}

	return response
}

// ReturnedListToGRPC is
func ReturnedListToGRPC(response abstractModel.PaginatedResponse[ReturnedResponse]) *order_v1.ReturnedListResponse {
	return &order_v1.ReturnedListResponse{
		ReturnedResponse: ReturnedInfoToGRPC(response.Items),
		Pagination: abstractModel.PaginationToGRPC(
			abstractModel.Page{
				CurrentPage:  response.CurrentPage,
				ItemsPerPage: response.ItemsPerPage,
			},
			response.TotalItems,
		),
	}
}

// AllInfoToGrpc is
func AllInfoToGrpc(allResponse []AllResponse) []*order_v1.OrderAllInfo {
	response := make([]*order_v1.OrderAllInfo, len(allResponse))
	for index, value := range allResponse {
		createdAt := value.CreatedAt
		updatedAt := value.UpdatedAt
		response[index] = &order_v1.OrderAllInfo{
			Order: &order_v1.Order{
				OrderID:  value.OrderID,
				ClientID: value.ClientID,
				Weight:   value.Weight,
				BoxID:    value.BoxID,
			},
			CreatedAt:  abstractModel.SafeTimestamp(&createdAt),
			UpdatedAt:  abstractModel.SafeTimestamp(&updatedAt),
			AcceptedAt: abstractModel.SafeTimestamp(value.AcceptedAt),
			IssuedAt:   abstractModel.SafeTimestamp(value.IssuedAt),
			ExpiresAt:  abstractModel.SafeTimestamp(value.ExpiresAt),
		}
	}

	return response
}

// ListToGRPC is
func ListToGRPC(response abstractModel.PaginatedResponse[AllResponse]) *order_v1.OrderListResponse {
	return &order_v1.OrderListResponse{
		OrderAllInfo: AllInfoToGrpc(response.Items),
		Pagination: abstractModel.PaginationToGRPC(
			abstractModel.Page{
				CurrentPage:  response.CurrentPage,
				ItemsPerPage: response.ItemsPerPage,
			},
			response.TotalItems,
		),
	}
}

// FromAcceptGRPC is
func FromAcceptGRPC(request *order_v1.RequestWithClientID) RequestWithClientID {
	return RequestWithClientID{
		OrderID:  request.OrderID,
		ClientID: request.ClientID,
	}
}
