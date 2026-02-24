// execution/result.go
package execution

import "time"

// Result is the output of executing an invocation.
type Result struct {
	Output   any
	Error    error
	Duration time.Duration
}