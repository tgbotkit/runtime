# Core Listeners

`tgbotkit-runtime` includes core listeners that drive the event-driven architecture by transforming generic events into more specific ones.

## Classifier Listener

The **Classifier** is responsible for analyzing incoming raw updates (`OnUpdate`) and emitting more specialized events.

-   **Listen to:** `OnUpdate`
-   **Emit:** `OnMessage`, `OnCallbackQuery` (and others as implemented).

Currently, it detects the following:
-   **Messages:** If the update contains a `Message`, it emits an `OnMessage` event. It also uses the `messagetype` package to determine the type of the message (e.g., `Text`, `Photo`, `Sticker`, etc.).

## Command Parser Listener

The **Command Parser** listens for text messages and identifies bot commands (e.g., `/start`).

-   **Listen to:** `OnMessage` (filtered for `messagetype.Text`)
-   **Emit:** `OnCommand`

### Features
-   **Command Identification:** Detects commands starting with `/` using Telegram's message entities.
-   **Bot Mentions:** Correctls handles commands with bot mentions (e.g., `/start@my_bot`). It only processes commands that either have no mention or are specifically mentioned to this bot.
-   **Argument Parsing:** Separates the command name from its arguments (the rest of the message text).
-   **Early Termination:** If a command is found, it returns `eventemitter.ErrBreak`. This stops other listeners from processing the `OnMessage` event further, which is useful for ensuring only the command handler executes.

## Internal vs External Listeners

These listeners are registered automatically during bot initialization in `runtime.New()`. While they are "internal" to the runtime's default configuration, they are implemented using the same public `eventemitter.Listener` interface that you use for your own bot logic.

### Disabling Core Listeners
Currently, core listeners are always registered. If you need a completely custom behavior, you could manually configure your own `EventEmitter` and pass it to `runtime.NewOptions`.
