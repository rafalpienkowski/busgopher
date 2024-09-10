package main

import (
	"fmt"
)

func main() {
	connection := connection{
		namespace:   "***",
		name:        "test",
		destination: "test-rpe",
	}
	message := newMessage()

	fmt.Println("Hello from busgopher!")
	fmt.Println("Connecting to '" + connection.name + "'")
	client := GetClient(connection)

	fmt.Println("Sending a message '" + message.body + "' to: '" + connection.destination + "'")
	SendMessage(connection.destination, message, client)

	fmt.Println("Done")
}
