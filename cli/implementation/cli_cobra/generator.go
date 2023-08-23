package cli_cobra

import (
	"fmt"
	"strings"

	"github.com/emirpasic/gods/stacks/linkedliststack"
	"github.com/spf13/cobra"

	"github.com/rohitramu/kpm/cli/model"
	"github.com/rohitramu/kpm/pkg/utils/log"
)

type ExecuteRootCommand func() error

func GetCobraImplementation(config *model.KpmConfig, rootCmd *model.Command) ExecuteRootCommand {
	var result ExecuteRootCommand

	var toVisit = linkedliststack.New()
	toVisit.Push(&cobraStackItem{toAdd: rootCmd})

	for stackItem, ok := toVisit.Pop(); ok; stackItem, ok = toVisit.Pop() {
		// Get next command.
		currentCobraStackItem, castOk := stackItem.(*cobraStackItem)
		if !castOk {
			log.Panicf("Failed to cast stack item to command")
		}

		// Convert to a Cobra command.
		var currentCobraCmd = convertToCobraCommand(config, currentCobraStackItem.toAdd)

		// Use the root command's "Execute()" method as the result.
		if currentCobraStackItem.parent == nil {
			result = currentCobraCmd.Execute
		} else {
			// Add this command as a subcommand on the parent.
			currentCobraStackItem.parent.AddCommand(currentCobraCmd)
		}

		// Mark subcommands as "toVisit".
		for _, subCmd := range currentCobraStackItem.toAdd.SubCommands {
			toVisit.Push(&cobraStackItem{
				parent: currentCobraCmd,
				toAdd:  subCmd,
			})
		}
	}

	return result
}

type cobraStackItem struct {
	parent *cobra.Command
	toAdd  *model.Command
}

func convertToCobraCommand(config *model.KpmConfig, modelCmd *model.Command) *cobra.Command {
	// Create result object and set basic info.
	var result = &cobra.Command{
		Use:   modelCmd.Name,
		Short: modelCmd.ShortDescription,
	}

	// Add aliases.
	if modelCmd.Alias != "" {
		result.Aliases = []string{modelCmd.Alias}
	}

	// Add flags.
	addFlags(result, modelCmd, config)

	// Args validation.
	var numMandatoryArgs = len(modelCmd.Args.MandatoryArgs)
	if modelCmd.Args.OptionalArg != nil {
		result.Args = cobra.RangeArgs(numMandatoryArgs, numMandatoryArgs+1)
	} else {
		result.Args = cobra.ExactArgs(numMandatoryArgs)
	}

	// Update "Use" string with args.
	var argsString = strings.Builder{}
	if len(modelCmd.Args.MandatoryArgs) > 0 {
		for _, arg := range modelCmd.Args.MandatoryArgs {
			argsString.WriteString(fmt.Sprintf(" <%s>", arg.Name))
		}
	}
	if modelCmd.Args.OptionalArg != nil {
		var arg = modelCmd.Args.OptionalArg
		argsString.WriteString(fmt.Sprintf(" [<%s>]", arg.Name))
	}
	result.Use = fmt.Sprintf("%s%s", result.Use, argsString.String())

	// Set pre-execute validation if there is any.
	if modelCmd.IsValidFunc != nil {
		result.PreRunE = func(cmd *cobra.Command, args []string) error {
			// Set args.
			setArgs(&modelCmd.Args, args)

			return modelCmd.IsValidFunc(config, modelCmd.Args)
		}
	}

	// Set behavior of command if there is any.
	if modelCmd.ExecuteFunc != nil {
		result.RunE = func(cmd *cobra.Command, args []string) error {
			// Set args.
			setArgs(&modelCmd.Args, args)

			return modelCmd.ExecuteFunc(config, modelCmd.Args)
		}
	}

	return result
}

func setArgs(argCollection *model.ArgCollection, args []string) {
	for i := 0; i < len(argCollection.MandatoryArgs); i++ {
		argCollection.MandatoryArgs[i].Value = args[i]
	}
	if argCollection.OptionalArg != nil && len(args) > len(argCollection.MandatoryArgs) {
		argCollection.OptionalArg.Value = args[len(argCollection.MandatoryArgs)]
	}
}

func addFlags(cobraCmd *cobra.Command, modelCmd *model.Command, config *model.KpmConfig) {
	// String flags.
	for _, modelFlag := range modelCmd.Flags.StringFlags {
		addFlag(cobraCmd.PersistentFlags().StringVarP, modelFlag, config)
	}

	// Bool flags.
	for _, modelFlag := range modelCmd.Flags.BoolFlags {
		addFlag(cobraCmd.PersistentFlags().BoolVarP, modelFlag, config)
	}
}

func addFlag[T any](
	addFlagPFunc func(p *T, name string, alias string, value T, usage string),
	modelFlag model.Flag[T],
	config *model.KpmConfig,
) {
	// Get alias.
	var alias string
	if modelFlag.GetAlias() != nil {
		alias = string(*modelFlag.GetAlias())
	}

	// Get default value.
	var defaultValue T
	defaultValue = modelFlag.GetDefaultValue(config)

	// Cobra can't handle a nil pointer.
	if modelFlag.GetValueRef() == nil {
		var temp = defaultValue
		modelFlag.SetValueRef(&temp)
	}

	addFlagPFunc(
		modelFlag.GetValueRef(),
		modelFlag.GetName(),
		alias,
		defaultValue,
		modelFlag.GetShortDescription(),
	)
}
