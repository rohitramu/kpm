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

func GetCobraImplementation(rootCmd *model.Command) ExecuteRootCommand {
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
		var currentCobraCmd = convertToCobraCommand(currentCobraStackItem.toAdd)

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

func convertToCobraCommand(modelCmd *model.Command) *cobra.Command {
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
	addFlags(result, modelCmd)

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

	// Set behavior of command if there is any.
	if modelCmd.ExecuteFunc != nil {
		result.RunE = func(cmd *cobra.Command, args []string) error {
			// Set args.
			for i := 0; i < len(modelCmd.Args.MandatoryArgs); i++ {
				modelCmd.Args.MandatoryArgs[i].Value = args[i]
			}
			if modelCmd.Args.OptionalArg != nil && len(args) > len(modelCmd.Args.MandatoryArgs) {
				modelCmd.Args.OptionalArg.Value = args[len(modelCmd.Args.MandatoryArgs)]
			}

			return modelCmd.ExecuteFunc(modelCmd.Args)
		}
	}

	return result
}

func addFlags(cobraCmd *cobra.Command, modelCmd *model.Command) {
	// String flags.
	for _, modelFlag := range modelCmd.Flags.StringFlags {
		addFlag(cobraCmd.PersistentFlags().StringVarP, modelFlag)
	}

	// Bool flags.
	for _, modelFlag := range modelCmd.Flags.BoolFlags {
		addFlag(cobraCmd.PersistentFlags().BoolVarP, modelFlag)
	}
}

func addFlag[T any](
	addFlagPFunc func(p *T, name string, alias string, value T, usage string),
	modelFlag model.Flag[T],
) {
	var alias string
	if modelFlag.GetAlias() != nil {
		alias = string(*modelFlag.GetAlias())
	}

	// Cobra can't handle a nil pointer.
	if modelFlag.GetValueRef() == nil {
		var temp T
		modelFlag.SetValueRef(&temp)
	}

	addFlagPFunc(
		modelFlag.GetValueRef(),
		modelFlag.GetName(),
		alias,
		modelFlag.GetDefaultValue(),
		modelFlag.GetShortDescription(),
	)
}
