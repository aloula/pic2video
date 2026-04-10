package renderjob

import "fmt"

type ErrorClass int

const (
ErrInvalidArguments ErrorClass = iota + 1
ErrInputValidation
ErrEnvironment
ErrExecution
)

type ClassifiedError struct {
	Class ErrorClass
	Msg   string
	Err   error
}

func (e *ClassifiedError) Error() string {
	if e.Err == nil {
		return e.Msg
	}
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

func ExitCode(err error) int {
	ce, ok := err.(*ClassifiedError)
	if !ok {
		return 1
	}
	switch ce.Class {
	case ErrInvalidArguments:
		return 2
	case ErrInputValidation:
		return 3
	case ErrEnvironment:
		return 4
	case ErrExecution:
		return 5
	default:
		return 1
	}
}
