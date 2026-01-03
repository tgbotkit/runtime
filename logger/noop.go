package logger

// Nop is a no-op logger implementation.
type Nop struct{}

// NewNop creates a new Nop logger.
func NewNop() Logger {
	return Nop{}
}

func (Nop) Errorf(format string, args ...interface{}) {}
func (Nop) Fatalf(format string, args ...interface{}) {}
func (Nop) Fatal(args ...interface{})                 {}
func (Nop) Infof(format string, args ...interface{})  {}
func (Nop) Info(args ...interface{})                  {}
func (Nop) Warnf(format string, args ...interface{})  {}
func (Nop) Debugf(format string, args ...interface{}) {}
func (Nop) Debug(args ...interface{})                 {}
