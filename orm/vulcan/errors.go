package vulcan

type InvalidStructError struct {
	message string
}

func (e *InvalidStructError) Error() string {
	return e.message
}
