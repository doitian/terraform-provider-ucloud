package client

type InvalidClientFieldError string

func (i InvalidClientFieldError) Error() string {
	return "Invalid client field: " + string(i)
}
