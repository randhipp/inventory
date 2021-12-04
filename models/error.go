package models

import "fmt"

type (
	Error struct {
		Message string     `json:"message"`
		Field   string     `json:"field"`
		Type    *ErrorType `json:"type"`
	}

	ErrorType string
)

func (e Error) Error() string {
	return fmt.Sprintf("%+s (%s)", e.Message, e.Field)
}
