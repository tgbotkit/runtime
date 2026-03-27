# Events & Handlers

The `tgbotkit-runtime` uses an event-driven system to process updates from Telegram. You can register handlers for different types of events to implement your bot's logic.

## Event Types

The following core events are emitted by the runtime:

| Event Name | Constant (`events.`) | Description |
|---|---|---|
| `onUpdate` | `OnUpdate` | Emitted for every update received from Telegram. |
| `onMessage` | `OnMessage` | Emitted when a message (text, photo, etc.) is received. |
| `onCommand` | `OnCommand` | Emitted when a command (e.g., `/start`) is detected. |

## Event Payloads

Each event comes with a specific payload structure:

### `UpdateEvent`
Used for `OnUpdate`.
- `Update`: The raw `*client.Update` object from the Telegram API.

### `MessageEvent`
Used for `OnMessage`.
- `Message`: The `*client.Message` object.
- `Type`: The `messagetype.MessageType` (e.g., `messagetype.Text`, `messagetype.Photo`).

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

Handlers return an `error`. If a handler returns an error, it will be logged (if the logger middleware is active), but it won't stop the processing of other handlers for the same event.
