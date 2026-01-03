package eventemitter

//go:generate go tool options-gen -out-filename=options.gen.go -from-struct=Options

type Options struct {
	errorHandler ErrorHandler `option:"optional"`
	stopOnError  bool         `option:"optional" default:"true"`
}
