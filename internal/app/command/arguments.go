package command

import (
	"time"

	"Homework-1/internal/model/cli"
)

// Arguments is
type Arguments struct {
	OrderID        cli.OrderID
	UserID         cli.UserID
	Page           int
	PageSize       int
	Duration       time.Duration
	AmountOfOrders int
	OrderIDs       []cli.OrderID
	Name           string
	Address        string
	Contact        string
}
