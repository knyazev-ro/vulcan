package vulcan

type InvalidStructError struct {
	message string
}

func (e *InvalidStructError) Error() string {
	return e.message
}

type FindByIdError struct {
	message string
}

func (e *FindByIdError) Error() string {
	return e.message
}
