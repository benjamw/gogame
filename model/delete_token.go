package model

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/benjamw/golibs/db"
	"google.golang.org/appengine/datastore"

	"github.com/benjamw/gogame/game"
)

// DeleteToken is an individual delete token for a player
type DeleteToken struct {
	Base
	PlayerKey *datastore.Key `datastore:"-" json:"-"`
	Token     string         `json:"-"`
	Expires   time.Time      `json:"-"`
}

const deleteTokenEntityType = "DeleteToken"

// EntityType returns the entity type
func (m *DeleteToken) EntityType() string {
	return deleteTokenEntityType
}

// PreSave sets some basic info before continuing on to Save
func (m *DeleteToken) PreSave(ctx context.Context) error {
	if m.GetKey() == nil {
		if m.PlayerKey == nil {
			return errors.New("missing parent player key")
		}
		m.SetIsNew(true)
		m.SetKey(datastore.NewIncompleteKey(ctx, m.EntityType(), m.PlayerKey))
	} else {
		return fmt.Errorf("cannot edit existing token")
	}

	if m.Expires.IsZero() {
		m.Expires = game.Now(ctx).Add(24 * time.Hour) // 1 day
	}

	return nil
}

// PostLoad populates the parent key from the loaded record's key
func (m *DeleteToken) PostLoad(ctx context.Context) error {
	if err := m.Base.PostLoad(ctx); err != nil {
		return err
	}

	m.PlayerKey = m.GetKey().Parent()

	return nil
}

// ByToken reads the DeleteToken record with the given token
func (m *DeleteToken) ByToken(ctx context.Context, token string) (myerr error) {
	var tokens []DeleteToken
	var keys []*datastore.Key
	keys, myerr = datastore.NewQuery(m.EntityType()).
		Filter("Token =", token).
		GetAll(ctx, &tokens)
	if myerr != nil {
		return
	}

	if len(tokens) == 0 {
		myerr = &db.UnfoundObjectError{
			EntityType: m.EntityType(),
			Key:        "token",
			Value:      token,
			Err:        nil,
		}
		return
	}

	if 1 < len(tokens) {
		myerr = &game.MultipleObjectError{
			EntityType: m.EntityType(),
			Key:        "token",
			Value:      token,
		}
		return
	}

	tokens[0].SetKey(keys[0])
	if myerr = tokens[0].PostLoad(ctx); myerr != nil {
		return
	}

	if game.Now(ctx).After(tokens[0].Expires) {
		myerr = game.NewUserError(nil, http.StatusBadRequest, "That token has expired.")
		return
	}

	_, myerr = m.ClearExpired(ctx)
	if myerr != nil {
		return
	}

	*m = tokens[0]

	return
}

func (m *DeleteToken) ClearExisting(ctx context.Context, playerKey *datastore.Key) (num int, myerr error) {
	var tokens []DeleteToken
	var keys []*datastore.Key
	keys, myerr = datastore.NewQuery(m.EntityType()).
		Ancestor(playerKey).
		GetAll(ctx, &tokens)
	if myerr != nil {
		return
	}

	num = 0
	for k := range tokens {
		tokens[k].SetKey(keys[k])
		if myerr = tokens[k].PostLoad(ctx); myerr != nil {
			return
		}

		if myerr = db.Delete(ctx, &tokens[k]); myerr != nil {
			return
		}

		num++
	}

	return
}

func (m *DeleteToken) ClearExpired(ctx context.Context) (num int, myerr error) {
	var tokens []DeleteToken
	var keys []*datastore.Key
	keys, myerr = datastore.NewQuery(m.EntityType()).
		Filter("Expires <", game.Now(ctx)).
		GetAll(ctx, &tokens)
	if myerr != nil {
		return
	}

	num = 0
	for k := range tokens {
		tokens[k].SetKey(keys[k])
		if myerr = tokens[k].PostLoad(ctx); myerr != nil {
			return
		}

		if myerr = db.Delete(ctx, &tokens[k]); myerr != nil {
			return
		}

		num++
	}

	return
}
