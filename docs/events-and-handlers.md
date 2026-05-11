# Events & Handlers

The `tgbotkit-runtime` uses an event-driven system to process updates from Telegram. You can register handlers for different types of events to implement your bot's logic.

## Event Types

The following core events are emitted by the runtime:

| Event Name | Constant (`events.`) | Description |
|---|---|---|
| `onUpdate` | `OnUpdate` | Emitted for every update received from Telegram. |
| `onMessage` | `OnMessage` | Emitted when a message (text, photo, etc.) is received. |
| `onEditedMessage` | `OnEditedMessage` | Emitted when a message is edited. |
| `onChannelPost` | `OnChannelPost` | Emitted when a channel post is received. |
| `onCallbackQuery` | `OnCallbackQuery` | Emitted when a callback query is received. |
| `onInlineQuery` | `OnInlineQuery` | Emitted when an inline query is received. |
| `onPoll` | `OnPoll` | Emitted when a poll update is received. |
| `onChatMember` | `OnChatMember` | Emitted when a chat member update is received. |
| `onMessageReaction` | `OnMessageReaction` | Emitted when a message reaction update is received. |
| `onCommand` | `OnCommand` | Emitted when a command (e.g., `/start`) is detected. |

## Event Payloads

Each event comes with a specific payload structure:

### `UpdateEvent`
Used for `OnUpdate`.
- `Update`: The raw `*client.Update` object from the Telegram API.

### `MessageEvent`
Used for message-like events such as `OnMessage`, `OnEditedMessage`, `OnChannelPost`, `OnBusinessMessage`, and `OnGuestMessage`.
- `Message`: The `*client.Message` object.
- `Type`: The `messagetype.MessageType` (e.g., `messagetype.Text`, `messagetype.Photo`).

Other Telegram update kinds use dedicated payloads such as `CallbackQueryEvent`, `InlineQueryEvent`, `PollEvent`, `ChatMemberEvent`, and `MessageReactionEvent`.

### `CommandEvent`
Used for `OnCommand`.
- `Message`: The `*client.Message` object that contained the command.
- `Command`: The command name (e.g., `start` for `/start`).
- `Args`: The text following the command.

## Registering Handlers

Handlers are registered via the `Handlers()` method on the `Bot` instance.

### `OnUpdate`
Handles every raw update.

```go
bot.Handlers().OnUpdate(func(ctx context.Context, event *events.UpdateEvent) error {
    log.Printf("Got update ID: %d", event.Update.UpdateId)
    return nil
})
```

### `OnMessage`
Handles all messages.

```go
bot.Handlers().OnMessage(func(ctx context.Context, event *events.MessageEvent) error {
    log.Printf("Got message from: %d", event.Message.From.Id)
    return nil
})
```

### `OnMessageType`
Handles messages of a specific type (e.g., only text, only photos). This is a convenience method that filters `OnMessage` events.

```go
bot.Handlers().OnMessageType(messagetype.Photo, func(ctx context.Context, event *events.MessageEvent) error {
    log.Printf("Got a photo!")
    return nil
})
```

Message-like events have dedicated registration methods such as `OnEditedMessage`, `OnChannelPost`, `OnBusinessMessage`, and `OnGuestMessage`.

Typed non-message updates have matching methods such as `OnCallbackQuery`, `OnInlineQuery`, `OnPoll`, `OnChatMember`, and `OnMessageReaction`.

### `OnCommand`
Handles bot commands.

```go
bot.Handlers().OnCommand(func(ctx context.Context, event *events.CommandEvent) error {
    if event.Command == "help" {
        // Send help message
    }
    return nil
})
```

## Handler Return Values

Handlers return an `error`. In the default bot runtime, ordinary handler errors are logged by the logger middleware. They are reported to an event emitter error handler only when one is configured. Ordinary errors do not stop other handlers for the same event by default.

Return `eventemitter.ErrBreak` when a handler or listener intentionally wants to stop propagation for that event. If you pass a custom event emitter to the bot, that emitter's own `stopOnError` option controls ordinary error propagation.

The default panic recovery middleware converts recovered panics into handler errors. Under the default runtime configuration, those recovered errors are logged and later handlers continue to run.
