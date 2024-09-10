package main

type busMessage struct {
    body string
    subject string
    customProperties map[string]any
}

func newMessage() busMessage {
    props := make(map[string]any)
    props["IsSynthetic"] = "false"
    return busMessage{ body: "testM", subject: "testl", customProperties: props }
}
