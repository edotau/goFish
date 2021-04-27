package concurrency

import (
	"errors"
	"testing"
)

func Test_ParamError(t *testing.T) {
	baseErr := errors.New("some error")
	err := &ParamError{Key: "x", Err: baseErr}

	assertEqual(t, err.Unwrap(), baseErr, "should unwrap proper error")
	assertEqual(t, err.Error(), "taskflow: parameter \"x\": some error", "should have proper message")
}
