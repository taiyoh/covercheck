package covercheck

import "context"

// Checker provides healthcheck implementation.
type Checker func(context.Context) error
