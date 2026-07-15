package respond

import (
	"context"
	"fmt"

	"github.com/tgbotkit/client"
)

// Responder sends common Telegram responses through the generated client.
type Responder struct {
	api client.ClientWithResponsesInterface
}

// New creates a responder backed by the generated Telegram API client.
func New(api client.ClientWithResponsesInterface) *Responder {
	return &Responder{api: api}
}

// SendText sends a text message to the target.
func (r *Responder) SendText(
	ctx context.Context,
	target ChatTarget,
	text string,
	opts ...SendTextOption,
) (*client.Message, error) {
	if r == nil || r.api == nil {
		return nil, ErrNilClient
	}

	body := client.SendMessageJSONRequestBody{
		Text: text,
	}
	target.applyTo(&body)
	applySendTextOptions(&body, opts)

	resp, err := r.api.SendMessageWithResponse(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("send message: empty response")
	}

	if resp.JSON200 == nil || !bool(resp.JSON200.Ok) {
		return nil, fmt.Errorf("send message: unexpected response: %s", resp.Status())
	}

	return &resp.JSON200.Result, nil
}

// SendTextInChat sends a text message to the same chat context as source.
func (r *Responder) SendTextInChat(
	ctx context.Context,
	source *client.Message,
	text string,
	opts ...SendTextOption,
) (*client.Message, error) {
	target, err := TargetFromMessage(source)
	if err != nil {
		return nil, err
	}

	return r.SendText(ctx, target, text, opts...)
}

// ReplyText sends a text reply to source.
func (r *Responder) ReplyText(
	ctx context.Context,
	source *client.Message,
	text string,
	opts ...SendTextOption,
) (*client.Message, error) {
	target, err := TargetFromMessage(source)
	if err != nil {
		return nil, err
	}

	opts = append([]SendTextOption{
		WithSendMessagePatch(func(body *client.SendMessageJSONRequestBody) {
			messageID := source.MessageId
			body.ReplyParameters = &client.ReplyParameters{
				MessageId: &messageID,
			}
		}),
	}, opts...)

	return r.SendText(ctx, target, text, opts...)
}

// AnswerCallback answers a callback query.
func (r *Responder) AnswerCallback(
	ctx context.Context,
	query *client.CallbackQuery,
	opts ...AnswerCallbackOption,
) error {
	if r == nil || r.api == nil {
		return ErrNilClient
	}

	if query == nil {
		return ErrNilCallbackQuery
	}

	body := client.AnswerCallbackQueryJSONRequestBody{
		CallbackQueryId: query.Id,
	}
	applyAnswerCallbackOptions(&body, opts)

	resp, err := r.api.AnswerCallbackQueryWithResponse(ctx, body)
	if err != nil {
		return fmt.Errorf("answer callback query: %w", err)
	}

	if resp == nil {
		return fmt.Errorf("answer callback query: empty response")
	}

	if resp.JSON200 == nil || !bool(resp.JSON200.Ok) || !resp.JSON200.Result {
		return fmt.Errorf("answer callback query: unexpected response: %s", resp.Status())
	}

	return nil
}

// AnswerCallbackText answers a callback query with a notification.
func (r *Responder) AnswerCallbackText(
	ctx context.Context,
	query *client.CallbackQuery,
	text string,
	opts ...AnswerCallbackOption,
) error {
	opts = append([]AnswerCallbackOption{
		func(body *client.AnswerCallbackQueryJSONRequestBody) {
			body.Text = &text
		},
	}, opts...)

	return r.AnswerCallback(ctx, query, opts...)
}

// AnswerCallbackAlert answers a callback query with an alert dialog.
func (r *Responder) AnswerCallbackAlert(
	ctx context.Context,
	query *client.CallbackQuery,
	text string,
	opts ...AnswerCallbackOption,
) error {
	opts = append([]AnswerCallbackOption{
		func(body *client.AnswerCallbackQueryJSONRequestBody) {
			body.Text = &text
			value := true
			body.ShowAlert = &value
		},
	}, opts...)

	return r.AnswerCallback(ctx, query, opts...)
}

func applySendTextOptions(body *client.SendMessageJSONRequestBody, opts []SendTextOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(body)
		}
	}
}

func applyAnswerCallbackOptions(
	body *client.AnswerCallbackQueryJSONRequestBody,
	opts []AnswerCallbackOption,
) {
	for _, opt := range opts {
		if opt != nil {
			opt(body)
		}
	}
}
