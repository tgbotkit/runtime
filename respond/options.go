package respond

import "github.com/tgbotkit/client"

// SendTextOption configures a SendMessage request built by Responder.
type SendTextOption func(*client.SendMessageJSONRequestBody)

// WithHTML sets HTML parse mode for sent text.
func WithHTML() SendTextOption {
	return WithParseMode("HTML")
}

// WithMarkdownV2 sets MarkdownV2 parse mode for sent text.
func WithMarkdownV2() SendTextOption {
	return WithParseMode("MarkdownV2")
}

// WithParseMode sets Telegram parse mode for sent text.
func WithParseMode(mode string) SendTextOption {
	return func(body *client.SendMessageJSONRequestBody) {
		body.ParseMode = &mode
	}
}

// WithSilent sends the message without notification.
func WithSilent() SendTextOption {
	return func(body *client.SendMessageJSONRequestBody) {
		value := true
		body.DisableNotification = &value
	}
}

// WithProtectedContent prevents forwarding and saving the sent message.
func WithProtectedContent() SendTextOption {
	return func(body *client.SendMessageJSONRequestBody) {
		value := true
		body.ProtectContent = &value
	}
}

// WithSendMessagePatch applies advanced SendMessage options not wrapped here.
func WithSendMessagePatch(patch func(*client.SendMessageJSONRequestBody)) SendTextOption {
	return func(body *client.SendMessageJSONRequestBody) {
		if patch != nil {
			patch(body)
		}
	}
}

// AnswerCallbackOption configures an AnswerCallbackQuery request built by Responder.
type AnswerCallbackOption func(*client.AnswerCallbackQueryJSONRequestBody)

// WithCallbackCache sets the callback answer cache duration in seconds.
func WithCallbackCache(seconds int) AnswerCallbackOption {
	return func(body *client.AnswerCallbackQueryJSONRequestBody) {
		body.CacheTime = &seconds
	}
}

// WithCallbackURL sets the URL opened by the client for a callback answer.
func WithCallbackURL(url string) AnswerCallbackOption {
	return func(body *client.AnswerCallbackQueryJSONRequestBody) {
		body.Url = &url
	}
}

// WithAnswerCallbackPatch applies advanced AnswerCallbackQuery options not wrapped here.
func WithAnswerCallbackPatch(patch func(*client.AnswerCallbackQueryJSONRequestBody)) AnswerCallbackOption {
	return func(body *client.AnswerCallbackQueryJSONRequestBody) {
		if patch != nil {
			patch(body)
		}
	}
}
