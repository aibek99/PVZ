package turnin

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"Homework-1/internal/app/command"
	"Homework-1/internal/model/cli"
	"Homework-1/internal/storage/order/file"
)

// Command is
type Command struct {
	command.Base
	store *file.Storage
	args  *command.Arguments
}

// New is
func New(store *file.Storage) *Command {
	cmd := Command{
		Base:  command.New("turnin", "turnin - returns the order to courier, requires 1 flag", flag.NewFlagSet("flag", flag.ExitOnError)),
		store: store,
	}
	var orderID int64
	cmd.FlagSet.Int64Var(&orderID, "orderid", 0, "provide the order ID")
	cmd.args = &command.Arguments{
		OrderID: cli.OrderID(orderID),
	}
	return &cmd
}

// CheckFormat is
func (t *Command) CheckFormat(args []string, callHelp bool) (*command.Arguments, error) {
	if len(args) > 0 && strings.HasSuffix(args[0], "help") {
		args = make([]string, 0)
		callHelp = true
	}

	var orderID int64
	turnInFlagSet := flag.NewFlagSet("flag", flag.ExitOnError)
	turnInFlagSet.Int64Var(&orderID, "orderid", 0, "provide the order ID")
	err := turnInFlagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("turnin.CheckFormat: %w", err)
	}

	if callHelp {
		return nil, nil
	}

	if orderID == 0 {
		return nil, errors.New("flag orderid is not specified")
	}
	if len(args) > 2 {
		return nil, errors.New("too many arguments for: " + t.Name() + " command")
	}

	return &command.Arguments{
		OrderID: cli.OrderID(orderID),
	}, nil
}

// Do is
func (t *Command) Do(args *command.Arguments) (*string, error) {
	order, err := t.store.Find(args.OrderID)
	if err != nil {
		return nil, fmt.Errorf("storage.Find: %w", err)
	}

	if order.IsDeleted {
		return nil, errors.New("the requested order doesn't exist")
	}

	if order.IsAccepted || time.Now().After(order.ExpireAt) {
		err = t.store.Update(cli.Order{
			ID:         order.ID,
			UserID:     order.UserID,
			ExpireAt:   order.ExpireAt,
			IsDeleted:  true,
			IsReturned: true,
			IsIssued:   false,
			IsAccepted: false,
			ReceivedAt: order.ReceivedAt,
			IssuedAt:   time.Now(),
		})
		if err != nil {
			return nil, fmt.Errorf("storage.Update: %w", err)
		}
		message := "Successfully returned order to courier\n"
		return &message, nil
	}
	return nil, errors.New("the requested order can't be returned to courier")
}
