package model

import (
	"github.com/rohitramu/kpm/pkg/utils/log"
)

type FlagBuilder[T any] interface {
	SetAlias(rune) FlagBuilder[T]
	SetShortDescription(string) FlagBuilder[T]
	SetDefaultValue(T) FlagBuilder[T]
	SetValidationFunc(FlagIsValidFunc[T]) FlagBuilder[T]
	Build() Flag[T]
}

type flagBuilder[T any] struct {
	value flag[T]
}

func NewFlagBuilder[T any](flagName string) FlagBuilder[T] {
	if flagName == "" {
		log.Panicf("Failed to create a FlagBuilder: flag name cannot be empty.")
	}

	return &flagBuilder[T]{
		value: flag[T]{name: flagName},
	}
}

func (thisBuilder *flagBuilder[T]) SetAlias(alias rune) FlagBuilder[T] {
	thisBuilder.value.alias = &alias

	return thisBuilder
}

func (thisBuilder *flagBuilder[T]) SetShortDescription(shortDescription string) FlagBuilder[T] {
	thisBuilder.value.shortDescription = shortDescription

	return thisBuilder
}

func (thisBuilder *flagBuilder[T]) SetDefaultValue(value T) FlagBuilder[T] {
	thisBuilder.value.defaultValue = value

	// Make sure that the default value doesn't change when the value itself changes.
	var temp = value
	thisBuilder.value.valueRef = &temp

	return thisBuilder
}

func (thisBuilder *flagBuilder[T]) SetValidationFunc(validationFunc FlagIsValidFunc[T]) FlagBuilder[T] {
	thisBuilder.value.isValidFunc = validationFunc

	return thisBuilder
}

func (thisBuilder *flagBuilder[T]) Build() Flag[T] {
	return &thisBuilder.value
}
