package asb

type InMemoryMessageSender struct {
	Namespace   string
	Destination string
	Message     Message
}

func (messageSender *InMemoryMessageSender) Send(
	namespace string,
	destination string,
	message Message,
) error {

    messageSender.Namespace = namespace
    messageSender.Destination = destination
    messageSender.Message = message

	return nil
}
