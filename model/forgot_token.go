package model

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/benjamw/golibs/db"
	"google.golang.org/appengine/datastore"

	"github.com/benjamw/gogame/game"
)

// ForgotToken is an individual forgot password token for a lookup
type ForgotToken struct {
	Base
	PlayerKey *datastore.Key `datastore:"-" json:"-"`
	Token     string         `json:"-"`
	Expires   time.Time      `json:"-"`
}

const forgotTokenEntityType = "ForgotToken"

// EntityType returns the entity type
func (m *ForgotToken) EntityType() string {
	return forgotTokenEntityType
}

// PreSave sets some basic info before continuing on to Save
func (m *ForgotToken) PreSave(ctx context.Context) error {
	if m.GetKey() == nil {
		if m.PlayerKey == nil {
			return errors.New("missing parent player key")
		}
		m.SetIsNew(true)
		m.SetKey(datastore.NewIncompleteKey(ctx, m.EntityType(), m.PlayerKey))
	}

	return nil
}

// PostLoad populates the parent key from the loaded record's key
func (m *ForgotToken) PostLoad(ctx context.Context) error {
	if err := m.Base.PostLoad(ctx); err != nil {
		return err
	}

	m.PlayerKey = m.GetKey().Parent()

	return nil
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

	if game.Now(ctx).After(tokens[0].Expires) {
		myerr = game.NewUserError(nil, http.StatusBadRequest, "That token has expired.")
		return
	}

	*m = tokens[0]

	return
}

// ByToken reads the ForgotToken record with the given token
func (m *ForgotToken) ByToken(ctx context.Context, token string) (myerr error) {
	_, myerr = m.ClearExpired(ctx)
	if myerr != nil {
		return
	}

	var tokens []ForgotToken
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

	*m = tokens[0]

	return
}

func (m *ForgotToken) ClearExisting(ctx context.Context, lookupKey *datastore.Key) (num int, myerr error) {
	var tokens []ForgotToken
	var keys []*datastore.Key
	keys, myerr = datastore.NewQuery(m.EntityType()).
		Ancestor(lookupKey).
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

func (m *ForgotToken) ClearExpired(ctx context.Context) (num int, myerr error) {
	var tokens []ForgotToken
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
