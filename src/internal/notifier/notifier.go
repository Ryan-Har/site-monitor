package notifier

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Notifier interface {
	SendTest() error
	SendMessage(msg string) error
}

type EmailNotifier struct {
	Email string
}

type WebhookNotifier struct {
	Url         string
	WebhookType string //discord or slack
}

// func NewEmailNotifier(opts ...func(*EmailNotifier)) Notifier {
// 	notifier := &EmailNotifier{}
// 	for _, opt := range opts {
// 		opt(notifier)
// 	}
// 	return notifier
// }

// func WithEmail(email string) func(*EmailNotifier) {
// 	return func(en *EmailNotifier) {
// 		en.Email = email
// 	}
// }

// func (en *EmailNotifier) SendTest() error {
// 	slog.Info("test", "email", en.Email)
// 	return nil
// }

func NewDiscordNotifier(opts ...func(*WebhookNotifier)) Notifier {
	notifier := &WebhookNotifier{
		WebhookType: "discord",
	}
	for _, opt := range opts {
		opt(notifier)
	}
	return notifier
}

func NewSlackNotifier(opts ...func(*WebhookNotifier)) Notifier {
	notifier := &WebhookNotifier{
		WebhookType: "slack",
	}
	for _, opt := range opts {
		opt(notifier)
	}
	return notifier
}

func WithUrl(url string) func(*WebhookNotifier) {
	return func(wn *WebhookNotifier) {
		wn.Url = url
	}
}

func (wn *WebhookNotifier) SendTest() error {
	if err := wn.SendMessage("Test Message"); err != nil {
		return err
	}
	return nil
}

type Message interface{}

type DiscordMessage struct {
	Content string `json:"content"`
}

type SlackMessage struct {
	Content string `json:"text"`
}

func getMessageStruct(webhookType string, msg string) Message {
	switch webhookType {
	case "discord":
		return DiscordMessage{
			Content: msg,
		}
	case "slack":
		return SlackMessage{
			Content: msg,
		}
	}
	return nil
}

func (wn *WebhookNotifier) SendMessage(msg string) error {
	msgStruct := getMessageStruct(wn.WebhookType, msg)
	jsonData, err := json.Marshal(msgStruct)
	if err != nil {
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, wn.Url, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	slog.Info("Webhook sent", "Status code", resp.StatusCode)

	return nil
}
