// Package events defines the events emitted by the bot and their associated payload types.
package events

// Constants for event names.
const (
	// OnUpdate is emitted when a new update is received from Telegram.
	OnUpdate = "onUpdate"
	// OnMessage is emitted when a new message is received, regardless of its type.
	// The specific type is available in the MessageEvent.Type field.
	OnMessage = "onMessage"
	// OnEditedMessage is emitted when a message is edited.
	OnEditedMessage = "onEditedMessage"
	// OnChannelPost is emitted when a new channel post is received.
	OnChannelPost = "onChannelPost"
	// OnEditedChannelPost is emitted when a channel post is edited.
	OnEditedChannelPost = "onEditedChannelPost"
	// OnBusinessMessage is emitted when a business message is received.
	OnBusinessMessage = "onBusinessMessage"
	// OnEditedBusinessMessage is emitted when a business message is edited.
	OnEditedBusinessMessage = "onEditedBusinessMessage"
	// OnGuestMessage is emitted when a guest message is received.
	OnGuestMessage = "onGuestMessage"
	// OnCommand is emitted when a command is received from a text message.
	OnCommand = "onCommand"
	// OnCallbackQuery is emitted when a callback query is received.
	OnCallbackQuery = "onCallbackQuery"
	// OnInlineQuery is emitted when an inline query is received.
	OnInlineQuery = "onInlineQuery"
	// OnChosenInlineResult is emitted when an inline result is chosen.
	OnChosenInlineResult = "onChosenInlineResult"
	// OnShippingQuery is emitted when a shipping query is received.
	OnShippingQuery = "onShippingQuery"
	// OnPreCheckoutQuery is emitted when a pre-checkout query is received.
	OnPreCheckoutQuery = "onPreCheckoutQuery"
	// OnPoll is emitted when a poll update is received.
	OnPoll = "onPoll"
	// OnPollAnswer is emitted when a poll answer update is received.
	OnPollAnswer = "onPollAnswer"
	// OnChatMember is emitted when a chat member update is received.
	OnChatMember = "onChatMember"
	// OnMyChatMember is emitted when the bot's chat member state changes.
	OnMyChatMember = "onMyChatMember"
	// OnChatJoinRequest is emitted when a chat join request is received.
	OnChatJoinRequest = "onChatJoinRequest"
	// OnChatBoost is emitted when a chat boost update is received.
	OnChatBoost = "onChatBoost"
	// OnRemovedChatBoost is emitted when a chat boost is removed.
	OnRemovedChatBoost = "onRemovedChatBoost"
	// OnMessageReaction is emitted when a message reaction update is received.
	OnMessageReaction = "onMessageReaction"
	// OnMessageReactionCount is emitted when a reaction-count update is received.
	OnMessageReactionCount = "onMessageReactionCount"
	// OnBusinessConnection is emitted when a business connection changes.
	OnBusinessConnection = "onBusinessConnection"
	// OnDeletedBusinessMessages is emitted when business messages are deleted.
	OnDeletedBusinessMessages = "onDeletedBusinessMessages"
	// OnPurchasedPaidMedia is emitted when paid media is purchased.
	OnPurchasedPaidMedia = "onPurchasedPaidMedia"
	// OnManagedBot is emitted when a managed bot update is received.
	OnManagedBot = "onManagedBot"
)
