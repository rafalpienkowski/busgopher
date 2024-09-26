package main

type busMessage struct {
    Name             string         `json:"name"`
	Body             string         `json:"body"`
	Subject          string         `json:"subject"`
	CustomProperties map[string]any `json:"customProperties"`
}
