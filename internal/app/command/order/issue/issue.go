package issue

import (
	"errors"
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"Homework-1/internal/app/command"
	"Homework-1/internal/model/cli"
	"Homework-1/internal/storage/order/file"
)

// Command is
type Command struct {
	command.Base
	store    *file.Storage
	orderIDs SliceInt
}

// SliceInt is
type SliceInt []int64

func (s *SliceInt) String() string {
	return fmt.Sprint(*s)
}

// Set is
func (s *SliceInt) Set(value string) error {
	if len(*s) > 0 {
		return errors.New("sliceint flag already set")
	}
	for _, ni := range strings.Split(value, ",") {
		num, err := strconv.ParseInt(ni, 10, 64)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt: %w", err)
		}
		*s = append(*s, num)
	}
	return nil
}

// New is
func New(store *file.Storage) *Command {
	cmd := Command{
		Base:  command.New("issue", "issue - finds the order of a client, and hands them to him, requires 1 flag", flag.NewFlagSet("flag", flag.ExitOnError)),
		store: store,
	}
	cmd.FlagSet.Var(&cmd.orderIDs, "slice", "provide slice of order' ID separated by comma(,) without spaces")
	return &cmd
}

// CheckFormat is
func (g *Command) CheckFormat(args []string, callHelp bool) (*command.Arguments, error) {
	if len(args) > 0 && strings.HasSuffix(args[0], "help") {
		args = make([]string, 0)
		callHelp = true
	}

	var orderIDs SliceInt
	issueFlagSet := flag.NewFlagSet("flag", flag.ExitOnError)
	issueFlagSet.Var(&orderIDs, "slice", "provide slice of order' ID separated by comma(,) without spaces")
	err := issueFlagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("FlagSet.Parse: %w", err)
	}

	if callHelp {
		return nil, nil
	}

	if len(orderIDs) == 0 {
		return nil, errors.New("flag slice (slice of order_id's) is not specified")
	}
	if len(args) > 2 {
		return nil, fmt.Errorf("too many arguments for: %s command", g.Name())
	}

	convertedOrderIDs := make([]cli.OrderID, 0, len(orderIDs))
	for _, val := range orderIDs {
		convertedOrderIDs = append(convertedOrderIDs, cli.OrderID(val))
	}

	return &command.Arguments{
		OrderIDs: convertedOrderIDs,
	}, nil
}

// Do is
func (g *Command) Do(args *command.Arguments) (*string, error) {
	activeOrders, err := g.store.List()
	if err != nil {
		return nil, fmt.Errorf("storage.List: %w", err)
	}

	var clientID *cli.UserID
	sort.Slice(activeOrders, func(i, j int) bool {
		return activeOrders[i].ID < activeOrders[j].ID
	})
	sort.Slice(args.OrderIDs, func(i, j int) bool {
		return args.OrderIDs[i] < args.OrderIDs[j]
	})

	index := 0
	for _, val := range activeOrders {
		if index < len(args.OrderIDs) && args.OrderIDs[index] == val.ID {
			switch {
			case clientID == nil:
				userIDCopy := val.UserID
				clientID = &userIDCopy
			case *clientID != val.UserID:
				return nil, errors.New("order belong to multiple client")
			case time.Now().After(val.ExpireAt):
				return nil, errors.New("some order are expired")
			case val.IsDeleted || val.IsIssued || val.IsAccepted || val.IsReturned:
				return nil, errors.New("some order weren't found")
			}
			index++
		}
		if len(args.OrderIDs) == index {
			err = g.store.IssueAll(args.OrderIDs)
			if err != nil {
				return nil, fmt.Errorf("storage.IssueAll: %w", err)
			}
			message := "All order were handed to client\n"
			return &message, nil
		}
	}
	return nil, errors.New("some order weren't found")
}
