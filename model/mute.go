package model

import (
	"context"

	"github.com/benjamw/golibs/db"
	"google.golang.org/appengine/datastore"
)

// Mute is a muted player entry
type Mute struct {
	Base
	PlayerKey *datastore.Key `datastore:"-" json:"-"`
	MutedKey  *datastore.Key `json:"-"`
}

// MuteList is a list of mutes
type MuteList []Mute

const muteEntityType = "Mute"

// EntityType returns the entity type
func (m *Mute) EntityType() string {
	return muteEntityType
}

// PreSave sets some basic info before continuing on to Save
func (m *Mute) PreSave(ctx context.Context) error {
	if m.GetKey() == nil {
		if m.PlayerKey == nil {
			return &db.MissingParentKeyError{}
		}
		m.SetIsNew(true)
		m.SetKey(datastore.NewIncompleteKey(ctx, m.EntityType(), m.PlayerKey))
	}

	if m.MutedKey == nil {
		return &db.MissingRequiredError{"MutedKey"}
	}

	return nil
}

// PostLoad populates the parent key from the loaded record's key
func (m *Mute) PostLoad(ctx context.Context) error {
	if err := m.Base.PostLoad(ctx); err != nil {
		return err
	}

	m.PlayerKey = m.key.Parent()

	return nil
}

// ByPlayer loads the mutes with the given parent player key
func (l *MuteList) ByPlayer(ctx context.Context, playerKey *datastore.Key) (num int, myerr error) {
	var mutes []Mute
	var keys []*datastore.Key
	keys, myerr = datastore.NewQuery(new(Mute).EntityType()).
		Ancestor(playerKey).
		GetAll(ctx, &mutes)
	if myerr != nil {
		return
	}

	num = 0
	for k := range mutes {
		mutes[k].SetKey(keys[k])
		if myerr = mutes[k].PostLoad(ctx); myerr != nil {
			return
		}

		num++
	}

	*l = mutes

	return
}
