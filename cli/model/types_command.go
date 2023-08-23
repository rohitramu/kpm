package model

type CommandExecuteFunc func(config *KpmConfig, args ArgCollection) error
type CommandIsValidFunc func(config *KpmConfig, args ArgCollection) error

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
