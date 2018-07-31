package webwire

import (
	"context"
	"fmt"
)

// TranslateContextError translates context errors to webwire error types
func TranslateContextError(err error) error {
	if err == context.DeadlineExceeded {
		return NewDeadlineExceededErr(err)
	} else if err == context.Canceled {
		return NewCanceledErr(err)
	}
	return fmt.Errorf("Unexpected context error: %s", err)
}
