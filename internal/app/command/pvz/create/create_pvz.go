package create

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"Homework-1/internal/app/command"
	"Homework-1/internal/model/cli"
	"Homework-1/internal/storage/pvz"
	"Homework-1/pkg/constants"
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
		Base:  command.New("create_pvz", "create_pvz - creates a pvz in database, requires 2 flags", flag.NewFlagSet("flag", flag.ExitOnError)),
		store: store,
	}
	var name string
	var address string
	var contact string
	cmd.FlagSet.StringVar(&name, "name", "", "provide the name of the PVZ")
	cmd.FlagSet.StringVar(&address, "address", "", "provide the name of the PVZ")
	cmd.FlagSet.StringVar(&contact, "contact", "", "provide the name of the PVZ")
	cmd.args = &command.Arguments{
		Name:    name,
		Address: address,
		Contact: contact,
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
	var address string
	var contact string

	createPvzFlagSet := flag.NewFlagSet("flag", flag.ExitOnError)
	createPvzFlagSet.StringVar(&pvzName, "name", "", "provide the name of the PVZ")
	createPvzFlagSet.StringVar(&address, "address", "", "provide the name of the PVZ")
	createPvzFlagSet.StringVar(&contact, "contact", "", "provide the name of the PVZ")

	err := createPvzFlagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("FlagSet.Parse: %w", err)
	}

	if callHelp {
		return nil, nil
	}

	if pvzName == "" || address == "" || contact == "" {
		return nil, errors.New("flags shouldn't be empty, there exists 3 different flags")
	}
	return &command.Arguments{
		Name:    pvzName,
		Address: address,
		Contact: contact,
	}, nil
}

// Do is
func (a *Command) Do(args *command.Arguments) (*string, error) {
	if args == nil {
		return nil, nil
	}
	newPVZ := cli.PVZ{
		Name:    args.Name,
		Address: args.Address,
		Contact: args.Contact,
	}
	err := a.store.Create(&newPVZ)
	if err != nil {
		return nil, fmt.Errorf("store.Create: %w", err)
	}
	message := constants.SuccessfullyCreatedPVZ
	return &message, nil
}
