package cli

import "time"

// OrderID is
type OrderID int64

// UserID is
type UserID int64

// Order is
type Order struct {
	ID         OrderID   `json:"id"`
	UserID     UserID    `json:"user_id"`
	ExpireAt   time.Time `json:"expire_at"`
	IsDeleted  bool      `json:"is_deleted"`
	IsReturned bool      `json:"is_returned"`
	IsIssued   bool      `json:"is_issued"`
	IsAccepted bool      `json:"is_accepted"`
	ReceivedAt time.Time `json:"receive_at"`
	IssuedAt   time.Time `json:"issue_at"`
}

// PVZ is
type PVZ struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}
