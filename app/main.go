package main

import (
	"fmt"
)

func main() {
    namespace := "***"
	queue := "test-rpe"
	message := newMessage()

	fmt.Println("Hello from busgopher!")
	fmt.Println("Connecting to '" + namespace + "'")
    client := GetClient(namespace)

	fmt.Println("Sending a message '" + message.body + "' to queue: '" + queue + "'")
    SendMessage(queue, message, client)

	fmt.Println("Done")
}
