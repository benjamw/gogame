package model

import (
	"context"

	"github.com/benjamw/golibs/db"
	"google.golang.org/appengine/datastore"
)

// Room is a chat room
type Room struct {
	Base
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

const roomEntityType = "Room"

// EntityType returns the entity type
func (m *Room) EntityType() string {
	return roomEntityType
}

// PreSave sets some basic info before continuing on to Save
func (m *Room) PreSave(ctx context.Context) error {
	if m.GetKey() == nil {
		m.SetIsNew(true)
		m.SetKey(makeRoomKey(ctx, m.ID))
	}

	return nil
}

// PostLoad sets some data after the struct is loaded from the db
func (m *Room) PostLoad(ctx context.Context) error {
	if err := m.Base.PostLoad(ctx); err != nil {
		return err
	}

	m.ID = m.key.IntID()

	return nil
}

// ByID loads the room with the given ID
func (m *Room) ByID(ctx context.Context, id int64) (myerr error) {
	key := makeRoomKey(ctx, id)
	room := Room{}
	_, myerr = db.Load(ctx, key, &room)
	if myerr != nil {
		return
	}

	room.SetKey(key)
	if myerr = room.PostLoad(ctx); myerr != nil {
		return
	}

	*m = room

	return
}

func makeRoomKey(ctx context.Context, id int64) *datastore.Key {
	return datastore.NewKey(ctx, roomEntityType, "", id, nil)
}
