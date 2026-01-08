// Package messagetype provides constants and utilities for identifying the type of a Telegram message.
// It allows for easy routing and handling of different message contents and service events.
package messagetype

import (
	"github.com/tgbotkit/client"
)

// MessageType represents the classification of a Telegram message.
type MessageType string

// Detect inspects a Telegram message and returns its most specific type.
// It prioritizes service messages (e.g., chat member changes) over standard content (e.g., text).
//
//nolint:gocyclo,cyclop,funlen,maintidx
func Detect(message *client.Message) MessageType {
	if message == nil {
		return Unknown
	}

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

	// Standard messages
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

// Standard Content Types.
const (
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
)

// Chat Lifecycle & Management.
const (
	NewChatMembers                MessageType = "new_chat_members"
	LeftChatMember                MessageType = "left_chat_member"
	NewChatTitle                  MessageType = "new_chat_title"
	NewChatPhoto                  MessageType = "new_chat_photo"
	DeleteChatPhoto               MessageType = "delete_chat_photo"
	GroupChatCreated              MessageType = "group_chat_created"
	SupergroupChatCreated         MessageType = "supergroup_chat_created"
	ChannelChatCreated            MessageType = "channel_chat_created"
	MessageAutoDeleteTimerChanged MessageType = "message_auto_delete_timer_changed"
	MigrateFromChatID             MessageType = "migrate_from_chat_id"
	MigrateToChatID               MessageType = "migrate_to_chat_id"
)

// Topic Management.
const (
	ForumTopicCreated         MessageType = "forum_topic_created"
	ForumTopicEdited          MessageType = "forum_topic_edited"
	ForumTopicClosed          MessageType = "forum_topic_closed"
	ForumTopicReopened        MessageType = "forum_topic_reopened"
	GeneralForumTopicHidden   MessageType = "general_forum_topic_hidden"
	GeneralForumTopicUnhidden MessageType = "general_forum_topic_unhidden"
)

// Video Chat Events.
const (
	VideoChatScheduled           MessageType = "video_chat_scheduled"
	VideoChatStarted             MessageType = "video_chat_started"
	VideoChatEnded               MessageType = "video_chat_ended"
	VideoChatParticipantsInvited MessageType = "video_chat_participants_invited"
)

// Payments & Financial.
const (
	SuccessfulPayment MessageType = "successful_payment"
	RefundedPayment   MessageType = "refunded_payment"
)

// Giveaways & Gifts.
const (
	GiveawayCompleted MessageType = "giveaway_completed"
	GiveawayCreated   MessageType = "giveaway_created"
	GiveawayWinners   MessageType = "giveaway_winners"
	Gift              MessageType = "gift"
	GiftUpgradeSent   MessageType = "gift_upgrade_sent"
	UniqueGift        MessageType = "unique_gift"
	BoostAdded        MessageType = "boost_added"
)

// Suggested Posts & Ads.
const (
	DirectMessagePriceChanged   MessageType = "direct_message_price_changed"
	PaidMessagePriceChanged     MessageType = "paid_message_price_changed"
	SuggestedPostApprovalFailed MessageType = "suggested_post_approval_failed"
	SuggestedPostApproved       MessageType = "suggested_post_approved"
	SuggestedPostDeclined       MessageType = "suggested_post_declined"
	SuggestedPostPaid           MessageType = "suggested_post_paid"
	SuggestedPostRefunded       MessageType = "suggested_post_refunded"
)

// Miscellaneous Service Events.
const (
	PinnedMessage           MessageType = "pinned_message"
	UsersShared             MessageType = "users_shared"
	ChatShared              MessageType = "chat_shared"
	WriteAccessAllowed      MessageType = "write_access_allowed"
	ProximityAlertTriggered MessageType = "proximity_alert_triggered"
	WebAppData              MessageType = "web_app_data"
	ChatBackgroundSet       MessageType = "chat_background_set"
	ChecklistTasksAdded     MessageType = "checklist_tasks_added"
	ChecklistTasksDone      MessageType = "checklist_tasks_done"
	PassportData            MessageType = "passport_data"
	ConnectedWebsite        MessageType = "connected_website"
)

// Special Types.
const (
	// Unknown is returned when the message type cannot be identified.
	Unknown MessageType = "unknown"
)
