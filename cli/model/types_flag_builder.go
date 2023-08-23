package model

import (
	"github.com/rohitramu/kpm/pkg/utils/log"
)

type DefaultValueFunc[T any] func(*KpmConfig) T

type FlagBuilder[T any] interface {
	SetAlias(rune) FlagBuilder[T]
	SetShortDescription(string) FlagBuilder[T]
	SetDefaultValueFunc(DefaultValueFunc[T]) FlagBuilder[T]
	SetValidationFunc(FlagIsValidFunc[T]) FlagBuilder[T]
	Build() Flag[T]
}

var _ FlagBuilder[any] = &flagBuilder[any]{}

type flagBuilder[T any] struct {
	value flag[T]
}

func NewFlagBuilder[T any](flagName string) FlagBuilder[T] {
	if flagName == "" {
		log.Panicf("Failed to create a FlagBuilder: flag name cannot be empty.")
	}

	return &flagBuilder[T]{
		value: flag[T]{
			name: flagName,
			// Set the defaultValueFunc so we don't get nil reference errors.
			defaultValueFunc: func(kc *KpmConfig) T {
				// Return the zero value for the type.
				var result T
				return result
			},
		},
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

func (thisBuilder *flagBuilder[T]) SetDefaultValueFunc(defaultValueFunc DefaultValueFunc[T]) FlagBuilder[T] {
	thisBuilder.value.defaultValueFunc = defaultValueFunc

	return thisBuilder
}

func (thisBuilder *flagBuilder[T]) SetValidationFunc(validationFunc FlagIsValidFunc[T]) FlagBuilder[T] {
	thisBuilder.value.isValidFunc = validationFunc

	return thisBuilder
}

func (thisBuilder *flagBuilder[T]) Build() Flag[T] {
	return &thisBuilder.value
}
