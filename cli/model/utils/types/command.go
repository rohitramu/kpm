package types

import (
	"github.com/rohitramu/kpm/cli/model/utils/config"
)

type CommandExecuteFunc func(config *config.KpmConfig, args ArgCollection) error
type CommandIsValidFunc func(config *config.KpmConfig, args ArgCollection) error

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
