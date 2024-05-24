package command

import (
	"bytes"
	"flag"
)

// Base is
type Base struct {
	name        string
	description string
	FlagSet     *flag.FlagSet
}

// New is
func New(name string, description string, FlagSet *flag.FlagSet) Base {
	return Base{
		name:        name,
		description: description,
		FlagSet:     FlagSet,
	}
}

// Name is
func (c *Base) Name() string {
	return c.name
}

// Description is
func (c *Base) Description() string {
	return c.description
}

// Help is
func (c *Base) Help() string {
	var buf bytes.Buffer
	c.FlagSet.SetOutput(&buf)
	c.FlagSet.Usage()
	return buf.String()
}
