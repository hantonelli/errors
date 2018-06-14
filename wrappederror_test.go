package errors

import (
	goerr "errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fields = map[string]interface{}{
	"key":  "value",
	"key2": 12,
}

func TestNewWithMessage(t *testing.T) {

	t.Run("handle message and no fields", func(t *testing.T) {
		message := fmt.Sprintf("error vii, %v", 22)

		var err error = NewWithMsgAndFields(message, nil)
		assert.Equal(t, err.Error(), "Message: error vii, 22. Location: /github.com/hantonelli/errors/wrappederror_test.go:23")
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
		assert.Equal(t, err.Error(), "Message: error vii, 22. Location: /github.com/hantonelli/errors/wrappederror_test.go:46")
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
		assert.Equal(t, err.Error(), "Message: error vii, 22. Location: /github.com/hantonelli/errors/wrappederror_test.go:61. Fields: map[key:value key2:12].")
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
		errPrevious := goerr.New("previous")
		errActual := fmt.Errorf("error vii, %v", 22)
		var err error = WithErrorAndFields(errPrevious, errActual, fields)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "Message: error vii, 22. Location: /github.com/hantonelli/errors/wrappederror_test.go:76. Fields: map[key:value key2:12]. <br> Message: previous.")
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
			 /github.com/hantonelli/errors/wrappederror_test.go:134
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

func TestContainsError(t *testing.T) {
	errLookFor := goerr.New("expected error")

	t.Run("return nil if error to look for is nil", func(t *testing.T) {
		actualErr, _, ok := ContainsError(errLookFor, nil)
		assert.False(t, ok)
		assert.Nil(t, actualErr)
	})

	t.Run("return nil if error to look for is nil", func(t *testing.T) {
		actualErr := goerr.New("other error")
		foundErr, _, ok := ContainsError(nil, actualErr)
		assert.False(t, ok)
		assert.Nil(t, foundErr)
	})

	t.Run("return error if we look for the same error that is provided", func(t *testing.T) {
		err := goerr.New("expected error")
		actualErr, actualFields, ok := ContainsError(errLookFor, err)
		assert.True(t, ok)
		assert.Equal(t, errLookFor, actualErr)
		assert.Equal(t, actualFields, map[string]interface{}{})
	})

	t.Run("return error if it is wrap", func(t *testing.T) {
		err := goerr.New("expected error")
		wrappedError1 := NewWithErrorAndFields(err, fields)

		actualErr, actualFields, ok := ContainsError(errLookFor, wrappedError1)
		assert.True(t, ok)
		assert.Equal(t, errLookFor, actualErr)
		if !reflect.DeepEqual(fields, actualFields) {
			t.Fatalf("expected fields to be %v, but got %v", fields, actualFields)
		}
	})

	t.Run("return error if it is wrap three times and is in the end", func(t *testing.T) {
		err := goerr.New("expected error")
		wrappedError1 := NewWithErrorAndFields(err, fields)
		err2 := goerr.New("err 2")
		wrappedError2 := WithError(wrappedError1, err2)
		err3 := goerr.New("err 3")
		wrappedError3 := WithError(wrappedError2, err3)

		actualErr, actualFields, ok := ContainsError(errLookFor, wrappedError3)
		assert.True(t, ok)
		assert.Equal(t, errLookFor, actualErr)
		if !reflect.DeepEqual(fields, actualFields) {
			t.Fatalf("expected fields to be %v, but got %v", fields, actualFields)
		}
	})

	t.Run("return error if it is wrap three times and is in the middle", func(t *testing.T) {
		err := goerr.New("err 1")
		wrappedError1 := NewWithError(err)
		err2 := goerr.New("expected error")
		wrappedError2 := WithErrorAndFields(wrappedError1, err2, fields)
		err3 := goerr.New("err 3")
		wrappedError3 := WithError(wrappedError2, err3)

		actualErr, actualFields, ok := ContainsError(errLookFor, wrappedError3)
		assert.True(t, ok)
		assert.Equal(t, errLookFor, actualErr)
		if !reflect.DeepEqual(fields, actualFields) {
			t.Fatalf("expected fields to be %v, but got %v", fields, actualFields)
		}
	})

	t.Run("return error if it is wrap three times and is the first", func(t *testing.T) {
		err := goerr.New("err 1")
		wrappedError1 := NewWithError(err)
		err2 := goerr.New("err 3")
		wrappedError2 := WithError(wrappedError1, err2)
		err3 := goerr.New("expected error")
		wrappedError3 := WithErrorAndFields(wrappedError2, err3, fields)

		actualErr, actualFields, ok := ContainsError(errLookFor, wrappedError3)
		assert.True(t, ok)
		assert.Equal(t, errLookFor, actualErr)
		if !reflect.DeepEqual(fields, actualFields) {
			t.Fatalf("expected fields to be %v, but got %v", fields, actualFields)
		}
	})
}

func TestContainsErrorPrefix(t *testing.T) {
	errMsg := "expected error to occur"
	errLookFor := goerr.New(errMsg)
	lookForPrefix := "expected error"

	t.Run("return nil if error to look for is nil", func(t *testing.T) {
		actualErr, _, ok := ContainsErrorPrefix(lookForPrefix, nil)
		assert.False(t, ok)
		assert.Nil(t, actualErr)
	})

	t.Run("return nil if error to look for is nil", func(t *testing.T) {
		actualErr := goerr.New("other error")
		foundErr, _, ok := ContainsErrorPrefix("", actualErr)
		assert.False(t, ok)
		assert.Nil(t, foundErr)
	})

	t.Run("return error if we look for the same error that is provided", func(t *testing.T) {
		err := goerr.New(errMsg)
		actualErr, actualFields, ok := ContainsErrorPrefix(lookForPrefix, err)
		assert.True(t, ok)
		assert.Equal(t, errLookFor, actualErr)
		assert.Equal(t, actualFields, map[string]interface{}{})
	})

	t.Run("return error if it is wrap", func(t *testing.T) {
		err := goerr.New(errMsg)
		wrappedError1 := NewWithErrorAndFields(err, fields)

		actualErr, actualFields, ok := ContainsErrorPrefix(lookForPrefix, wrappedError1)
		assert.True(t, ok)
		assert.Equal(t, errLookFor, actualErr)
		if !reflect.DeepEqual(fields, actualFields) {
			t.Fatalf("expected fields to be %v, but got %v", fields, actualFields)
		}
	})

	t.Run("return error if it is wrap three times and is in the end", func(t *testing.T) {
		err := goerr.New(errMsg)
		wrappedError1 := NewWithErrorAndFields(err, fields)
		err2 := goerr.New("err 2")
		wrappedError2 := WithError(wrappedError1, err2)
		err3 := goerr.New("err 3")
		wrappedError3 := WithError(wrappedError2, err3)

		actualErr, actualFields, ok := ContainsErrorPrefix(lookForPrefix, wrappedError3)
		assert.True(t, ok)
		assert.Equal(t, errLookFor, actualErr)
		if !reflect.DeepEqual(fields, actualFields) {
			t.Fatalf("expected fields to be %v, but got %v", fields, actualFields)
		}
	})

	t.Run("return error if it is wrap three times and is in the middle", func(t *testing.T) {
		err := goerr.New("err 1")
		wrappedError1 := NewWithError(err)
		err2 := goerr.New(errMsg)
		wrappedError2 := WithErrorAndFields(wrappedError1, err2, fields)
		err3 := goerr.New("err 3")
		wrappedError3 := WithError(wrappedError2, err3)

		actualErr, actualFields, ok := ContainsErrorPrefix(lookForPrefix, wrappedError3)
		assert.True(t, ok)
		assert.Equal(t, errLookFor, actualErr)
		if !reflect.DeepEqual(fields, actualFields) {
			t.Fatalf("expected fields to be %v, but got %v", fields, actualFields)
		}
	})

	t.Run("return error if it is wrap three times and is the first", func(t *testing.T) {
		err := goerr.New("err 1")
		wrappedError1 := NewWithError(err)
		err2 := goerr.New("err 3")
		wrappedError2 := WithError(wrappedError1, err2)
		err3 := goerr.New(errMsg)
		wrappedError3 := WithErrorAndFields(wrappedError2, err3, fields)

		actualErr, actualFields, ok := ContainsErrorPrefix(lookForPrefix, wrappedError3)
		assert.True(t, ok)
		assert.Equal(t, errLookFor, actualErr)
		if !reflect.DeepEqual(fields, actualFields) {
			t.Fatalf("expected fields to be %v, but got %v", fields, actualFields)
		}
	})
}
