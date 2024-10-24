package asb

import (
	"bytes"
	"encoding/json"
	"text/template"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Body string `json:"body"`

	//Broker properties
	CorrelationID string `json:"correlationId"`
	MessageID     string `json:"messageId"`
	ReplayTo      string `json:"replyTo"`
	Subject       string `json:"subject"`

	CustomProperties map[string]any `json:"customProperties"`
}

func (msg *Message) Print() string {

	prettyMsgBytes, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(prettyMsgBytes)
}

func (msg *Message) TransformBody() (string, error) {

	t := template.Must(template.New("example").Funcs(template.FuncMap{
		"utcNow": func() string { return time.Now().UTC().Format(time.RFC3339) },
		"utcNowPlus": func(minutes int) string {
			return time.Now().UTC().Add(time.Duration(minutes) * time.Minute).Format(time.RFC3339)
		},
		"generateUUID": func() string { return uuid.New().String() },
	}).Parse(msg.Body))

	var output bytes.Buffer
	err := t.Execute(&output, nil)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
