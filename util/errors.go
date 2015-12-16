package util

// WrapErr packs error with message
func WrapErr(message string, innerError error) *WrappedError {
	return &WrappedError{message, innerError}
}

// WrappedError is a pack of message and inner-error
type WrappedError struct {
	Message    string
	InnerError error
}

func (w *WrappedError) Error() string {
	return w.Message + " (inner:" + w.InnerError.Error() + ")"
}

// UnwrapError unpacks WrappedError
func UnwrapError(err error) error {
	if wErr, ok := err.(*WrappedError); ok {
		return UnwrapError(wErr.InnerError)
	}
	return err
}
