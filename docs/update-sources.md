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

## Webhooks

Webhooks are more efficient for high-traffic bots. Telegram sends an HTTP POST request to your server whenever a new update is available.

### Configuration
To use webhooks, you need a publicly accessible URL (HTTPS is required by Telegram).

```go
// 1. Initialize the webhook source
wh, _ := webhook.New(webhook.NewOptions(
    webhook.WithURL("https://your-bot.example.com/webhook"),
    webhook.WithSecretToken("your-secret-token"), // Recommended
))

// 2. Pass it to the bot
bot, _ := runtime.New(runtime.NewOptions(
    token,
    runtime.WithUpdateSource(wh),
))

// 3. Start your HTTP server
go func() {
    log.Fatal(http.ListenAndServe(":8080", wh))
}()

bot.Run(ctx)
```

### Advantages of Webhooks
- **Real-time:** Updates are received immediately.
- **Resource efficient:** No need for constant polling.

### Secret Token
It's highly recommended to use `WithSecretToken`. This token is sent by Telegram in the `X-Telegram-Bot-Api-Secret-Token` header. The `webhook` package automatically validates this header to ensure requests are actually coming from Telegram.

## Custom Update Sources
You can implement the `UpdateSource` interface yourself if you have a custom way of receiving updates (e.g., from a message queue).
