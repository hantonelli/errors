package errors

type GenericErrorI interface {
	error
	IsGenericError() bool
}

type GenericError struct {}

func (g GenericError) Error() string {
	return "GenericError"
}

func (g GenericError) IsGenericError() bool {
	return true
}

func ContainsGenericError(err error) (GenericErrorI, map[string]interface{}, bool) {
	ce, isExpectedType := err.(GenericErrorI)
	if isExpectedType {
		return ce, map[string]interface{}{}, true
	}
	we, isWrappedError := err.(WrappedError)
	if !isWrappedError {
		return nil, map[string]interface{}{}, false
	}
	ge1, isActualExpectedType := we.GetActual().(GenericError)
	if isActualExpectedType {
		return ge1, we.GetFields(), true
	}
	if we.GetPrevious() != nil {
		ge2, isPreviousExpectedType := we.GetPrevious().(GenericError)
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
