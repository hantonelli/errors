package errors

import (
	"errors"
	"strings"
)

func cleanSpaces(txt string) string {
	newTxt := strings.Replace(txt, "  ", " ", -1)
	newTxt = strings.Replace(newTxt, "\n", "", -1)
	newTxt = strings.Replace(newTxt, "\t", "", -1)
	return newTxt
}

func getThirdWrap() error {
	secondWrap := getSecondWrap()
	errThirdWrap := errors.New("third wrap")
	return WithErrorAndFields(secondWrap, errThirdWrap, map[string]interface{}{
		"third-wrap-string": "test-string",
		"third-wrap-number": 123,
	})
}

func getSecondWrap() error {
	firstWrap := getFirstWrapped()
	errSecondWrap := errors.New("second wrap")
	return WithErrorAndFields(firstWrap, errSecondWrap, map[string]interface{}{
		"second-wrap-string": "test-string",
		"second-wrap-number": 123,
	})
}

func getFirstWrapped() error {
	externalError := getExternalError()
	return NewWithErrorAndFields(externalError, map[string]interface{}{
		"first-wrap-string": "test-string",
		"first-wrap-number": 123,
	})
}

func getExternalError() error {
	return errors.New("previous")
}
