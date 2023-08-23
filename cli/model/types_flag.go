package model

type FlagCollection struct {
	StringFlags []Flag[string]
	BoolFlags   []Flag[bool]
}

type FlagIsValidFunc[T any] func(flagName string, flagValueRef *T) error

type Flag[T any] interface {
	GetName() string
	GetAlias() *rune
	GetShortDescription() string
	GetDefaultValue(*KpmConfig) T
	GetValueRef() *T
	SetValueRef(*T)
	GetValueOrDefault(*KpmConfig) T
	GetIsValidFunc() FlagIsValidFunc[T]
}

var _ Flag[any] = &flag[any]{}

type flag[T any] struct {
	name             string
	alias            *rune
	shortDescription string
	defaultValueFunc DefaultValueFunc[T]
	valueRef         *T
	isValidFunc      FlagIsValidFunc[T]
}

func (this *flag[T]) GetName() string {
	return this.name
}

func (this *flag[T]) GetAlias() *rune {
	return this.alias
}

func (this *flag[T]) GetShortDescription() string {
	return this.shortDescription
}

func (this *flag[T]) GetDefaultValue(config *KpmConfig) T {
	return this.defaultValueFunc(config)
}

func (this *flag[T]) GetValueOrDefault(config *KpmConfig) T {
	if this.valueRef != nil {
		return *this.valueRef
	}

	return this.GetDefaultValue(config)
}

func (this *flag[T]) GetValueRef() *T {
	return this.valueRef
}

func (this *flag[T]) SetValueRef(valueRef *T) {
	this.valueRef = valueRef
}

func (this *flag[T]) GetIsValidFunc() FlagIsValidFunc[T] {
	return this.isValidFunc
}
