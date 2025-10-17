package errors

import (
	"fmt"
	"strings"
)

type RemoteErrCodeError struct {
	*ErrCodeError
	Status *int
	Source interface{}
	Meta   map[string]interface{}
}

type RemoteMultiCauseError struct {
	*RemoteErrCodeError
	Causes []*RemoteErrCodeError
}

func NewRemoteErrCodeError(id string, code ErrorCode, detail string, meta map[string]interface{},
	status *int, source interface{}) *RemoteErrCodeError {
	return &RemoteErrCodeError{
		ErrCodeError: &ErrCodeError{
			Id:         id,
			ErrorCode:  code,
			Detail:     detail,
			StackTrace: nil,
			Cause:      nil,
		},
		Status: status,
		Source: source,
		Meta:   meta,
	}
}

func NewRemoteMultiCauseError(id string, code ErrorCode, detail string, meta map[string]interface{},
	status *int, source interface{}, causes []*RemoteErrCodeError) *RemoteMultiCauseError {
	return &RemoteMultiCauseError{
		RemoteErrCodeError: NewRemoteErrCodeError(id, code, detail, meta, status, source),
		Causes:             causes,
	}
}

func (e *RemoteMultiCauseError) GetStackTrace() string {
	s := e.Error() + "\n"
	prefix := " "
	for i, cause := range e.Causes {
		s += fmt.Sprintf("Caused by (%d/%d): ", i+1, len(e.Causes))
		stackTrace := cause.GetStackTrace()
		nlCount := strings.Count(stackTrace, "\n")
		s += strings.Replace(stackTrace, "\n", "\n"+prefix, nlCount-1)
	}
	return s
}
