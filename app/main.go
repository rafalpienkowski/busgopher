package main

import (
	"fmt"

)

func main() {
    connections, err := loadConnections()
    if err != nil {
        fmt.Println("Can't load connections")
        fmt.Println(err.Error())
        return
    }
    activeConnection := connections[0]

    messages, err := loadMessages()
    if err != nil {
        fmt.Println("Can't load messages")
        fmt.Println(err.Error())
        return
    }
	message := messages[0]

	fmt.Println("Hello from busgopher!")
	fmt.Println("Connecting to '" + activeConnection.Name + "'")
	client := GetClient(activeConnection)

	fmt.Println("Sending a message '" + message.Body + "' to: '" + activeConnection.Destination + "'")
	SendMessage(activeConnection.Destination, message, client)

	fmt.Println("Done")
}
