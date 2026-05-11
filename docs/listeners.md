# Core Listeners

`tgbotkit-runtime` includes core listeners that drive the event-driven architecture by transforming generic events into more specific ones.

## Classifier Listener

The **Classifier** is responsible for analyzing incoming raw updates (`OnUpdate`) and emitting more specialized events.

-   **Listen to:** `OnUpdate`
-   **Emit:** `OnMessage`, `OnEditedMessage`, `OnChannelPost`, `OnCallbackQuery`, `OnInlineQuery`, poll, chat member, business, reaction, and payment-related events.

Currently, it detects the following:
-   **Message-like updates:** `message`, `edited_message`, channel posts, business messages, and guest messages are emitted as `MessageEvent` payloads with `messagetype` classification.
-   **Query updates:** callback, inline, chosen inline result, shipping, and pre-checkout queries are emitted as typed payloads.
-   **State updates:** polls, chat member changes, join requests, boosts, reactions, business connections, paid media purchases, and managed bot updates are emitted as typed payloads.

## Command Parser Listener

The **Command Parser** listens for text messages and identifies bot commands (e.g., `/start`).

-   **Listen to:** `OnMessage` (filtered for `messagetype.Text`)
-   **Emit:** `OnCommand`

### Features
-   **Command Identification:** Detects commands starting with `/` using Telegram's message entities.
-   **Bot Mentions:** Correctly handles commands with bot mentions (e.g., `/start@my_bot`). It only processes commands that either have no mention or are specifically mentioned to this bot.
-   **Argument Parsing:** Separates the command name from its arguments (the rest of the message text).
-   **Early Termination:** If a command is found, it returns `eventemitter.ErrBreak`. This stops later `OnMessage` listeners for that event, while the emitted `OnCommand` handlers still run.

## Internal vs External Listeners

These listeners are registered automatically during bot initialization in `runtime.New()`. While they are "internal" to the runtime's default configuration, they are implemented using the same public `eventemitter.Listener` interface that you use for your own bot logic.

### Disabling Core Listeners
Core listeners are registered by default. Use `runtime.WithDefaultListenersEnabled(false)` for a fully custom event pipeline. Use `runtime.WithDefaultMiddlewareEnabled(false)` if you also need to opt out of the built-in context, logging, and panic-recovery middleware.
