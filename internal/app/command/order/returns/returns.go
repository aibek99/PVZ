package returns

import (
	"flag"
	"fmt"
	"strings"

	"Homework-1/internal/app/command"
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
		Base:  command.New("returns", "returns - returns the list of returned order, requires 2 flags (optional)", flag.NewFlagSet("flag", flag.ExitOnError)),
		store: store,
	}
	var page int
	var pageSize int
	cmd.FlagSet.IntVar(&page, "page", 1, "indicates the page number")
	cmd.FlagSet.IntVar(&pageSize, "pagesize", 5, "indicates the size of page")
	cmd.args = &command.Arguments{
		Page:     page,
		PageSize: pageSize,
	}
	return &cmd
}

// CheckFormat is
func (r *Command) CheckFormat(args []string, callHelp bool) (*command.Arguments, error) {
	if len(args) > 0 && strings.HasSuffix(args[0], "help") {
		args = make([]string, 0)
		callHelp = true
	}

	var page int
	var pageSize int
	returnsFlagSet := flag.NewFlagSet("flag", flag.ExitOnError)
	returnsFlagSet.IntVar(&page, "page", 1, "indicates the page number")
	returnsFlagSet.IntVar(&pageSize, "pagesize", 5, "indicates the size of page")

	err := returnsFlagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("FlagSet.Parse: %w", err)
	}

	if callHelp {
		return nil, nil
	}

	return &command.Arguments{
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// Do is
func (r *Command) Do(args *command.Arguments) (*string, error) {
	listOfReturns := r.store.ListReturns(args.PageSize)
	message := fmt.Sprintf("Page: %d,  Page size: %d\n", args.Page, args.PageSize)
	for index, val := range listOfReturns {
		message = message + fmt.Sprintf("%d. Order ID: %d;   User ID: %d\n", index+1, val.ID, val.UserID)
	}
	return &message, nil
}
