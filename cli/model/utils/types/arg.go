package types

type ArgCollection struct {
	MandatoryArgs []*Arg
	OptionalArg   *Arg
}

type ArgIsValidFunc func(value string) error

type Arg struct {
	Name             string
	ShortDescription string
	Value            string
	IsValidFunc      ArgIsValidFunc
}
