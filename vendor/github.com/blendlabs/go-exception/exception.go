package exception

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

// Exception is an error with a stack trace.
type Exception struct {
	message        string
	stackTrace     []string
	innerException *Exception
}

// MarshalJSON is a custom json marshaler.
func (e *Exception) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Decompose())
}

// Decompose returns a decomposed version of the exception object.
func (e *Exception) Decompose() map[string]interface{} {
	values := map[string]interface{}{}
	values["message"] = e.message
	values["stack_trace"] = e.StackTrace()

	if e.innerException != nil {
		values["inner_exception"] = e.innerException.Decompose()
	}
	return values
}

// Message returns the exception message.
func (e *Exception) Message() string {
	if e.innerException == nil {
		return e.message
	}
	return fmt.Sprintf("%s Inner Exception: %s", e.message, e.innerException.Message())
}

// StackTrace returns the exception stack trace.
func (e *Exception) StackTrace() []string {
	return e.stackTrace
}

// StackString returns the stack trace formated nicely.
func (e *Exception) StackString() string {
	return formatStackTrace(e.stackTrace)
}

// InnerException returns the nested exception.
func (e *Exception) InnerException() *Exception {
	return e.innerException
}

// Error implements the `error` interface
func (e *Exception) Error() string {
	message := fmt.Sprintf("Exception: %s", e.message)
	message = message + fmt.Sprintf("\n%11s", "At: ")
	message = message + formatStackTrace(prefixLines(spaces(11), e.StackTrace()))

	if e.innerException == e {
		panic("exception loop cycle length 1")
	}

	if e.innerException != nil {
		innerErrorMessage := e.innerException.Error()
		message = message + fmt.Sprintf("\n\nWrapped Exception: %s", innerErrorMessage)
	}
	return message
}

// New returns a new exception by `Sprint`ing the messageComponents.
func New(messageComponents ...interface{}) error {
	message := fmt.Sprint(messageComponents...)
	if len(message) == 0 {
		message = "An Exception Occurred"
	}
	return &Exception{message: message, stackTrace: callerInfo()}
}

// Newf returns a new exception by `Sprintf`ing the format and the args.
func Newf(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	if len(message) == 0 {
		message = "An Exception Occurred"
	}
	return &Exception{message: message, stackTrace: callerInfo()}
}

// Wrap wraps an exception, will return error-typed `nil` if the exception is nil.
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	if typedEx, isException := err.(*Exception); isException {
		return typedEx
	}
	return WrapError(err)
}

// WrapPrefix wraps an exception and allows user to add a custom message to the error.
// This is useful when we want to retain the original stack trace but also add customized
// message at the same time. Will return error-typed `nil` if the exception is nil.
func WrapPrefix(err error, prefix string) error {
	if err == nil {
		return nil
	}

	typedEx := Wrap(err).(*Exception)
	newMessage := fmt.Sprintf("%s: %s", prefix, typedEx.message)

	return &Exception{
		message:    newMessage,
		stackTrace: typedEx.stackTrace,
	}
}

// WrapMany is vestigal and is an API compatability shim for `Nest(...)`.
func WrapMany(err ...error) error {
	return Nest(err...)
}

// Nest nests an arbitrary number of exceptions.
func Nest(err ...error) error {
	var ex *Exception
	var last *Exception
	var didSet bool //(*Exception)(nil) != nil

	for _, e := range err {
		if e != nil {
			wrappedError := Wrap(e).(*Exception)
			if wrappedError != nil && wrappedError != ex {
				if ex == nil {
					ex = wrappedError
					last = wrappedError
				} else {
					last.innerException = wrappedError
					last = wrappedError
				}
				didSet = true
			}
		}
	}
	if didSet {
		return ex
	}
	return nil
}

// WrapError is a shortcut method for wrapping an error by calling .Message() on it.
func WrapError(err error) error {
	if err == nil {
		return nil
	}
	return New(err.Error())
}

// GetStackTrace is a utility method to get the current stack trace at call time.
func GetStackTrace() string {
	return formatStackTrace(callerInfo())
}

// IsException is a helper function that returns if an error is an exception.
func IsException(err error) bool {
	if _, typedOk := err.(*Exception); typedOk {
		return true
	}
	return false
}

// AsException is a helper method that returns an error as an exception.
func AsException(err error) *Exception {
	if typed, typedOk := err.(*Exception); typedOk {
		return typed
	}
	return nil
}

func prefixLines(prefix string, lines []string) []string {
	outputLines := []string{}
	for i, line := range lines {
		if i == 0 {
			outputLines = append(outputLines, line)
		} else {
			outputLines = append(outputLines, fmt.Sprintf("%s%s", prefix, line))

		}
	}
	return outputLines
}

func spaces(num int) string {
	return repeatString(" ", num)
}

func repeatString(token string, num int) string {
	str := ""
	for i := 0; i < num; i++ {
		str = str + token
	}
	return str
}

func formatStackTrace(stack []string) string {
	return strings.Join(stack, "\n")
}

func callerInfo() []string {
	callers := []string{}

	for i := 0; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			return callers
		}

		if file == "<autogenerated>" {
			break
		}

		parts := strings.Split(file, "/")
		dir := parts[len(parts)-2]
		file = parts[len(parts)-1]

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}
		fullyQualifiedName := f.Name()
		segments := strings.Split(fullyQualifiedName, ".")
		functionName := segments[len(segments)-1]

		if dir != "go-exception" {
			caller := fmt.Sprintf("%s:%d %s()", file, line, functionName)
			callers = append(callers, caller)
		}
	}

	return callers
}
