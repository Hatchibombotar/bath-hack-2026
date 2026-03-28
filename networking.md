# Client
```go
    // Send message to client
    message, err := json.Marshal(&GenericAction{Action: "blah blah"})
    if err != nil {
        panic(err)
    }

    g.SendMessage(message)
```