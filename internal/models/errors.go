package models

import "fmt"

type FormatError struct {
	Err error
}

type AlreadyExistError struct {
	Err error
}

type LoginPasswordError struct {
	Err error
}

type IncorrectInputError struct {
	Err error
}

func (fe *FormatError) Error() string {
	return fmt.Sprintf("Format error [%s]", fe.Err)
}

func (fe *AlreadyExistError) Error() string {
	return fmt.Sprintf("Login already exist [%s]", fe.Err)
}

func (fe *LoginPasswordError) Error() string {
	return fmt.Sprintf("Login/password pair error [%s]", fe.Err)
}

func (fe *IncorrectInputError) Error() string {
	return fmt.Sprintf("Incorrect input error [%s]", fe.Err)
}
