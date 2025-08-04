package models

import "fmt"

type FormatError struct {
	Err error
}

func (fe *FormatError) Error() string {
	return fmt.Sprintf("Format error [%s]", fe.Err)
}
