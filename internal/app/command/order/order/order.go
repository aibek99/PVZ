package order

import (
	"errors"
	"flag"
	"fmt"
	"sort"
	"strings"

	"Homework-1/internal/app/command"
	"Homework-1/internal/model/cli"
	"Homework-1/internal/storage/order/file"
)

// TooManyArguments is
const TooManyArguments = "too many arguments for: "

// Command is
type Command struct {
	command.Base
	store *file.Storage
	args  *command.Arguments
}

// New is
func New(store *file.Storage) *Command {
	cmd := Command{
		Base:  command.New("order", "order - provides list of order of a user, command has 2 flags, 1 is mandatory, 1 is optional", flag.NewFlagSet("flag", flag.ExitOnError)),
		store: store,
	}
	var userID int64
	var amountOfOrders int
	cmd.FlagSet.Int64Var(&userID, "userid", 0, "provide user ID")
	cmd.FlagSet.IntVar(&amountOfOrders, "amount", 0, "optional flag to choose the last amount of order")
	cmd.args = &command.Arguments{
		UserID:         cli.UserID(userID),
		AmountOfOrders: amountOfOrders,
	}
	return &cmd
}

// CheckFormat is
func (o *Command) CheckFormat(args []string, callHelp bool) (*command.Arguments, error) {
	if len(args) > 0 && strings.HasSuffix(args[0], "help") {
		args = make([]string, 0)
		callHelp = true
	}

	var userID int64
	var amountOfOrders int
	ordersFlagSet := flag.NewFlagSet("flag", flag.ExitOnError)
	ordersFlagSet.Int64Var(&userID, "userid", 0, "provide user ID")
	ordersFlagSet.IntVar(&amountOfOrders, "amount", 0, "optional flag to choose the last amount of order")
	err := ordersFlagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("FlatSet.Parse: %w", err)
	}
	if callHelp {
		return nil, nil
	}

	switch {
	case userID == 0:
		return nil, errors.New("flag userid is not specified")
	case len(args) > 2 && amountOfOrders == 0:
		return nil, errors.New(TooManyArguments + o.Name() + " command, without flag -amount")
	case amountOfOrders != 0 && len(args) > 4:
		return nil, errors.New(TooManyArguments + o.Name() + " command")
	}

	return &command.Arguments{
		UserID:         cli.UserID(userID),
		AmountOfOrders: amountOfOrders,
	}, nil
}

// Do is
func (o *Command) Do(args *command.Arguments) (*string, error) {
	activeOrders, err := o.store.List()
	if err != nil {
		return nil, fmt.Errorf("storage.List: %w", err)
	}

	sort.Slice(activeOrders, func(i, j int) bool {
		return activeOrders[i].ReceivedAt.After(activeOrders[j].ReceivedAt)
	})

	if args.AmountOfOrders == 0 {
		args.AmountOfOrders = len(activeOrders)
	}

	foundOrders := make([]cli.Order, 0, args.AmountOfOrders)
	for _, val := range activeOrders {
		if val.UserID == args.UserID {
			foundOrders = append(foundOrders, val)
		}
		if len(foundOrders) == args.AmountOfOrders {
			break
		}
	}

	message := "Found order: \n"
	for _, val := range foundOrders {
		message = message + fmt.Sprintf("%d ", val.ID)
	}
	message = message + "\n"
	return &message, nil
}
