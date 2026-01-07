package listeners

import (
	"context"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/messagetype"
)

func Classifier(emitter eventemitter.EventEmitter) eventemitter.Listener {
	return eventemitter.ListenerFunc(func(ctx context.Context, payload any) error {
		if event, ok := payload.(*events.UpdateEvent); ok {
			if event != nil && event.Update != nil {
				classifyUpdate(ctx, emitter, event)
			}
		}
		return nil
	})
}

// classifyUpdate inspects the update and emits corresponding events.
func classifyUpdate(ctx context.Context, emitter eventemitter.EventEmitter, event *events.UpdateEvent) {
	update := event.Update

	if update.Message != nil {
		msgType := getMessageType(update.Message)
		messageEvent := &events.MessageEvent{
			Message: update.Message,
			Type:    msgType,
		}
		emitter.Emit(ctx, events.OnMessage, messageEvent)
	}
}

func getMessageType(message *client.Message) messagetype.MessageType {
	if message == nil {
		return messagetype.Unknown
	}

	// The order is important: check for specific service messages first,
	// then for standard content types.
	switch {
	// Service messages
	case message.NewChatMembers != nil:
		return messagetype.NewChatMembers
	case message.LeftChatMember != nil:
		return messagetype.LeftChatMember
	case message.NewChatTitle != nil:
		return messagetype.NewChatTitle
	case message.NewChatPhoto != nil:
		return messagetype.NewChatPhoto
	case message.DeleteChatPhoto != nil:
		return messagetype.DeleteChatPhoto
	case message.GroupChatCreated != nil:
		return messagetype.GroupChatCreated
	case message.SupergroupChatCreated != nil:
		return messagetype.SupergroupChatCreated
	case message.ChannelChatCreated != nil:
		return messagetype.ChannelChatCreated
	case message.MessageAutoDeleteTimerChanged != nil:
		return messagetype.MessageAutoDeleteTimerChanged
	case message.MigrateToChatId != nil:
		return messagetype.MigrateToChatID
	case message.MigrateFromChatId != nil:
		return messagetype.MigrateFromChatID
	case message.PinnedMessage != nil:
		return messagetype.PinnedMessage
	case message.SuccessfulPayment != nil:
		return messagetype.SuccessfulPayment
	case message.RefundedPayment != nil:
		return messagetype.RefundedPayment
	case message.UsersShared != nil:
		return messagetype.UsersShared
	case message.ChatShared != nil:
		return messagetype.ChatShared
	case message.WriteAccessAllowed != nil:
		return messagetype.WriteAccessAllowed
	case message.ProximityAlertTriggered != nil:
		return messagetype.ProximityAlertTriggered
	case message.ForumTopicCreated != nil:
		return messagetype.ForumTopicCreated
	case message.ForumTopicEdited != nil:
		return messagetype.ForumTopicEdited
	case message.ForumTopicClosed != nil:
		return messagetype.ForumTopicClosed
	case message.ForumTopicReopened != nil:
		return messagetype.ForumTopicReopened
	case message.GeneralForumTopicHidden != nil:
		return messagetype.GeneralForumTopicHidden
	case message.GeneralForumTopicUnhidden != nil:
		return messagetype.GeneralForumTopicUnhidden
	case message.VideoChatScheduled != nil:
		return messagetype.VideoChatScheduled
	case message.VideoChatStarted != nil:
		return messagetype.VideoChatStarted
	case message.VideoChatEnded != nil:
		return messagetype.VideoChatEnded
	case message.VideoChatParticipantsInvited != nil:
		return messagetype.VideoChatParticipantsInvited
	case message.WebAppData != nil:
		return messagetype.WebAppData
	case message.BoostAdded != nil:
		return messagetype.BoostAdded
	case message.ChatBackgroundSet != nil:
		return messagetype.ChatBackgroundSet
	case message.ChecklistTasksAdded != nil:
		return messagetype.ChecklistTasksAdded
	case message.ChecklistTasksDone != nil:
		return messagetype.ChecklistTasksDone
	case message.DirectMessagePriceChanged != nil:
		return messagetype.DirectMessagePriceChanged
	case message.Gift != nil:
		return messagetype.Gift
	case message.GiftUpgradeSent != nil:
		return messagetype.GiftUpgradeSent
	case message.GiveawayCompleted != nil:
		return messagetype.GiveawayCompleted
	case message.GiveawayCreated != nil:
		return messagetype.GiveawayCreated
	case message.GiveawayWinners != nil:
		return messagetype.GiveawayWinners
	case message.PaidMessagePriceChanged != nil:
		return messagetype.PaidMessagePriceChanged
	case message.SuggestedPostApprovalFailed != nil:
		return messagetype.SuggestedPostApprovalFailed
	case message.SuggestedPostApproved != nil:
		return messagetype.SuggestedPostApproved
	case message.SuggestedPostDeclined != nil:
		return messagetype.SuggestedPostDeclined
	case message.SuggestedPostPaid != nil:
		return messagetype.SuggestedPostPaid
	case message.SuggestedPostRefunded != nil:
		return messagetype.SuggestedPostRefunded
	case message.UniqueGift != nil:
		return messagetype.UniqueGift
	case message.PassportData != nil:
		return messagetype.PassportData
	case message.ConnectedWebsite != nil:
		return messagetype.ConnectedWebsite

	// Standard content types
	case message.Text != nil:
		return messagetype.Text
	case message.Animation != nil:
		return messagetype.Animation
	case message.Audio != nil:
		return messagetype.Audio
	case message.Document != nil:
		return messagetype.Document
	case message.Photo != nil:
		return messagetype.Photo
	case message.Sticker != nil:
		return messagetype.Sticker
	case message.Story != nil:
		return messagetype.Story
	case message.Video != nil:
		return messagetype.Video
	case message.VideoNote != nil:
		return messagetype.VideoNote
	case message.Voice != nil:
		return messagetype.Voice
	case message.Contact != nil:
		return messagetype.Contact
	case message.Dice != nil:
		return messagetype.Dice
	case message.Game != nil:
		return messagetype.Game
	case message.Poll != nil:
		return messagetype.Poll
	case message.Venue != nil:
		return messagetype.Venue
	case message.Location != nil:
		return messagetype.Location
	case message.Invoice != nil:
		return messagetype.Invoice
	case message.Checklist != nil:
		return messagetype.Checklist
	case message.PaidMedia != nil:
		return messagetype.PaidMedia
	case message.Giveaway != nil:
		return messagetype.Giveaway

	default:
		return messagetype.Unknown
	}
}
