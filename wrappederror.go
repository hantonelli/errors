package errors

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"os"
	"runtime"
	"sort"
	"strings"
)

var (
	gopath string
	goroot string
)

func init() {
	gopath = os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	goroot = os.Getenv("GOROOT")
	if goroot == "" {
		goroot = build.Default.GOROOT
	}
	gopath = gopath + "/src"
	goroot = goroot + "/src"
}

// location describes a source code location.
type location struct {
	file string
	line int
}

// String returns a location where the error was wrap, in the format filename.go:123 format.
func (loc location) String() string {
	return fmt.Sprintf("%s:%d", loc.file, loc.line)
}

// WrappedError specifies the interface for a wrapped error.
type WrappedError interface {
	error
	IsWrappedError() bool
	GetPrevious() error
	GetActual() error
	GetFields() map[string]interface{}
	GetAllFields() map[string]interface{}
	GetStacktrace() string
}

// WrappedErrorImpl is a wrapper for an error chain that allow to specify errors fields.
type WrappedErrorImpl struct {
	actual     error
	previous   error
	stacktrace string

	location location
	fields   map[string]interface{}
}

// IsWrappedError returns always true for this error type.
func (e WrappedErrorImpl) IsWrappedError() bool {
	return true
}

// GetPrevious returns the previous error in the chain.
func (e WrappedErrorImpl) GetPrevious() error {
	return e.previous
}

// GetActual returns the actual wrapped error in the chain.
func (e WrappedErrorImpl) GetActual() error {
	return e.actual
}

// GetFields returns the fields associated with the actual error.
func (e WrappedErrorImpl) GetFields() map[string]interface{} {
	return e.fields
}

// Error returns stack of all the wrapped error messages and it associated fields.
func (e *WrappedErrorImpl) Error() string {
	if e.previous != nil {
		if _, ok := e.previous.(WrappedError); ok {
			return fmt.Sprintf("%s <br> %s", printActual(e), e.previous.Error())
		}
		return fmt.Sprintf("%s <br> Message: %s.", printActual(e), e.previous.Error())
	}
	return printActual(e)
}

func printActual(e *WrappedErrorImpl) string {
	if e.fields != nil && 0 < len(e.fields) {
		return fmt.Sprintf("Message: %v. Location: %v. Fields: %v.", e.actual.Error(), e.location, printFields(e.fields))
	}
	return fmt.Sprintf("Message: %v. Location: %v", e.actual.Error(), e.location)
}

func printFields(fields map[string]interface{}) string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	first := true
	var b bytes.Buffer
	fmt.Fprint(&b, "map[")
	for _, k := range keys {
		if first {
			first = false
		} else {
			b.WriteString(" ")
		}
		fmt.Fprintf(&b, "%s:%v", k, fields[k])
	}
	fmt.Fprint(&b, "]")

	return b.String()
}

// GetAllFields returns a map of the fields for all the errors that had been wrap in the chain.
func (e WrappedErrorImpl) GetAllFields() map[string]interface{} {
	if e.previous != nil {
		if we, ok := e.previous.(WrappedError); ok {
			previousFields := we.GetAllFields()
			newFields := e.fields
			for k := range previousFields {
				newFields[k] = previousFields[k]
			}
			return newFields
		}
		return e.fields
	}
	return e.fields
}

// GetStacktrace returns the stack trace of the first error in the chain.
func (e WrappedErrorImpl) GetStacktrace() string {
	if e.previous != nil {
		if we, ok := e.previous.(WrappedError); ok {
			return we.GetStacktrace()
		}
	}
	return e.stacktrace
}

// NewWithMessage returns a new WrappedErrorImpl with the provided message and fields.
func NewWithMessage(message string, fields map[string]interface{}) *WrappedErrorImpl {
	actual := errors.New(message)
	return createWrappedError(nil, actual, fields)
}

// NewWithError returns a new WrappedErrorImpl with the provided error and fields.
func NewWithError(actual error, fields map[string]interface{}) *WrappedErrorImpl {
	return createWrappedError(nil, actual, fields)
}

// WithError takes the previous error, the actual error and the fields associated with it and returns a new WrappedErrorImpl.
func WithError(previous error, actual error, fields map[string]interface{}) *WrappedErrorImpl {
	return createWrappedError(previous, actual, fields)
}

func createWrappedError(previous error, actual error, fields map[string]interface{}) *WrappedErrorImpl {
	loc := getLocation()
	if actual == nil {
		return nil
	}
	if fields == nil {
		fields = map[string]interface{}{}
	}
	var stacktrace string
	if previous == nil {
		stacktrace = getStacktrace()
	}
	return &WrappedErrorImpl{
		actual:     actual,
		previous:   previous,
		fields:     fields,
		location:   loc,
		stacktrace: stacktrace,
	}
}

func getLocation() location {
	_, file, line, _ := runtime.Caller(3)
	file = cleanFilePath(file)
	return location{file: file, line: line}
}

func cleanFilePath(file string) string {
	if strings.HasPrefix(file, gopath) {
		file = strings.TrimPrefix(file, gopath)
	} else {
		if strings.HasPrefix(file, goroot) {
			file = strings.TrimPrefix(file, goroot)
		}
	}
	return file
}

func getStacktrace() string {
	n := 2
	var b bytes.Buffer
	first := true
	for i := 0; i < 20; i++ {
		_, file, line, ok := runtime.Caller(n + 1)
		if !ok {
			return b.String()
		}
		file = cleanFilePath(file)
		if first {
			first = false
		} else {
			fmt.Fprint(&b, " ")
		}
		fmt.Fprintf(&b, "%s:%d", file, line)
		n++
	}
	return b.String()
}
