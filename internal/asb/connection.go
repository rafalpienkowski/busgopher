package asb

type Connection struct {
	Name         string   `json:"name"`
	Namespace    string   `json:"namespace"`
	Destinations []string `json:"destinations"`
}

type MessageSender interface {
    Send(namespace string, destination string, message Message) error
}
