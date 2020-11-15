package api

import "fmt"

type Error struct {
	Code        string
	DisplayText string
}

func NewError(code, displayText string) *Error {
	return &Error{Code: code, DisplayText: displayText}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.DisplayText)
}
