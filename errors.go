package goagi

type agiError struct {
	s string
}

func errorNew(text string) *agiError {
	return &agiError{text}
}

func (e *agiError) withInfo(text string) *agiError {
	e.s = e.s + ": " + text
	return e
}

func (e *agiError) Error() string {
	return e.s
}
