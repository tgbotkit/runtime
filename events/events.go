package events

// Constants for event names.
const (
	OnBeforeStart = "onBeforeStart"
	OnStart       = "onStart"
	OnStop        = "onStop"

	// OnUpdateReceived is emitted when a new update is received from Telegram.
	OnUpdateReceived = "onUpdateReceived"
	// OnMessageReceived is emitted when a new message is received.
	OnMessageReceived = "onMessageReceived"
	// OnAudioMessageReceived is emitted when an audio message is received.
	OnAudioMessageReceived = "onAudioMessageReceived"
	// OnContactMessageReceived is emitted when a contact message is received.
	OnContactMessageReceived = "onContactMessageReceived"
	// OnDocumentMessageReceived is emitted when a document message is received.
	OnDocumentMessageReceived = "onDocumentMessageReceived"
	// OnLocationMessageReceived is emitted when a location message is received.
	OnLocationMessageReceived = "onLocationMessageReceived"
	// OnPhotoMessageReceived is emitted when a photo message is received.
	OnPhotoMessageReceived = "onPhotoMessageReceived"
	// OnStickerMessageReceived is emitted when a sticker message is received.
	OnStickerMessageReceived = "onStickerMessageReceived"
	// OnTextMessageReceived is emitted when a text message is received.
	OnTextMessageReceived = "onTextMessageReceived"
	// OnVideoMessageReceived is emitted when a video message is received.
	OnVideoMessageReceived = "onVideoMessageReceived"
	// OnVenueMessageReceived is emitted when a venue message is received.
	OnVenueMessageReceived = "onVenueMessageReceived"
	// OnVoiceMessageReceived is emitted when a voice message is received.
	OnVoiceMessageReceived = "onVoiceMessageReceived"
	// OnVideoNoteMessageReceived is emitted when a video note message is received.
	OnVideoNoteMessageReceived = "onVideoNoteMessageReceived"
	// OnServiceMessageReceived is emitted when a service message is received.
	OnServiceMessageReceived = "onServiceMessageReceived"
	// OnUnknownMessageReceived is emitted when an unknown type of message is received.
	OnUnknownMessageReceived = "onUnknownTypeMessageReceived"

	// OnCommand is emitted when a command is received.
	OnCommand = "onCommand"
)
