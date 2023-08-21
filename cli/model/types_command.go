package model

type CommandExecuteFunc func(args ArgCollection) error
type CommandIsValidFunc func(args ArgCollection) (bool, error)

type Command struct {
	Name             string
	Alias            string
	ShortDescription string
	SubCommands      []*Command
	Flags            FlagCollection
	Args             ArgCollection
	IsValidFunc      CommandIsValidFunc
	ExecuteFunc      CommandExecuteFunc
}
