package runtime

import (
	"github.com/metalagman/appkit/lifecycle"
)

// Runnable represents a long-running process that blocks until the context is canceled or an error occurs.
type Runnable = lifecycle.Runnable

// RunnableFunc is a function adapter for the Runnable interface.
type RunnableFunc = lifecycle.RunnableFunc

// Lifecycle represents a component that has a distinct start and stop phase.
type Lifecycle = lifecycle.Lifecycle

// RunnableOption configures the ToRunnable adapter.
type RunnableOption = lifecycle.Option

// ToRunnable converts a Lifecycle to a Runnable.
var ToRunnable = lifecycle.ToRunnable

// WithStartTimeout sets the timeout for the Start operation.
var WithStartTimeout = lifecycle.WithStartTimeout

// WithStopTimeout sets the timeout for the Stop operation.
var WithStopTimeout = lifecycle.WithStopTimeout
