package fail

// This code is liberally borrowed from github.com/go-errors/errors
// https://github.com/go-errors/errors/blob/master/error.go

import (
	"bytes"
	"fmt"
	"runtime"

	"github.com/go-errors/errors"
)

// MaxStackDepth is the maximum number of stack frames to report.
var MaxStackDepth = 50

// StackTraceError wraps an error with the caller stack trace. Wrapping an error
// with a stack trace should occur at the boundary of your application code and
// external code. If an error comes from your own code, don't trace it, but if
// an error comes from third party code, such as opening a file, or running a db
// query, you should trace it so you know exactly where in your system the error
// occurred, and what path was taken to get there.
type StackTraceError struct {
	error
	frames  []errors.StackFrame
	stack   []uintptr
	context map[string]interface{}
}

// TraceFace is an interface that represents StackTraceError. By defining an
// interface, we're able to return a nil interface rather than a nil pointer.
// This allows comparison with nil, and therefore doesn't break error checking.
type TraceFace interface {
	Error() string
	WithContext(map[string]interface{}) TraceFace
	WithContextField(string, interface{}) TraceFace
	ErrorContext() map[string]interface{}
}

// Trace is a convenience wrapper for creating a StackTraceError around an error.
// Note this can be called with nil errors, for convenience.
func Trace(err error) TraceFace {
	if err == nil {
		return nil
	}
	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2, stack[:])
	return &StackTraceError{
		error:   err,
		context: map[string]interface{}{},
		stack:   stack[:length],
	}
}

// WithContext adds metadata to a stack trace error for context.
func (err *StackTraceError) WithContext(ctx map[string]interface{}) TraceFace {
	if err == nil {
		return err
	}
	if len(err.context) == 0 {
		err.context = ctx
	} else {
		for key, value := range ctx {
			err.context[key] = value
		}
	}
	return err
}

// WithContextField adds a metadatum to a stack trace error for context.
func (err *StackTraceError) WithContextField(key string, value interface{}) TraceFace {
	if err == nil {
		return err
	}
	err.context[key] = value
	return err
}

// ErrorContext implements the api.ErrorContext interface.
func (err *StackTraceError) ErrorContext() map[string]interface{} {
	return err.context
}

// Error implements the error interface, with the error, type and stack trace.
func (err *StackTraceError) Error() string {
	buffer := bytes.NewBufferString(fmt.Sprintf("%T %s", err.error, err.error.Error()))
	for key, value := range err.context {
		buffer.WriteString(fmt.Sprintf(" %s=%v", key, value))
	}
	buffer.WriteString("\n")
	buffer.Write(err.Stack())
	return buffer.String()
}

// Stack returns the callstack formatted the same way that go does
// in runtime/debug.Stack()
func (err StackTraceError) Stack() []byte {
	buf := bytes.Buffer{}

	for _, frame := range err.StackFrames() {
		buf.WriteString(frame.String())
	}

	return buf.Bytes()
}

// StackFrames returns an array of errors.StackFrames containing information
// about the stack.
func (err *StackTraceError) StackFrames() []errors.StackFrame {
	if err.frames == nil {
		err.frames = make([]errors.StackFrame, len(err.stack))
		for i, pc := range err.stack {
			err.frames[i] = errors.NewStackFrame(pc)
		}
	}
	return err.frames
}
