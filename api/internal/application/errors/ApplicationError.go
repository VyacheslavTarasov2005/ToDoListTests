package errors

type ApplicationError struct {
	StatusCode int
	Code       string
	Errors     map[string]string
}

func (e ApplicationError) Error() string {
	return e.Code
}
