package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/thediveo/enumflag/v2"
	"golang.org/x/exp/constraints"
)

type KpmFlagValue[T any] interface {
	GetValue() T
	SetValue(T) error
	GetValueRef() *T
}

type KpmFlag interface {
	Add(*pflag.FlagSet)
	IsSetByUser(*cobra.Command) bool
	GetFlagName() string
	GetShortDescription() string
}

type KpmEnumFlag interface {
	GetPFlagValue() pflag.Value
}

type kpmFlagBase[T any] struct {
	KpmFlag
	KpmFlagValue[T]

	flagName         string
	shortDescription string
	value            T
}

type kpmEnumFlagBase[T constraints.Integer] struct {
	KpmEnumFlag
	*kpmFlagBase[T]

	enumFlag pflag.Value
}

func newKpmFlagBase[T any](flagName string, shortDescription string, value T) *kpmFlagBase[T] {
	var result = &kpmFlagBase[T]{
		flagName:         flagName,
		shortDescription: shortDescription,
		value:            value,
	}

	result.KpmFlag = result
	result.KpmFlagValue = result

	return result
}

func newKpmEnumFlagBase[T constraints.Integer](flagName string, shortDescription string, value T, enumMapping map[T][]string) *kpmEnumFlagBase[T] {
	var valueInternal = value
	var enumFlag = enumflag.New(
		&value,
		"logLevel",
		enumMapping,
		enumflag.EnumCaseInsensitive)
	var base = newKpmFlagBase[T](flagName, shortDescription, valueInternal)

	var result = &kpmEnumFlagBase[T]{
		kpmFlagBase: base,
		enumFlag:    enumFlag,
	}

	return result
}

func (instance *kpmFlagBase[T]) GetFlagName() string {
	return instance.flagName
}

func (instance *kpmFlagBase[T]) GetShortDescription() string {
	return instance.shortDescription
}

func (instance *kpmFlagBase[T]) GetValue() T {
	return instance.value
}

func (instance *kpmFlagBase[T]) GetValueRef() *T {
	return &instance.value
}

func (instance *kpmFlagBase[T]) SetValue(val T) error {
	instance.value = val

	return nil
}

func (instance *kpmFlagBase[T]) IsSetByUser(cmd *cobra.Command) bool {
	var flag = cmd.Flags().Lookup(instance.GetFlagName())
	return flag.Changed
}

func (instance *kpmEnumFlagBase[T]) GetPFlagValue() pflag.Value {
	return instance.enumFlag
}
