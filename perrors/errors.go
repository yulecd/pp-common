package perrors

type CustomError interface {
	error
	Code() int
}

const CodeInvalidParams = 400

type ValidateError struct {
	err error
}

func NewValidateError(err error) CustomError {
	return &ValidateError{
		err: err,
	}
}

func (e *ValidateError) Is(err error) bool {
	if _, ok := err.(*ValidateError); ok {
		return true
	}

	return false
}

func (e *ValidateError) Error() string {
	return e.err.Error()
}

func (e *ValidateError) Code() int {
	return CodeInvalidParams
}
