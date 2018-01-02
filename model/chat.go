package model

import (
	"context"
	"time"

	"github.com/benjamw/golibs/db"
	"google.golang.org/appengine/datastore"

	"github.com/benjamw/gogame/game"
)

type Chat struct {
	Base
	RoomKey   *datastore.Key
	PlayerKey *datastore.Key
	Message   string
	Created   time.Time
}

const chatEntityType = "Chat"

// EntityType returns the entity type
func (m *Chat) EntityType() string {
	return chatEntityType
}

// PreSave sets some basic info before continuing on to Save
func (m *Chat) PreSave(ctx context.Context) error {
	if m.GetKey() == nil {
		if m.RoomKey == nil {
			return &db.MissingParentKeyError{}
		}

		m.SetIsNew(true)
		m.SetKey(datastore.NewIncompleteKey(ctx, m.EntityType(), m.RoomKey))
	}

	if m.PlayerKey == nil {
		return &db.MissingRequiredError{"PlayerKey"}
	}

	if m.Message == "" {
		return &db.MissingRequiredError{"Message"}
	}

	if m.Created.IsZero() {
		m.Created = game.Now(ctx)
	}

	return nil
}

// PostLoad sets some data after the struct is loaded from the db
func (m *Chat) PostLoad(ctx context.Context) error {
	if err := m.Base.PostLoad(ctx); err != nil {
		return err
	}

	m.RoomKey = m.key.Parent()

	return nil
}

// ChatList is a slice of related Chats
type ChatList []Chat

// ByRoomID loads all the chats for the room with the given ID
func (l *ChatList) ByRoomID(ctx context.Context, id int64) (myerr error) {
	key := makeRoomKey(ctx, id)

	var chats []Chat
	var keys []*datastore.Key
	keys, myerr = datastore.NewQuery(chatEntityType).
		Ancestor(key).
		Order("-Created"). // DESC
		GetAll(ctx, &chats)
	if myerr != nil {
		return
	}

	for k := range keys {
		chats[k].SetKey(keys[k])
		if myerr = chats[k].PostLoad(ctx); myerr != nil {
			return
		}
	}

	*l = chats

	return
}

// ByRoomIDAfter loads all the chats for the room with the given ID that came in after the given time
func (l *ChatList) ByRoomIDAfter(ctx context.Context, id int64, after time.Time) (myerr error) {
	key := makeRoomKey(ctx, id)

	var chats []Chat
	var keys []*datastore.Key
	keys, myerr = datastore.NewQuery(chatEntityType).
		Ancestor(key).
		Filter("Created >=", after).
		Order("-Created"). // DESC
		GetAll(ctx, &chats)
	if myerr != nil {
		return
	}

	for k := range keys {
		chats[k].SetKey(keys[k])
		if myerr = chats[k].PostLoad(ctx); myerr != nil {
			return
		}
	}

	*l = chats

	return
}
