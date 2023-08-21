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
	GetDefaultValue() T
	GetValueRef() *T
	SetValueRef(*T)
	GetValueOrDefault() T
	GetIsValidFunc() FlagIsValidFunc[T]
}

type flag[T any] struct {
	name             string
	alias            *rune
	shortDescription string
	defaultValue     T
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

func (this *flag[T]) GetDefaultValue() T {
	return this.defaultValue
}

func (this *flag[T]) GetValueOrDefault() T {
	if this.valueRef != nil {
		return *this.valueRef
	}

	return this.defaultValue
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
