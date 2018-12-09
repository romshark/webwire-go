package webwire

import (
	"context"
	"fmt"
)

// TranslateContextError translates context errors to webwire error types
func TranslateContextError(err error) error {
	if err == context.DeadlineExceeded {
		return DeadlineExceededErr{Cause: err}
	} else if err == context.Canceled {
		return CanceledErr{Cause: err}
	}
	return fmt.Errorf("unexpected context error: %s", err)
}
