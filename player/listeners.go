package player

import (
	"context"

	"github.com/benjamw/golibs/hooks"

	"github.com/benjamw/gogame/model"
	"github.com/benjamw/gogame/session"
)

func listen() {
	hooks.Listen("Register", &RegistrationListener{listenRegister}, 1000)
	hooks.Listen("Login", &LoginListener{listenLogin}, 1000)
	hooks.Listen("Update", &UpdateListener{listenUpdate}, 1000)
}

func listenRegister(ctx context.Context, plyr model.Player, pass string) (bool, error) {

	// do something

	return true, nil
}

func listenLogin(ctx context.Context, s session.Cookier) (bool, error) {

	// do something

	return true, nil
}

func listenUpdate(ctx context.Context, old model.Player, plyr model.Player, pass string) (bool, error) {

	// do something

	return true, nil
}
