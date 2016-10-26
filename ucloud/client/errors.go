package client

import (
	"fmt"
)

type InvalidClientFieldError string

func (icfe InvalidClientFieldError) Error() string {
	return "Invalid client field: " + string(icfe)
}

type BadRetCodeError struct {
	Action  string
	RetCode int
	Message string
}

func (brce *BadRetCodeError) Error() string {
	return fmt.Sprintf("Bad RetCode %d:%s %s", brce.RetCode, brce.Action, brce.Message)
}
