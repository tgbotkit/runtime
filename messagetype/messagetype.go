// Package messagetype provides constants and utilities for identifying the type of a Telegram message.
package messagetype

import (
	"github.com/tgbotkit/client"
)

// MessageType used for message routing.
type MessageType string

// Detect inspects the message and returns its type.
func Detect(message *client.Message) MessageType {
	if message == nil {
		return Unknown
	}

	// The order is important: check for specific service messages first,
	// then for standard content types.
	switch {
	// Service messages
	case message.NewChatMembers != nil:
		return NewChatMembers
	case message.LeftChatMember != nil:
		return LeftChatMember
	case message.NewChatTitle != nil:
		return NewChatTitle
	case message.NewChatPhoto != nil:
		return NewChatPhoto
	case message.DeleteChatPhoto != nil:
		return DeleteChatPhoto
	case message.GroupChatCreated != nil:
		return GroupChatCreated
	case message.SupergroupChatCreated != nil:
		return SupergroupChatCreated
	case message.ChannelChatCreated != nil:
		return ChannelChatCreated
	case message.MessageAutoDeleteTimerChanged != nil:
		return MessageAutoDeleteTimerChanged
	case message.MigrateToChatId != nil:
		return MigrateToChatID
	case message.MigrateFromChatId != nil:
		return MigrateFromChatID
	case message.PinnedMessage != nil:
		return PinnedMessage
	case message.SuccessfulPayment != nil:
		return SuccessfulPayment
	case message.RefundedPayment != nil:
		return RefundedPayment
	case message.UsersShared != nil:
		return UsersShared
	case message.ChatShared != nil:
		return ChatShared
	case message.WriteAccessAllowed != nil:
		return WriteAccessAllowed
	case message.ProximityAlertTriggered != nil:
		return ProximityAlertTriggered
	case message.ForumTopicCreated != nil:
		return ForumTopicCreated
	case message.ForumTopicEdited != nil:
		return ForumTopicEdited
	case message.ForumTopicClosed != nil:
		return ForumTopicClosed
	case message.ForumTopicReopened != nil:
		return ForumTopicReopened
	case message.GeneralForumTopicHidden != nil:
		return GeneralForumTopicHidden
	case message.GeneralForumTopicUnhidden != nil:
		return GeneralForumTopicUnhidden
	case message.VideoChatScheduled != nil:
		return VideoChatScheduled
	case message.VideoChatStarted != nil:
		return VideoChatStarted
	case message.VideoChatEnded != nil:
		return VideoChatEnded
	case message.VideoChatParticipantsInvited != nil:
		return VideoChatParticipantsInvited
	case message.WebAppData != nil:
		return WebAppData
	case message.BoostAdded != nil:
		return BoostAdded
	case message.ChatBackgroundSet != nil:
		return ChatBackgroundSet
	case message.ChecklistTasksAdded != nil:
		return ChecklistTasksAdded
	case message.ChecklistTasksDone != nil:
		return ChecklistTasksDone
	case message.DirectMessagePriceChanged != nil:
		return DirectMessagePriceChanged
	case message.Gift != nil:
		return Gift
	case message.GiftUpgradeSent != nil:
		return GiftUpgradeSent
	case message.GiveawayCompleted != nil:
		return GiveawayCompleted
	case message.GiveawayCreated != nil:
		return GiveawayCreated
	case message.GiveawayWinners != nil:
		return GiveawayWinners
	case message.PaidMessagePriceChanged != nil:
		return PaidMessagePriceChanged
	case message.SuggestedPostApprovalFailed != nil:
		return SuggestedPostApprovalFailed
	case message.SuggestedPostApproved != nil:
		return SuggestedPostApproved
	case message.SuggestedPostDeclined != nil:
		return SuggestedPostDeclined
	case message.SuggestedPostPaid != nil:
		return SuggestedPostPaid
	case message.SuggestedPostRefunded != nil:
		return SuggestedPostRefunded
	case message.UniqueGift != nil:
		return UniqueGift
	case message.PassportData != nil:
		return PassportData
	case message.ConnectedWebsite != nil:
		return ConnectedWebsite

	// Standard content types
	case message.Text != nil:
		return Text
	case message.Animation != nil:
		return Animation
	case message.Audio != nil:
		return Audio
	case message.Document != nil:
		return Document
	case message.Photo != nil:
		return Photo
	case message.Sticker != nil:
		return Sticker
	case message.Story != nil:
		return Story
	case message.Video != nil:
		return Video
	case message.VideoNote != nil:
		return VideoNote
	case message.Voice != nil:
		return Voice
	case message.Contact != nil:
		return Contact
	case message.Dice != nil:
		return Dice
	case message.Game != nil:
		return Game
	case message.Poll != nil:
		return Poll
	case message.Venue != nil:
		return Venue
	case message.Location != nil:
		return Location
	case message.Invoice != nil:
		return Invoice
	case message.Checklist != nil:
		return Checklist
	case message.PaidMedia != nil:
		return PaidMedia
	case message.Giveaway != nil:
		return Giveaway

	default:
		return Unknown
	}
}

// Constants for message types, derived from the client.Message struct.
const (
	// Standard content message types.
	Text      MessageType = "text"
	Animation MessageType = "animation"
	Audio     MessageType = "audio"
	Document  MessageType = "document"
	Photo     MessageType = "photo"
	Sticker   MessageType = "sticker"
	Story     MessageType = "story"
	Video     MessageType = "video"
	VideoNote MessageType = "video_note"
	Voice     MessageType = "voice"
	Contact   MessageType = "contact"
	Dice      MessageType = "dice"
	Game      MessageType = "game"
	Poll      MessageType = "poll"
	Venue     MessageType = "venue"
	Location  MessageType = "location"
	Invoice   MessageType = "invoice"
	Checklist MessageType = "checklist"
	PaidMedia MessageType = "paid_media"
	Giveaway  MessageType = "giveaway"

	// Service message types.
	BoostAdded                   MessageType = "boost_added"
	ChatBackgroundSet            MessageType = "chat_background_set"
	ChatShared                   MessageType = "chat_shared"
	ChecklistTasksAdded          MessageType = "checklist_tasks_added"
	ChecklistTasksDone           MessageType = "checklist_tasks_done"
	ConnectedWebsite             MessageType = "connected_website"
	DeleteChatPhoto              MessageType = "delete_chat_photo"
	DirectMessagePriceChanged    MessageType = "direct_message_price_changed"
	ForumTopicClosed             MessageType = "forum_topic_closed"
	ForumTopicCreated            MessageType = "forum_topic_created"
	ForumTopicEdited             MessageType = "forum_topic_edited"
	ForumTopicReopened           MessageType = "forum_topic_reopened"
	GeneralForumTopicHidden      MessageType = "general_forum_topic_hidden"
	GeneralForumTopicUnhidden    MessageType = "general_forum_topic_unhidden"
	Gift                         MessageType = "gift"
	GiftUpgradeSent              MessageType = "gift_upgrade_sent"
	GiveawayCompleted            MessageType = "giveaway_completed"
	GiveawayCreated              MessageType = "giveaway_created"
	GiveawayWinners              MessageType = "giveaway_winners"
	GroupChatCreated             MessageType = "group_chat_created"
	LeftChatMember               MessageType = "left_chat_member"
	MessageAutoDeleteTimerChanged  MessageType = "message_auto_delete_timer_changed"
	MigrateFromChatID            MessageType = "migrate_from_chat_id"
	MigrateToChatID              MessageType = "migrate_to_chat_id"
	NewChatMembers               MessageType = "new_chat_members"
	NewChatPhoto                 MessageType = "new_chat_photo"
	NewChatTitle                 MessageType = "new_chat_title"
	PaidMessagePriceChanged      MessageType = "paid_message_price_changed"
	PassportData                 MessageType = "passport_data"
	PinnedMessage                MessageType = "pinned_message"
	ProximityAlertTriggered      MessageType = "proximity_alert_triggered"
	RefundedPayment              MessageType = "refunded_payment"
	SuccessfulPayment            MessageType = "successful_payment"
	SuggestedPostApprovalFailed  MessageType = "suggested_post_approval_failed"
	SuggestedPostApproved        MessageType = "suggested_post_approved"
	SuggestedPostDeclined        MessageType = "suggested_post_declined"
	SuggestedPostPaid            MessageType = "suggested_post_paid"
	SuggestedPostRefunded        MessageType = "suggested_post_refunded"
	SupergroupChatCreated        MessageType = "supergroup_chat_created"
	ChannelChatCreated           MessageType = "channel_chat_created"
	UniqueGift                   MessageType = "unique_gift"
	UsersShared                  MessageType = "users_shared"
	VideoChatEnded               MessageType = "video_chat_ended"
	VideoChatParticipantsInvited MessageType = "video_chat_participants_invited"
	VideoChatScheduled           MessageType = "video_chat_scheduled"
	VideoChatStarted             MessageType = "video_chat_started"
	WebAppData                   MessageType = "web_app_data"
	WriteAccessAllowed           MessageType = "write_access_allowed"

	// Unknown message type.
	Unknown MessageType = "unknown"
)