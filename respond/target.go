package respond

import (
	"fmt"

	"github.com/tgbotkit/client"
)

// ChatTarget identifies where a message should be sent.
type ChatTarget struct {
	ChatID                int64
	MessageThreadID       *int
	DirectMessagesTopicID *int
	BusinessConnectionID  *string
}

// TargetFromMessage builds a send target from a source Telegram message.
func TargetFromMessage(message *client.Message) (ChatTarget, error) {
	if message == nil {
		return ChatTarget{}, ErrNilMessage
	}

	target := ChatTarget{
		ChatID:               message.Chat.Id,
		MessageThreadID:      message.MessageThreadId,
		BusinessConnectionID: message.BusinessConnectionId,
	}

	if message.DirectMessagesTopic != nil {
		topicID, err := intFromInt64(message.DirectMessagesTopic.TopicId)
		if err != nil {
			return ChatTarget{}, err
		}

		target.DirectMessagesTopicID = &topicID
	}

	return target, nil
}

func (t ChatTarget) applyTo(body *client.SendMessageJSONRequestBody) {
	body.ChatId = t.ChatID
	body.MessageThreadId = t.MessageThreadID
	body.DirectMessagesTopicId = t.DirectMessagesTopicID
	body.BusinessConnectionId = t.BusinessConnectionID
}

func intFromInt64(value int64) (int, error) {
	maxInt := int64(int(^uint(0) >> 1))
	minInt := -maxInt - 1

	if value > maxInt || value < minInt {
		return 0, fmt.Errorf("direct messages topic id out of range: %d", value)
	}

	return int(value), nil
}
