package model

import (
	"context"

	"github.com/benjamw/golibs/db"
	"google.golang.org/appengine/datastore"

	"github.com/benjamw/gogame/game"
)

// ForgotToken is an individual forgot password token for a player
type ForgotToken struct {
	Token
}

const forgotTokenEntityType = "ForgotToken"

// EntityType returns the entity type
func (m *ForgotToken) EntityType() string {
	return forgotTokenEntityType
}

// ByEmail reads the ForgotToken record with the given email
func (m *ForgotToken) ByEmail(ctx context.Context, email string) (myerr error) {
	_, myerr = m.ClearExpired(ctx)
	if myerr != nil {
		return
	}

	var player Player
	if err := player.ByEmail(ctx, email); err != nil {
		myerr = &db.UnfoundObjectError{
			EntityType: player.EntityType(),
			Key:        "email",
			Value:      email,
		}
		return
	}

	var tokens []ForgotToken
	var keys []*datastore.Key
	keys, myerr = datastore.NewQuery(m.EntityType()).
		Ancestor(player.GetKey()).
		GetAll(ctx, &tokens)
	if myerr != nil {
		return
	}

	if len(tokens) == 0 {
		myerr = &db.UnfoundObjectError{
			EntityType: m.EntityType(),
			Key:        "email",
			Value:      email,
		}
		return
	}

	if 1 < len(tokens) {
		myerr = &game.MultipleObjectError{
			EntityType: player.EntityType(),
			Key:        "email",
			Value:      email,
		}
		return
	}

	tokens[0].SetKey(keys[0])
	if myerr = tokens[0].PostLoad(ctx); myerr != nil {
		return
	}

	*m = tokens[0]

	return
}
