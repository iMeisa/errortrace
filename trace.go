package errortrace

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

// NewTrace creates a new ErrorTrace object with an error
func NewTrace(err error) ErrorTrace {
	newTrace := ErrorTrace{
		Error: err,
	}

	newTrace.addTrace(2)

	return newTrace
}

// NilTrace creates a new ErrorTrace with nil error
func NilTrace() ErrorTrace {
	return ErrorTrace{
		Error: nil,
	}
}

// ErrorTrace holds the error and traces back to all places where the error occurred
type ErrorTrace struct {
	Error  error
	traces []trace
}

func (t *ErrorTrace) addTrace(skip int) {
	pc, file, line, _ := runtime.Caller(skip)

	splitFuncDir := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	funcName := splitFuncDir[len(splitFuncDir)-1]

	// Get working directory
	path, err := os.Getwd()
	if err != nil {
		log.Println("Error reading path during addTrace() call: ", err)
		return
	}

	t.traces = append(t.traces, trace{file: file[len(path):], funcName: funcName, line: line})

}

// HasError return bool if error is nil,
// adds trace if true
func (t *ErrorTrace) HasError() bool {
	hasError := t.Error != nil

	if hasError {
		t.addTrace(2)
	}

	return hasError
}

// Read reads back all the locations where the error went through
func (t *ErrorTrace) Read() {
	t.addTrace(2)

	fmt.Println(strings.Repeat("-", 30))
	log.Println("\nERROR: \n\t", t.Error)

	// Tell user to add traces
	if len(t.traces) < 1 {
		log.Println("NO TRACES FOUND")
		return
	}

	originTrace := t.traces[0]

	fmt.Println("ORIGIN: \n\t", originTrace.file, " -> ", originTrace.funcName+": line", originTrace.line)

	if len(t.traces) > 1 {
		fmt.Println("CALLED BY:")

		for _, followingTrace := range t.traces[1:] {
			fmt.Println("\t", followingTrace.file, " -> ", followingTrace.funcName+": line", followingTrace.line)
		}
	}

	fmt.Println(strings.Repeat("-", 30))
}

func (t *ErrorTrace) ErrorString() string {
	return t.Error.Error()
}

// trace holds valid traceback info
type trace struct {
	file     string
	funcName string
	line     int
}
