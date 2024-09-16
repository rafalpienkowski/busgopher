package main

import (
	"fmt"

)

var activeConnection busConnection
var activeMessage busMessage

func main() {
	fmt.Println("Hello from busgopher!")
    connections, err := loadConnections()
    if err != nil {
        fmt.Println("Can't load connections")
        fmt.Println(err.Error())
        return
    }
    activeConnection = connections[0]

    messages, err := loadMessages()
    if err != nil {
        fmt.Println("Can't load messages")
        fmt.Println(err.Error())
        return
    }
	activeMessage = messages[0]

	fmt.Println("Connecting to '" + activeConnection.Name + "'")
	client := GetClient(activeConnection)

	fmt.Println("Sending a message '" + activeMessage.Body + "' to: '" + activeConnection.Destination + "'")
	SendMessage(activeConnection.Destination, activeMessage, client)

	fmt.Println("Done")
}
