# Update Sources

`tgbotkit-runtime` supports two main ways of receiving updates from Telegram: **Long Polling** and **Webhooks**.

The `UpdateSource` interface (in `interfaces.go`) abstracts these mechanisms:

```go
type UpdateSource interface {
    lifecycle.Lifecycle
    UpdateChan() <-chan client.Update
}
```

## Long Polling (Default)

Long polling is the simplest way to get started. Your bot periodically asks Telegram for new updates.

### Configuration
By default, `runtime.New` will initialize a poller if no other source is provided. You can also configure it manually:

```go
poller, _ := updatepoller.NewPoller(updatepoller.NewOptions(
    apiClient,
    updatepoller.WithOffsetStore(offsetstore.NewInMemoryOffsetStore(0)),
    updatepoller.WithPollingInterval(time.Second),
    updatepoller.WithTimeout(30*time.Second),
    updatepoller.WithRequestTimeout(35*time.Second),
    updatepoller.WithBufferSize(100),
    updatepoller.WithAllowedUpdates([]string{"message", "callback_query"}),
))

bot, _ := runtime.New(runtime.NewOptions(
    token,
    runtime.WithUpdateSource(poller),
))
```

### Offset Storage
The poller needs to keep track of the last update it received. This is done via an `OffsetStore`.
- `InMemoryOffsetStore`: Simplest, but loses state on restart.
- You can implement your own `OffsetStore` (e.g., using Redis or a database) for production use.

The built-in poller advances the offset after every fetched update has been successfully queued into the runtime update channel. It does not wait for application handlers to finish. Treat this as at-most-once-after-queue delivery: production handlers should be idempotent, and production bots should use durable offset storage.

## Webhooks

Webhooks are more efficient for high-traffic bots. Telegram sends an HTTP POST request to your server whenever a new update is available.

### Configuration
To use webhooks, you need a publicly accessible URL (HTTPS is required by Telegram).

```go
// 1. Initialize the webhook source
wh, _ := webhook.New(webhook.NewOptions(
    webhook.WithClient(apiClient),
    webhook.WithUrl("https://your-bot.example.com/webhook"),
    webhook.WithToken("your-secret-token"), // Recommended
    webhook.WithAllowedUpdates([]string{"message", "callback_query"}),
    webhook.WithMaxBodyBytes(1<<20),
))

// 2. Pass it to the bot
bot, _ := runtime.New(runtime.NewOptions(
    token,
    runtime.WithUpdateSource(wh),
))

// 3. Start your HTTP server
server := &http.Server{
    Addr:    ":8080",
    Handler: wh,
}

go func() {
    if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
        log.Fatal(err)
    }
}()

runErr := bot.Run(ctx)

shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
defer cancelShutdown()
if err := server.Shutdown(shutdownCtx); err != nil {
    log.Printf("server shutdown error: %v", err)
}
if runErr != nil {
    log.Fatal(runErr)
}
```

If another deployment component manages `setWebhook`, pass `webhook.WithWebhookRegistrationEnabled(false)` instead of partial registration options.

### Advantages of Webhooks
- **Real-time:** Updates are received immediately.
- **Resource efficient:** No need for constant polling.

### Secret Token
It's highly recommended to use `WithToken`. This token is sent by Telegram in the `X-Telegram-Bot-Api-Secret-Token` header. The `webhook` package automatically validates this header to ensure requests are actually coming from Telegram.

## Custom Update Sources
You can implement the `UpdateSource` interface yourself if you have a custom way of receiving updates (e.g., from a message queue).
