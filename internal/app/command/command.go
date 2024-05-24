package command

// Commander is
type Commander interface {
	Name() string
	Description() string
	CheckFormat(args []string, callHelp bool) (*Arguments, error)
	Do(args *Arguments) (*string, error)
	Help() string
}
