package errors

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"errors"
)

func TestContainsGenericError(t *testing.T) {

	errOther := errors.New("other error")
	ge := GenericError{}
	firstWrapType := errors.New("first")
	firstFields :=  map[string]interface{}{
		"first-wrap-string": "test-string",
		"first-wrap-number": 123,
	}
	secondWrapType := errors.New("second")
	secondFields := map[string]interface{}{
		"second-wrap-string": "test-string",
		"second-wrap-number": 123,
	}

	t.Run("should return nil when it does not exists", func(t *testing.T){
		_, _, ok := ContainsGenericError(errOther)
		assert.False(t, ok)
	})

	t.Run("should return GE when it is not wrap", func(t *testing.T){
		extractedGE, fields, ok := ContainsGenericError(ge)
		assert.True(t, ok)
		assert.Equal(t, ge, extractedGE)
		assert.Equal(t, map[string]interface{}{}, fields)
	})

	t.Run("should return nil when other error is wrap", func(t *testing.T){
		err := errors.New("other error")
		firstWrap :=  WithError(err, firstWrapType, firstFields)
		assert.NotNil(t, firstWrap)

		_, _, ok := ContainsGenericError(firstWrap)
		assert.False(t, ok)
	})

	t.Run("should return GE when it is wrap once", func(t *testing.T){
		firstWrap :=  WithError(ge, firstWrapType, firstFields)
		assert.NotNil(t, firstWrap)

		ge, fields, ok := ContainsGenericError(firstWrap)
		assert.True(t, ok)
		assert.Equal(t, ge, ge)
		assert.True(t, ge.IsGenericError())
		assert.Equal(t, "GenericError", ge.Error())
		if !reflect.DeepEqual(firstFields, fields) {
			t.Fatalf("expected fields to be %v, but got %v", firstFields, fields)
		}
	})

	t.Run("should return GE when it is wrap twice", func(t *testing.T){
		firstWrap :=  WithError(errOther, ge, firstFields)
		assert.NotNil(t, firstWrap)

		secondWrap :=  WithError(firstWrap, secondWrapType, secondFields)
		assert.NotNil(t, secondWrap)

		ge, fields, ok := ContainsGenericError(secondWrap)
		assert.True(t, ok)
		assert.NotNil(t, ge)
		assert.True(t, ge.IsGenericError())
		if !reflect.DeepEqual(firstFields, fields) {
			t.Fatalf("expected fields to be %v, but got %v", firstFields, fields)
		}
	})

	t.Run("should return GE when it is in the middle", func(t *testing.T){
		firstWrap :=  WithError(ge, firstWrapType, firstFields)
		assert.NotNil(t, firstWrap)

		secondWrap :=  WithError(firstWrap, secondWrapType, secondFields)
		assert.NotNil(t, secondWrap)

		ge, fields, ok := ContainsGenericError(secondWrap)
		assert.True(t, ok)
		assert.NotNil(t, ge)
		assert.True(t, ge.IsGenericError())
		if !reflect.DeepEqual(firstFields, fields) {
			t.Fatalf("expected fields to be %v, but got %v", firstFields, fields)
		}
	})

	t.Run("should return GE when it is at last", func(t *testing.T){
		firstWrap :=  WithError(errOther, firstWrapType, firstFields)
		assert.NotNil(t, firstWrap)

		secondWrap :=  WithError(firstWrap, ge, secondFields)
		assert.NotNil(t, secondWrap)

		ge, fields, ok := ContainsGenericError(secondWrap)
		assert.True(t, ok)
		assert.NotNil(t, ge)
		assert.True(t, ge.IsGenericError())
		if !reflect.DeepEqual(secondFields, fields) {
			t.Fatalf("expected fields to be %v, but got %v", secondFields, fields)
		}
	})
}