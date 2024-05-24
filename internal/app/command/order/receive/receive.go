package receive

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
		Base:  command.New("receive", "receive - gets order from courier to PVZ, requires 3 flags", flag.NewFlagSet("flag", flag.ExitOnError)),
		store: store,
	}
	var orderID int64
	var userID int64
	var duration int
	cmd.FlagSet.Int64Var(&orderID, "orderid", 0, "provide order ID")
	cmd.FlagSet.Int64Var(&userID, "userid", 0, "provide user ID")
	cmd.FlagSet.IntVar(&duration, "expire", 0, "provide duration of expiration in days (after how much days order expires)")
	cmd.args = &command.Arguments{
		OrderID:  cli.OrderID(orderID),
		UserID:   cli.UserID(userID),
		Duration: time.Duration(duration) * 24 * time.Hour,
	}
	return &cmd
}

// CheckFormat is
func (r *Command) CheckFormat(args []string, callHelp bool) (*command.Arguments, error) {
	if len(args) > 0 && strings.HasSuffix(args[0], "help") {
		args = make([]string, 0)
		callHelp = true
	}

	var orderID int64
	var userID int64
	var duration int
	receiveFlagSet := flag.NewFlagSet("flag", flag.ExitOnError)
	receiveFlagSet.Int64Var(&orderID, "orderid", 0, "provide order ID")
	receiveFlagSet.Int64Var(&userID, "userid", 0, "provide user ID")
	receiveFlagSet.IntVar(&duration, "expire", 0, "provide duration of expiration in days (after how much days order expires)")
	err := receiveFlagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("FlagSet.Parse: %w", err)
	}

	if callHelp {
		return nil, nil
	}

	flags := []bool{orderID == 0, userID == 0, duration == 0}
	flagNames := []string{" userid", " orderid", " expire"}
	errorMessage := "flag(s)"
	comma := 0
	for i, val := range flags {
		if val {
			if comma > 0 {
				errorMessage += ","
			}
			errorMessage += flagNames[i]
			comma++
		}
	}

	switch {
	case comma > 0:
		errorMessage += " is (are) not specified"
		return nil, errors.New(errorMessage)
	case len(args) > 6:
		return nil, errors.New("too many arguments for: " + r.Name() + " command")
	case duration <= 0:
		return nil, errors.New("expiration date must me in future")
	}

	return &command.Arguments{
		OrderID:  cli.OrderID(orderID),
		UserID:   cli.UserID(userID),
		Duration: time.Duration(duration) * 24 * time.Hour,
	}, nil
}

// Do is
func (r *Command) Do(args *command.Arguments) (*string, error) {
	receiveAt := time.Now()
	expireAt := receiveAt.Add(args.Duration)

	err := r.store.Create(&cli.Order{
		ID:         args.OrderID,
		UserID:     args.UserID,
		ExpireAt:   expireAt,
		IsDeleted:  false,
		IsReturned: false,
		IsIssued:   false,
		IsAccepted: false,
		ReceivedAt: receiveAt,
		IssuedAt:   expireAt,
	})
	if err != nil {
		return nil, fmt.Errorf("storage.Create: %w", err)
	}

	message := fmt.Sprintf("Successfully received order: %d\n", args.OrderID)
	return &message, nil
}
