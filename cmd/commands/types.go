package commands

import (
	"github.com/rohitramu/kpm/cmd/flags"
	"github.com/spf13/cobra"
)

type kpmCommandBuilder struct {
	state *cobra.Command
}

func newKpmCommandBuilder(cobraCommand *cobra.Command) *kpmCommandBuilder {
	return &kpmCommandBuilder{
		state: cobraCommand,
	}
}

func (cmdBuilder *kpmCommandBuilder) Build() *cobra.Command {
	return cmdBuilder.state
}

func (cmdBuilder *kpmCommandBuilder) AddLocalFlags(kpmFlags ...flags.KpmFlag) *kpmCommandBuilder {
	for _, flag := range kpmFlags {
		flag.Add(cmdBuilder.state.Flags())
	}

	return cmdBuilder
}

func (cmdBuilder *kpmCommandBuilder) AddPersistentFlags(kpmFlags ...flags.KpmFlag) *kpmCommandBuilder {
	for _, flag := range kpmFlags {
		flag.Add(cmdBuilder.state.PersistentFlags())
	}

	return cmdBuilder
}

func (cmdBuilder *kpmCommandBuilder) AddSubcommands(subcommands ...*cobra.Command) *kpmCommandBuilder {
	cmdBuilder.state.AddCommand(subcommands...)

	return cmdBuilder
}

func (cmdBuilder *kpmCommandBuilder) AddCommandGroups(group ...*cobra.Group) *kpmCommandBuilder {
	cmdBuilder.state.AddGroup(group...)

	return cmdBuilder
}
