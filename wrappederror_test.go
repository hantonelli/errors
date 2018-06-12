package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"strings"
)

var fields = map[string]interface{}{
	"key":  "value",
	"key2": 12,
}

func TestNewWithMessage(t *testing.T) {

	t.Run("handle message and no fields", func(t *testing.T) {
		message := fmt.Sprintf("error vii, %v", 22)

		var err error = NewWithMsgAndFields(message, nil)
		assert.Equal(t, err.Error(), "Message: error vii, 22. Location: /github.com/hantonelli/errors/wrappederror_test.go:22")
		assert.NotNil(t, err)
		we, ok := err.(WrappedError)
		if !ok {
			t.Error("Error is not WrappedError")
		}
		expectedFields := map[string]interface{}{}
		if !reflect.DeepEqual(expectedFields, we.GetAllFields()) {
			t.Fatalf("expected we.GetAllFields() to be %v, but got %v", expectedFields, we.GetAllFields())
		}
	})
}

func TestNewWithError(t *testing.T) {

	t.Run("return nil if error provided is nil", func(t *testing.T) {
		err := NewWithErrorAndFields(nil, nil)
		assert.Nil(t, err)
	})

	t.Run("should handle error and no fields", func(t *testing.T) {
		errActual := fmt.Errorf("error vii, %v", 22)
		var err error = NewWithErrorAndFields(errActual, nil)
		assert.Equal(t, err.Error(), "Message: error vii, 22. Location: /github.com/hantonelli/errors/wrappederror_test.go:45")
		assert.NotNil(t, err)
		we, ok := err.(WrappedError)
		if !ok {
			t.Error("Error is not WrappedError")
		}
		expectedFields := map[string]interface{}{}
		if !reflect.DeepEqual(expectedFields, we.GetAllFields()) {
			t.Fatalf("expected we.GetAllFields() to be %v, but got %v", expectedFields, we.GetAllFields())
		}
	})

	t.Run("should handle error and fields", func(t *testing.T) {
		errActual := fmt.Errorf("error vii, %v", 22)
		var err error = NewWithErrorAndFields(errActual, fields)
		assert.Equal(t, err.Error(), "Message: error vii, 22. Location: /github.com/hantonelli/errors/wrappederror_test.go:60. Fields: map[key:value key2:12].")
		assert.NotNil(t, err)
		we, ok := err.(WrappedError)
		if !ok {
			t.Error("Error is not WrappedError")
		}
		if !reflect.DeepEqual(fields, we.GetAllFields()) {
			t.Fatalf("expected we.GetAllFields() to be %v, but got %v", fields, we.GetAllFields())
		}
	})

	t.Run("should wrap a previous error", func(t *testing.T) {
		errPrevious := errors.New("previous")
		errActual := fmt.Errorf("error vii, %v", 22)
		var err error = WithErrorAndFields(errPrevious, errActual, fields)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "Message: error vii, 22. Location: /github.com/hantonelli/errors/wrappederror_test.go:75. Fields: map[key:value key2:12]. <br> Message: previous.")
		assert.NotNil(t, err)
		we, ok := err.(WrappedError)
		if !ok {
			t.Error("Error is not WrappedError")
		}
		if !reflect.DeepEqual(fields, we.GetAllFields()) {
			t.Fatalf("expected we.GetAllFields() to be %v, but got %v", fields, we.GetAllFields())
		}
	})

	t.Run("should handle 3 level wrap", func(t *testing.T) {
		err := getThirdWrap()
		assert.NotNil(t, err)
		expectedMsg := `
			Message: third wrap.
			 Location: /github.com/hantonelli/errors/wrappederror_helper_test.go:18.
			 Fields: map[third-wrap-number:123 third-wrap-string:test-string]. <br>

			 Message: second wrap.
			 Location: /github.com/hantonelli/errors/wrappederror_helper_test.go:27.
			 Fields: map[second-wrap-number:123 second-wrap-string:test-string]. <br>

			 Message: previous.
			 Location: /github.com/hantonelli/errors/wrappederror_helper_test.go:35.
			 Fields: map[first-wrap-number:123 first-wrap-string:test-string].
		`
		expectedMsg = cleanSpaces(expectedMsg)
		assert.Equal(t, err.Error(), expectedMsg)

		we, ok := err.(WrappedError)
		if !ok {
			t.Error("Error is not WrappedError")
		}
		expectedFields := map[string]interface{}{
			"third-wrap-string":  "test-string",
			"third-wrap-number":  123,
			"second-wrap-string": "test-string",
			"second-wrap-number": 123,
			"first-wrap-string":  "test-string",
			"first-wrap-number":  123,
		}
		if !reflect.DeepEqual(expectedFields, we.GetAllFields()) {
			t.Fatalf("expected we.GetAllFields() to be %v, but got %v", expectedFields, we.GetAllFields())
		}
	})

	t.Run("should return stacktrace", func(t *testing.T) {
		expectedStacktrace := `
			/github.com/hantonelli/errors/wrappederror_helper_test.go:35
			 /github.com/hantonelli/errors/wrappederror_helper_test.go:25
			 /github.com/hantonelli/errors/wrappederror_helper_test.go:16
			 /github.com/hantonelli/errors/wrappederror_test.go:133
		`
		expectedStacktrace = cleanSpaces(expectedStacktrace)

		err := getThirdWrap()
		assert.NotNil(t, err)
		we, ok := err.(WrappedError)
		if !ok {
			t.Error("Error is not WrappedError")
		}
		actualStacktrace := we.GetStacktrace()
		assert.True(t, strings.HasPrefix(actualStacktrace, expectedStacktrace))
	})
}
