package errors

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func testErrorCode(code int) ErrorCode {
	return ErrorCode{Code: "TEST-" + strconv.Itoa(code), Title: "Test " + strconv.Itoa(code)}
}

func TestErrCodeError(t *testing.T) {
	assertions := require.New(t)
	cause := NewError(testErrorCode(1001), "test 1001 detail", nil)
	err := NewError(testErrorCode(1002), "test 1002 detail", cause)
	errId := err.GetId()
	actualMessage := err.Error()
	stackTrace := err.GetStackTrace()
	fmt.Print(stackTrace)
	assertions.Equal(fmt.Sprintf("ErrCodeError [TEST-1002][%s] test 1002 detail", errId), actualMessage)
	assertions.True(strings.Contains(stackTrace, fmt.Sprintf("Caused by: ErrCodeError [TEST-1001][%s] test 1001 detail", cause.GetId())))
}

func TestMultiCauseError(t *testing.T) {
	assertions := require.New(t)
	err1 := NewError(testErrorCode(1001), "test 1001 detail", nil)
	err2 := NewError(testErrorCode(1002), "test 1002 detail", nil)
	multiErr := NewDefaultMultiCauseError([]ErrCodeErr{err1, err2})
	multiErrId := multiErr.GetId()
	actualMessage := multiErr.Error()
	stackTrace := multiErr.GetStackTrace()
	fmt.Print(stackTrace)
	assertions.Equal(fmt.Sprintf("MultiCauseError [NC-COMMON-2100][%s] multiple independent errors have happened", multiErrId), actualMessage)
	assertions.True(strings.Contains(stackTrace, fmt.Sprintf("Caused by (1/2): ErrCodeError [TEST-1001][%s] test 1001 detail", err1.GetId())))
	assertions.True(strings.Contains(stackTrace, fmt.Sprintf("Caused by (2/2): ErrCodeError [TEST-1002][%s] test 1002 detail", err2.GetId())))
}

type CustomErr struct {
	*ErrCodeError
	Source string
}

func TestCustomErrCodeError(t *testing.T) {
	assertions := require.New(t)
	err := New(CustomErr{}, testErrorCode(1001), "test 1001 detail")
	errId := err.GetId()
	actualMessage := err.Error()
	stackTrace := err.GetStackTrace()
	fmt.Print(stackTrace)
	assertions.Equal(fmt.Sprintf("CustomErr [TEST-1001][%s] test 1001 detail", errId), actualMessage)
	assertions.True(strings.HasPrefix(stackTrace, fmt.Sprintf("CustomErr [TEST-1001][%s] test 1001 detail", errId)))
}

func TestCustomErrCodeErrorWithTemplate(t *testing.T) {
	assertions := require.New(t)
	err := New(CustomErr{Source: "/test"}, testErrorCode(1001), "test 1001 detail")
	errId := err.GetId()
	actualMessage := err.Error()
	stackTrace := err.GetStackTrace()
	fmt.Print(stackTrace)
	assertions.Equal("/test", err.Source)
	assertions.Equal(fmt.Sprintf("CustomErr [TEST-1001][%s] test 1001 detail", errId), actualMessage)
	assertions.True(strings.HasPrefix(stackTrace, fmt.Sprintf("CustomErr [TEST-1001][%s] test 1001 detail", errId)))
}
