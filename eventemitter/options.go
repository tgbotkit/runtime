package eventemitter

//go:generate go tool options-gen -out-filename=options.gen.go -from-struct=Options

// Options defines the configuration for an EventEmitter.
type Options struct {
	stopOnError  bool         `default:"true"    option:"optional"`
	errorHandler ErrorHandler `option:"optional"`
}
