package accept

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

const defaultTimeToReturn = 48 * time.Hour

// Command is
type Command struct {
	command.Base
	store *file.Storage
	args  *command.Arguments
}

// New is
func New(store *file.Storage) *Command {
	cmd := Command{
		Base:  command.New("accept", "accept - accepts returned order from a client, requires 2 flags", flag.NewFlagSet("flag", flag.ExitOnError)),
		store: store,
	}
	var userID int64
	var orderID int64
	cmd.FlagSet.Int64Var(&userID, "userid", 0, "provide user ID")
	cmd.FlagSet.Int64Var(&orderID, "orderid", 0, "provide order ID")
	cmd.args = &command.Arguments{
		OrderID: cli.OrderID(orderID),
		UserID:  cli.UserID(userID),
	}
	return &cmd
}

// CheckFormat is
func (a *Command) CheckFormat(args []string, callHelp bool) (*command.Arguments, error) {
	if len(args) > 0 && strings.HasSuffix(args[0], "help") {
		args = make([]string, 0)
		callHelp = true
	}

	var userID int64
	var orderID int64
	acceptFlagSet := flag.NewFlagSet("flag", flag.ExitOnError)
	acceptFlagSet.Int64Var(&userID, "userid", 0, "provide user ID")
	acceptFlagSet.Int64Var(&orderID, "orderid", 0, "provide order ID")
	err := acceptFlagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("FlagSet.Parse: %w", err)
	}
	if callHelp {
		return nil, nil
	}

	switch {
	case userID == 0 && orderID == 0:
		return nil, errors.New("flags userid and orderid are not specified")
	case userID == 0:
		return nil, errors.New("flag userid is not specified")
	case orderID == 0:
		return nil, errors.New("flag orderid is not specified")
	case len(args) > 4:
		return nil, errors.New("too many arguments for: " + a.Name() + " command")
	}

	return &command.Arguments{
		OrderID: cli.OrderID(orderID),
		UserID:  cli.UserID(userID),
	}, nil
}

// Do is
func (a *Command) Do(args *command.Arguments) (*string, error) {
	order, err := a.store.Find(args.OrderID)
	if err != nil {
		return nil, fmt.Errorf("storage.Find: %w", err)
	}

	if args.UserID != order.UserID {
		return nil, errors.New("user not recognized as owner of this order")
	}
	if !(order.IsIssued && order.IsDeleted) || order.IsReturned {
		return nil, errors.New("order was not found in database")
	}

	t := time.Now().Add(-defaultTimeToReturn)
	if order.IssuedAt.Before(t) {
		return nil, errors.New("can't be accepted back, two days passed already")
	}

	err = a.store.Update(cli.Order{
		ID:         args.OrderID,
		UserID:     args.UserID,
		ExpireAt:   order.ExpireAt,
		IsDeleted:  false,
		IsReturned: false,
		IsIssued:   false,
		IsAccepted: true,
		ReceivedAt: order.ReceivedAt,
		IssuedAt:   order.ExpireAt,
	})

	if err != nil {
		return nil, fmt.Errorf("storage.Update: %w", err)
	}
	message := "Order was successfully accepted\n"
	return &message, nil
}
