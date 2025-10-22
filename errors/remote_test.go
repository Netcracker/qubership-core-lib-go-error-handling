package errors

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRemoteErrCodeError_Basic(t *testing.T) {
	code := ErrorCode{"R001", "remote-failure"}
	status := 500
	source := "remote-service"
	meta := map[string]interface{}{"region": "eu-west"}

	err := NewRemoteErrCodeError("id123", code, "something broke", meta, &status, source)

	require.NotNil(t, err)
	assert.Equal(t, "R001", err.ErrorCode.Code)
	assert.Equal(t, "something broke", err.Detail)
	require.NotNil(t, err.Status)
	assert.Equal(t, 500, *err.Status)
	assert.Equal(t, source, err.Source)
	assert.Equal(t, "eu-west", err.Meta["region"])
	assert.Equal(t, "id123", err.Id)
}

func TestNewRemoteMultiCauseError_Basic(t *testing.T) {
	code := ErrorCode{"RMC1", "remote-multi"}
	status := 404
	source := "gateway"
	meta := map[string]interface{}{"env": "dev"}

	cause1 := NewRemoteErrCodeError("c1", ErrorCode{"R-A", "a"}, "cause1", nil, &status, source)
	cause2 := NewRemoteErrCodeError("c2", ErrorCode{"R-B", "b"}, "cause2", nil, &status, source)

	err := NewRemoteMultiCauseError("id999", code, "top-level", meta, &status, source, []*RemoteErrCodeError{cause1, cause2})

	require.NotNil(t, err)
	assert.Equal(t, "RMC1", err.ErrorCode.Code)
	assert.Equal(t, "top-level", err.Detail)
	require.Len(t, err.Causes, 2)
	assert.Equal(t, "c1", err.Causes[0].Id)
	assert.Equal(t, "cause2", err.Causes[1].Detail)
}

func TestRemoteMultiCauseError_GetStackTrace(t *testing.T) {
	status := 503
	code := ErrorCode{"REM-ERR", "multi-err"}
	source := "remote"
	meta := map[string]interface{}{"key": "val"}

	cause1 := NewRemoteErrCodeError("c1", ErrorCode{"E1", "child1"}, "child-detail-1", nil, &status, source)
	cause2 := NewRemoteErrCodeError("c2", ErrorCode{"E2", "child2"}, "child-detail-2", nil, &status, source)

	err := NewRemoteMultiCauseError("main", code, "parent-detail", meta, &status, source, []*RemoteErrCodeError{cause1, cause2})
	trace := err.GetStackTrace()

	assert.Contains(t, trace, "parent-detail")
	assert.Contains(t, trace, "Caused by (1/2):")
	assert.Contains(t, trace, "Caused by (2/2):")
	assert.Contains(t, trace, "child-detail-1")
	assert.Contains(t, trace, "child-detail-2")
	assert.GreaterOrEqual(t, strings.Count(trace, "Caused by"), 2)
}

func TestRemoteMultiCauseError_GetStackTrace_NewlinesHandled(t *testing.T) {
	status := 200
	code := ErrorCode{"NLCODE", "nl-title"}
	source := "srv"

	cause := NewRemoteErrCodeError("c", ErrorCode{"SUB", "child"}, "child-detail\nline2", nil, &status, source)
	err := NewRemoteMultiCauseError("main", code, "parent", nil, &status, source, []*RemoteErrCodeError{cause})

	trace := err.GetStackTrace()
	assert.Contains(t, trace, "\n ") // Indentation for multiline
	assert.Contains(t, trace, "line2")
}
