package firebase

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"github.com/arifai/zenith/cmd/wire/logger"
	"go.uber.org/zap"
)

type MessagingService struct{ *Messaging }

var log = logger.ProvideLogger()

func NewMessagingService(messaging *Messaging) *MessagingService {
	return &MessagingService{messaging}
}

func (m *MessagingService) SendMessage(data map[string]string, token, title, body string) error {
	message := &messaging.Message{
		Token:        token,
		Data:         data,
		Notification: &messaging.Notification{Title: title, Body: body},
	}

	response, err := m.Client.Send(context.Background(), message)
	if err != nil {
		return err
	}

	log.Info("successful send push notification", zap.String("response", response))

	return nil
}
