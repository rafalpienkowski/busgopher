package main

type busMessage struct {
    body string
    subject string
}

func newMessage() busMessage {
    return busMessage{ body: "testM", subject: "testl" }
}
