package webhook

//go:generate go tool options-gen -out-filename=options.gen.go -from-struct=Options

// Options is the options for the Webhook handler.
type Options struct {
	Token string `option:"optional"`
}
