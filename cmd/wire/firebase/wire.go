//go:build wireinject

package firebase

import (
	"github.com/arifai/zenith/pkg/firebase"
	"github.com/google/wire"
)

func ProvideFirebase(file string) (*firebase.Messaging, error) {
	wire.Build(firebase.New)
	return &firebase.Messaging{}, nil
}

func ProvideFirebaseMessagingService(file string) (*firebase.MessagingService, error) {
	wire.Build(ProvideFirebase, firebase.NewMessagingService)
	return &firebase.MessagingService{}, nil
}
