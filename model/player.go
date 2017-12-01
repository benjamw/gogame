package model

import (
	"context"
	"time"

	"github.com/benjamw/golibs/db"
	"google.golang.org/appengine/datastore"

	"github.com/benjamw/gogame/game"
)

// Player is a user entity
type Player struct {
	Base
	Username     string        `json:"username"`
	Email        string        `json:"email"`
	PasswordHash string        `datastore:",noindex" json:"-"`
	Timezone     time.Location `json:"-"`
	IsAdmin      bool          `json:"is_admin"`
	Created      time.Time     `json:"-"`
	Approved     time.Time     `json:"-"`
}

const playerEntityType = "Player"

// EntityType returns the entity type
func (m *Player) EntityType() string {
	return playerEntityType
}

// PreSave sets some basic info before continuing on to Save
func (m *Player) PreSave(ctx context.Context) error {
	if m.GetKey() == nil {
		m.SetIsNew(true)
		m.SetKey(datastore.NewIncompleteKey(ctx, m.EntityType(), nil))
	}

	if m.Created.IsZero() {
		m.Created = time.Now()
	}

	return nil
}

func (m *Player) ByEmail(ctx context.Context, email string) (myerr error) {
	var people []Player
	var keys []*datastore.Key
	keys, myerr = datastore.NewQuery(m.EntityType()).
		Filter("Email =", email).
		GetAll(ctx, &people)
	if myerr != nil {
		return
	}

	if len(people) == 0 {
		myerr = &db.UnfoundObjectError{
			EntityType: m.EntityType(),
			Key:        "email",
			Value:      email,
		}
		return
	}

	if 1 < len(people) {
		myerr = &game.MultipleObjectError{
			EntityType: m.EntityType(),
			Key:        "email",
			Value:      email,
		}
		return
	}

	people[0].SetKey(keys[0])
	if myerr = people[0].PostLoad(ctx); myerr != nil {
		return
	}

	*m = people[0]

	return nil
}
