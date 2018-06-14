package errors

// GenericError is an interface of the error to use as an example of how to write a Contains function.
type GenericError interface {
	error
	IsGenericError() bool
}

type genericError struct{}

func (g genericError) Error() string {
	return "GenericError"
}

// IsGenericError returns always true and is used to identify an error type
func (g genericError) IsGenericError() bool {
	return true
}

// ContainsGenericError takes an error and returns a concrete error if it is present in the error chain.
func ContainsGenericError(err error) (GenericError, map[string]interface{}, bool) {
	ce, isExpectedType := err.(GenericError)
	if isExpectedType {
		return ce, map[string]interface{}{}, true
	}
	we, isWrappedError := err.(WrappedError)
	if !isWrappedError {
		return nil, map[string]interface{}{}, false
	}
	ge1, isActualExpectedType := we.GetActual().(genericError)
	if isActualExpectedType {
		return ge1, we.GetFields(), true
	}
	if we.GetPrevious() != nil {
		ge2, isPreviousExpectedType := we.GetPrevious().(genericError)
		if isPreviousExpectedType {
			return ge2, we.GetFields(), true
		}
		prev, isPreviousWrappedError := we.GetPrevious().(WrappedError)
		if isPreviousWrappedError {
			return ContainsGenericError(prev)
		}
	}
	return nil, map[string]interface{}{}, false
}
