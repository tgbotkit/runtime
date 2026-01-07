package webhook

import "github.com/tgbotkit/client"

//go:generate go tool options-gen -out-filename=options.gen.go -from-struct=Options

// Options is the options for the Webhook handler.
type Options struct {
	token  string                              `option:"optional"`
	url    string                              `option:"optional"`
	client client.ClientWithResponsesInterface `option:"optional"`
}
