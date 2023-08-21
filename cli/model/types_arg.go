package model

type ArgCollection struct {
	MandatoryArgs []*Arg
	OptionalArg   *Arg
}

type ArgIsValidFunc func() (bool, error)

type Arg struct {
	Name             string
	ShortDescription string
	Value            string
	IsValidFunc      ArgIsValidFunc
}
