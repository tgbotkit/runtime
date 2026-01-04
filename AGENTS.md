# Project Overview

`tgbotkit-runtime` is a Go-based runtime environment designed for building Telegram bots. It provides a robust foundation for handling API interactions, event dispatching, and component management, allowing developers to focus on bot logic.

## Key Features

-   **OpenAPI Client:** Utilizes a pre-generated client (`github.com/tgbotkit/client`) from the Telegram Bot API OpenAPI specification, ensuring full API coverage and type safety.
-   **Event-Driven Architecture:** Built around a pluggable `EventEmitter` system. Components emit events, and handlers subscribe to them.
-   **Modular Components:** Core functionalities like update polling, command parsing, and message classification are implemented as distinct, reusable components.
-   **Flexible Configuration:** Uses the functional options pattern for type-safe and flexible configuration of the bot and its components.
-   **Pluggable Logging:** Interfaces for logging allow for easy swapping of implementations (e.g., `slog`, `zerolog`, or no-op).

## Architecture

The system is orchestrated by the `Bot` struct in `bot.go`, which serves as the central point for configuration and control.

The lifecycle of the bot is as follows:

1.  **Initialization (`New`):**
    -   A `Bot` instance is created using a flexible options pattern.
    -   Essential services like the API `client`, `logger`, and `eventEmitter` are initialized with sensible defaults if not provided by the user.
    -   Internal components, such as the `classifier` and `commandparser`, are subscribed to the event system to process incoming updates.

2.  **Execution (`Run`):**
    -   The `Run` method starts the bot's main processing loop.
    -   It establishes an **Update Source** to receive updates from Telegram. By default, it uses a long-polling mechanism (`updatepoller`) backed by an in-memory offset store, but this can be replaced with another source, like a `webhook`.
    -   The `receiveLoop` continuously fetches updates from the source.

3.  **Event-Driven Processing:**
    -   For each update received, the `receiveLoop` emits a generic `OnUpdateReceived` event.
    -   This triggers a **Processing Pipeline**:
        -   **Classifier:** Listens for `OnUpdateReceived` and emits more specific events (e.g., `OnTextMessageReceived`, `OnCallbackQueryReceived`).
        -   **Command Parser:** Listens for text messages, detects commands (e.g., `/start`), and emits command-specific events.
        -   **User-Defined Handlers:** Custom logic, registered via `AddHandler`, subscribes to any of these events to implement the bot's features.

This architecture decouples the source of updates from the logic that processes them, creating a modular and extensible system.

## Directory Structure

-   `/` - Root directory containing the main `Bot` logic and entry points.
    -   `bot.go`: Main bot initialization and run loop.
    -   `interfaces.go`: Core interfaces like `UpdateSource`.
    -   `options.gen.go`: Generated configuration options.
-   `internal/` - Internal packages not intended for external use.
    -   `classifier/`: Analyzes updates to trigger more specific events.
    -   `commandparser/`: Parses text for bot commands.
-   `updatepoller/`: Implements long-polling mechanism for updates.
    -   `offsetstore/`: Interfaces and implementations for storing update offsets.
-   `eventemitter/` - The event bus implementation (`Sync` emitter provided). Defines `Listener` and `EventEmitter` interfaces and generic typed helpers.
-   `events/` - Definitions of event topics (strings) and payload structures.
-   `examples/` - Example bot implementations (`pingpong`, `webhook`).
-   `logger/` - Logging abstractions and adapters (`zerolog`, `slog`, `noop`).
-   `webhook/` - Webhook handling logic.

## Development

-   **Dependencies:**
    -   `github.com/tgbotkit/client`: The generated Telegram Bot API client.
    -   `github.com/kazhuravlev/options-gen`: Tool for generating functional options.
-   **Conventions:**
    -   Follow Go idioms.
    -   Use `options-gen` for configuration structs.
    -   Components should implement the `Component` interface if they have a lifecycle.