package logger

// Nop is a logger that does nothing.
type Nop struct{}

// NewNop creates a new Nop.
func NewNop() Nop {
	return Nop{}
}

// Errorf logs a message at the error level.
func (Nop) Errorf(_ string, _ ...interface{}) {}

// Fatalf logs a message at the fatal level and calls os.Exit(1).
func (Nop) Fatalf(_ string, _ ...interface{}) {}

// Fatal logs a message at the fatal level and calls os.Exit(1).
func (Nop) Fatal(_ ...interface{})            {}

// Infof logs a message at the info level.
func (Nop) Infof(_ string, _ ...interface{})  {}

// Info logs a message at the info level.
func (Nop) Info(_ ...interface{})             {}

// Warnf logs a message at the warn level.
func (Nop) Warnf(_ string, _ ...interface{})  {}

// Debugf logs a message at the debug level.
func (Nop) Debugf(_ string, _ ...interface{}) {}

// Debug logs a message at the debug level.
func (Nop) Debug(_ ...interface{})            {}