package classifier

import (
	"context"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
)

// Classifier component.
type Classifier struct{}

// New creates a new Classifier.
func New() *Classifier {
	return &Classifier{}
}

func (c *Classifier) Subscribe(bot events.BotContext) {
	ee := bot.EventEmitter()
	eventemitter.On[events.UpdateEvent](ee, events.OnUpdateReceived, c.onUpdateReceived)
	eventemitter.On[events.MessageEvent](ee, events.OnMessageReceived, c.onMessageReceived)
}

func (c *Classifier) onUpdateReceived(ctx context.Context, event *events.UpdateEvent) error {
	if event == nil || event.Update == nil || event.Update.Message == nil {
		return nil
	}
	event.Bot.EventEmitter().Emit(ctx, events.OnMessageReceived, &events.MessageEvent{Bot: event.Bot, Message: event.Update.Message})
	return nil
}

func (c *Classifier) onMessageReceived(ctx context.Context, event *events.MessageEvent) error {
	if event == nil || event.Message == nil {
		return nil
	}
	message := event.Message

	ee := event.Bot.EventEmitter()

	eventType, ok := messageTypeToEvent[messageType(message)]
	if !ok {
		eventType = events.OnUnknownMessageReceived
	}
	ee.Emit(ctx, eventType, event)
	return nil
}

// MessageType used for message routing.
type MessageType string

// Constants for message types.
const (
	// MessageTypeAudio is an audio message.
	MessageTypeAudio MessageType = "audio"
	// MessageTypeContact is a contact message.
	MessageTypeContact MessageType = "contact"
	// MessageTypeDocument is a document message.
	MessageTypeDocument MessageType = "document"
	// MessageTypeLocation is a location message.
	MessageTypeLocation MessageType = "location"
	// MessageTypePhoto is a photo message.
	MessageTypePhoto MessageType = "photo"
	// MessageTypeSticker is a sticker message.
	MessageTypeSticker MessageType = "sticker"
	// MessageTypeText is a text message.
	MessageTypeText MessageType = "text"
	// MessageTypeVenue is a venue message.
	MessageTypeVenue MessageType = "venue"
	// MessageTypeVideo is a video message.
	MessageTypeVideo MessageType = "video"
	// MessageTypeVoice is a voice message.
	MessageTypeVoice MessageType = "voice"
	// MessageTypeVideoNote is a video note message.
	MessageTypeVideoNote MessageType = "video_note"
	// MessageTypeService is a service message.
	MessageTypeService MessageType = "service"
	// MessageTypeUnknown is an unknown message type.
	MessageTypeUnknown MessageType = "unknown"
)

var messageTypeToEvent = map[MessageType]string{
	MessageTypeAudio:     events.OnAudioMessageReceived,
	MessageTypeContact:   events.OnContactMessageReceived,
	MessageTypeDocument:  events.OnDocumentMessageReceived,
	MessageTypeLocation:  events.OnLocationMessageReceived,
	MessageTypePhoto:     events.OnPhotoMessageReceived,
	MessageTypeSticker:   events.OnStickerMessageReceived,
	MessageTypeText:      events.OnTextMessageReceived,
	MessageTypeVenue:     events.OnVenueMessageReceived,
	MessageTypeVideo:     events.OnVideoMessageReceived,
	MessageTypeVoice:     events.OnVoiceMessageReceived,
	MessageTypeVideoNote: events.OnVideoNoteMessageReceived,
	MessageTypeService:   events.OnServiceMessageReceived,
}

func messageType(message *client.Message) MessageType {
	if message == nil {
		return MessageTypeUnknown
	}
	switch {
	case message.Audio != nil:
		return MessageTypeAudio
	case message.Contact != nil:
		return MessageTypeContact
	case message.Document != nil:
		return MessageTypeDocument
	case message.Location != nil:
		return MessageTypeLocation
	case message.Photo != nil:
		return MessageTypePhoto
	case message.Sticker != nil:
		return MessageTypeSticker
	case message.Text != nil:
		return MessageTypeText
	case message.Venue != nil:
		return MessageTypeVenue
	case message.Video != nil:
		return MessageTypeVideo
	case message.Voice != nil:
		return MessageTypeVoice
	case message.VideoNote != nil:
		return MessageTypeVideoNote
	case isServiceMessage(message):
		return MessageTypeService
	default:
		return MessageTypeUnknown
	}
}

func isServiceMessage(message *client.Message) bool {
	return (message.NewChatMembers != nil && len(*message.NewChatMembers) > 0) ||
		message.LeftChatMember != nil ||
		(message.NewChatTitle != nil && *message.NewChatTitle != "") ||
		(message.NewChatPhoto != nil && len(*message.NewChatPhoto) > 0) ||
		(message.DeleteChatPhoto != nil && *message.DeleteChatPhoto) ||
		(message.GroupChatCreated != nil && *message.GroupChatCreated) ||
		(message.SupergroupChatCreated != nil && *message.SupergroupChatCreated) ||
		(message.ChannelChatCreated != nil && *message.ChannelChatCreated) ||
		message.PinnedMessage != nil
}
