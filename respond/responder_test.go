package respond_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/respond"
)

type mockClient struct {
	client.ClientWithResponsesInterface
	sendFunc   func(context.Context, client.SendMessageJSONRequestBody) (*client.SendMessageResponse, error)
	answerFunc func(context.Context, client.AnswerCallbackQueryJSONRequestBody) (*client.AnswerCallbackQueryResponse, error)
}

func (m *mockClient) SendMessageWithResponse(
	ctx context.Context,
	body client.SendMessageJSONRequestBody,
	_ ...client.RequestEditorFn,
) (*client.SendMessageResponse, error) {
	return m.sendFunc(ctx, body)
}

func (m *mockClient) AnswerCallbackQueryWithResponse(
	ctx context.Context,
	body client.AnswerCallbackQueryJSONRequestBody,
	_ ...client.RequestEditorFn,
) (*client.AnswerCallbackQueryResponse, error) {
	return m.answerFunc(ctx, body)
}

func TestResponderSendText(t *testing.T) {
	t.Parallel()

	var got client.SendMessageJSONRequestBody
	responder := respond.New(&mockClient{
		sendFunc: func(_ context.Context, body client.SendMessageJSONRequestBody) (*client.SendMessageResponse, error) {
			got = body

			return sendMessageResponse(client.Message{Chat: client.Chat{Id: body.ChatId}, Text: &body.Text}), nil
		},
	})

	threadID := 7
	directTopicID := 11
	businessConnectionID := "business-1"
	target := respond.ChatTarget{
		ChatID:                42,
		MessageThreadID:       &threadID,
		DirectMessagesTopicID: &directTopicID,
		BusinessConnectionID:  &businessConnectionID,
	}

	msg, err := responder.SendText(
		context.Background(),
		target,
		"<b>hello</b>",
		respond.WithHTML(),
		respond.WithSilent(),
		respond.WithProtectedContent(),
	)
	if err != nil {
		t.Fatalf("SendText() unexpected error: %v", err)
	}
	if msg == nil {
		t.Fatal("SendText() message is nil")
	}
	if got.ChatId != target.ChatID {
		t.Fatalf("ChatId=%d, want %d", got.ChatId, target.ChatID)
	}
	if got.Text != "<b>hello</b>" {
		t.Fatalf("Text=%q, want %q", got.Text, "<b>hello</b>")
	}
	if got.MessageThreadId == nil || *got.MessageThreadId != threadID {
		t.Fatalf("MessageThreadId=%v, want %d", got.MessageThreadId, threadID)
	}
	if got.DirectMessagesTopicId == nil || *got.DirectMessagesTopicId != directTopicID {
		t.Fatalf("DirectMessagesTopicId=%v, want %d", got.DirectMessagesTopicId, directTopicID)
	}
	if got.BusinessConnectionId == nil || *got.BusinessConnectionId != businessConnectionID {
		t.Fatalf("BusinessConnectionId=%v, want %q", got.BusinessConnectionId, businessConnectionID)
	}
	if got.ParseMode == nil || *got.ParseMode != "HTML" {
		t.Fatalf("ParseMode=%v, want HTML", got.ParseMode)
	}
	if got.DisableNotification == nil || !*got.DisableNotification {
		t.Fatalf("DisableNotification=%v, want true", got.DisableNotification)
	}
	if got.ProtectContent == nil || !*got.ProtectContent {
		t.Fatalf("ProtectContent=%v, want true", got.ProtectContent)
	}
}

func TestResponderSendTextInChatAndReplyText(t *testing.T) {
	t.Parallel()

	var got []client.SendMessageJSONRequestBody
	responder := respond.New(&mockClient{
		sendFunc: func(_ context.Context, body client.SendMessageJSONRequestBody) (*client.SendMessageResponse, error) {
			got = append(got, body)

			return sendMessageResponse(client.Message{Chat: client.Chat{Id: body.ChatId}, Text: &body.Text}), nil
		},
	})

	threadID := 3
	topic := &client.DirectMessagesTopic{TopicId: 9}
	businessConnectionID := "business-2"
	source := &client.Message{
		BusinessConnectionId: &businessConnectionID,
		Chat:                 client.Chat{Id: 100},
		DirectMessagesTopic:  topic,
		MessageId:            55,
		MessageThreadId:      &threadID,
	}

	if _, err := responder.SendTextInChat(context.Background(), source, "same chat"); err != nil {
		t.Fatalf("SendTextInChat() unexpected error: %v", err)
	}
	if _, err := responder.ReplyText(context.Background(), source, "reply"); err != nil {
		t.Fatalf("ReplyText() unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("send calls=%d, want 2", len(got))
	}

	assertSourceTarget(t, got[0], source)
	assertSourceTarget(t, got[1], source)
	if got[0].ReplyParameters != nil {
		t.Fatalf("SendTextInChat ReplyParameters=%#v, want nil", got[0].ReplyParameters)
	}
	if got[1].ReplyParameters == nil || got[1].ReplyParameters.MessageId == nil || *got[1].ReplyParameters.MessageId != source.MessageId {
		t.Fatalf("ReplyParameters=%#v, want message id %d", got[1].ReplyParameters, source.MessageId)
	}
}

func TestResponderSendTextPatch(t *testing.T) {
	t.Parallel()

	var got client.SendMessageJSONRequestBody
	responder := respond.New(&mockClient{
		sendFunc: func(_ context.Context, body client.SendMessageJSONRequestBody) (*client.SendMessageResponse, error) {
			got = body

			return sendMessageResponse(client.Message{Chat: client.Chat{Id: body.ChatId}, Text: &body.Text}), nil
		},
	})

	_, err := responder.SendText(
		context.Background(),
		respond.ChatTarget{ChatID: 1},
		"hello",
		respond.WithSendMessagePatch(func(body *client.SendMessageJSONRequestBody) {
			body.Text = "patched"
		}),
	)
	if err != nil {
		t.Fatalf("SendText() unexpected error: %v", err)
	}
	if got.Text != "patched" {
		t.Fatalf("Text=%q, want patched", got.Text)
	}
}

func TestResponderAnswerCallback(t *testing.T) {
	t.Parallel()

	var got client.AnswerCallbackQueryJSONRequestBody
	responder := respond.New(&mockClient{
		answerFunc: func(_ context.Context, body client.AnswerCallbackQueryJSONRequestBody) (*client.AnswerCallbackQueryResponse, error) {
			got = body

			return answerCallbackResponse(true), nil
		},
	})

	query := &client.CallbackQuery{Id: "callback-1"}
	if err := responder.AnswerCallbackAlert(
		context.Background(),
		query,
		"saved",
		respond.WithCallbackCache(5),
		respond.WithCallbackURL("https://example.com"),
	); err != nil {
		t.Fatalf("AnswerCallbackAlert() unexpected error: %v", err)
	}

	if got.CallbackQueryId != query.Id {
		t.Fatalf("CallbackQueryId=%q, want %q", got.CallbackQueryId, query.Id)
	}
	if got.Text == nil || *got.Text != "saved" {
		t.Fatalf("Text=%v, want saved", got.Text)
	}
	if got.ShowAlert == nil || !*got.ShowAlert {
		t.Fatalf("ShowAlert=%v, want true", got.ShowAlert)
	}
	if got.CacheTime == nil || *got.CacheTime != 5 {
		t.Fatalf("CacheTime=%v, want 5", got.CacheTime)
	}
	if got.Url == nil || *got.Url != "https://example.com" {
		t.Fatalf("Url=%v, want https://example.com", got.Url)
	}
}

func TestResponderErrors(t *testing.T) {
	t.Parallel()

	if _, err := respond.TargetFromMessage(nil); !errors.Is(err, respond.ErrNilMessage) {
		t.Fatalf("TargetFromMessage() err=%v, want ErrNilMessage", err)
	}

	responder := respond.New(nil)
	if _, err := responder.SendText(context.Background(), respond.ChatTarget{}, "hello"); !errors.Is(err, respond.ErrNilClient) {
		t.Fatalf("SendText() err=%v, want ErrNilClient", err)
	}
	if err := responder.AnswerCallback(context.Background(), &client.CallbackQuery{}); !errors.Is(err, respond.ErrNilClient) {
		t.Fatalf("AnswerCallback() err=%v, want ErrNilClient", err)
	}

	responder = respond.New(&mockClient{})
	if err := responder.AnswerCallback(context.Background(), nil); !errors.Is(err, respond.ErrNilCallbackQuery) {
		t.Fatalf("AnswerCallback() err=%v, want ErrNilCallbackQuery", err)
	}
}

func assertSourceTarget(t *testing.T, got client.SendMessageJSONRequestBody, source *client.Message) {
	t.Helper()

	if got.ChatId != source.Chat.Id {
		t.Fatalf("ChatId=%d, want %d", got.ChatId, source.Chat.Id)
	}
	if got.MessageThreadId == nil || *got.MessageThreadId != *source.MessageThreadId {
		t.Fatalf("MessageThreadId=%v, want %d", got.MessageThreadId, *source.MessageThreadId)
	}
	if got.DirectMessagesTopicId == nil || *got.DirectMessagesTopicId != int(source.DirectMessagesTopic.TopicId) {
		t.Fatalf("DirectMessagesTopicId=%v, want %d", got.DirectMessagesTopicId, source.DirectMessagesTopic.TopicId)
	}
	if got.BusinessConnectionId == nil || *got.BusinessConnectionId != *source.BusinessConnectionId {
		t.Fatalf("BusinessConnectionId=%v, want %q", got.BusinessConnectionId, *source.BusinessConnectionId)
	}
}

func sendMessageResponse(message client.Message) *client.SendMessageResponse {
	return &client.SendMessageResponse{
		HTTPResponse: &http.Response{StatusCode: http.StatusOK, Status: "200 OK"},
		JSON200: &struct {
			Ok     client.SendMessage200Ok `json:"ok"`
			Result client.Message          `json:"result"`
		}{
			Ok:     true,
			Result: message,
		},
	}
}

func answerCallbackResponse(result bool) *client.AnswerCallbackQueryResponse {
	return &client.AnswerCallbackQueryResponse{
		HTTPResponse: &http.Response{StatusCode: http.StatusOK, Status: "200 OK"},
		JSON200: &struct {
			Ok     client.AnswerCallbackQuery200Ok `json:"ok"`
			Result bool                            `json:"result"`
		}{
			Ok:     true,
			Result: result,
		},
	}
}
