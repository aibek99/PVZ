package get

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"Homework-1/internal/app/command"
	"Homework-1/internal/storage/pvz"
)

// Command is
type Command struct {
	command.Base
	store pvz.Storage
	args  *command.Arguments
}

// New is
func New(store pvz.Storage) *Command {
	cmd := Command{
		Base:  command.New("get_pvz", "get_pvz - gets a pvz from database, requires 1 flag", flag.NewFlagSet("flag", flag.ExitOnError)),
		store: store,
	}
	var name string
	cmd.FlagSet.StringVar(&name, "name", "", "provide the name of the PVZ")
	cmd.args = &command.Arguments{
		Name: name,
	}
	return &cmd
}

// CheckFormat is
func (a *Command) CheckFormat(args []string, callHelp bool) (*command.Arguments, error) {
	if len(args) > 0 && strings.HasSuffix(args[0], "help") {
		args = make([]string, 0)
		callHelp = true
	}

	var pvzName string

	getPvzFlagSet := flag.NewFlagSet("flag", flag.ExitOnError)
	getPvzFlagSet.StringVar(&pvzName, "name", "", "provide the name of the PVZ")

	err := getPvzFlagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("FlagSet.Parse: %w", err)
	}

	if callHelp {
		return nil, nil
	}

	if pvzName == "" {
		return nil, errors.New("flag name wasn't provided")
	}
	return &command.Arguments{
		Name: pvzName,
	}, nil
}

// Do is
func (a *Command) Do(args *command.Arguments) (*string, error) {
	if args == nil {
		return nil, nil
	}
	pvzModel, err := a.store.Find(args.Name)
	if err != nil {
		return nil, fmt.Errorf("store.Find: %w", err)
	}
	message := fmt.Sprintf("\nPVZ Name: %s\nPVZ address: %s\nPVZ Contact: %s\n", pvzModel.Name, pvzModel.Address, pvzModel.Contact)
	return &message, nil
}
