# Middleware

Middlewares allow you to intercept and modify the execution of event handlers. They are useful for cross-cutting concerns like logging, error recovery, and context injection.

## How it Works

A middleware in `tgbotkit-runtime` is a function that wraps a `Listener`. It has the following signature:

```go
type Middleware interface {
    Handle(next Listener) Listener
}
```

When an event is emitted, it passes through all registered middlewares before reaching the actual listeners.

## Built-in Middlewares

`tgbotkit-runtime` comes with several built-in middlewares that are enabled by default:

### `ContextInjector`
Injects the `Bot` instance into the `context.Context`. This allows any handler or listener to access the bot's API client and other services.

### `Logger`
Logs the processing of every event and any errors returned by handlers.

### `Recoverer`
Recovers from panics in any listener or handler, preventing the entire bot process from crashing.

## Registering Middleware

You can register your own middleware using the `EventEmitter().Use()` method. You can apply middleware to specific events or to all events using the `*` wildcard.

```go
bot, _ := runtime.New(runtime.NewOptions(token))

// Apply middleware to all events
bot.EventEmitter().Use("*", MyCustomMiddleware())

// Apply middleware only to OnCommand events
bot.EventEmitter().Use(events.OnCommand, OnlyCommandMiddleware())
```

## Creating Custom Middleware

To create a custom middleware, you can use the `eventemitter.MiddlewareFunc` adapter:

```go
func MyCustomMiddleware() eventemitter.Middleware {
    return eventemitter.MiddlewareFunc(func(next eventemitter.Listener) eventemitter.Listener {
        return eventemitter.ListenerFunc(func(ctx context.Context, payload any) error {
            // Logic before the handler
            log.Println("Before handler...")

            err := next.Handle(ctx, payload)

            // Logic after the handler
            log.Println("After handler...")

            return err
        })
    })
}
```

## Global vs Local Middleware

Currently, all middlewares are registered via the `EventEmitter`. The built-in middlewares are registered globally for all events (`*`) during bot initialization.
