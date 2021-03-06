package chat

import (
	"context"
	"strconv"
	"time"

	"github.com/benjamw/gogame/config"
	"github.com/benjamw/golibs/db"
	"google.golang.org/appengine/datastore"

	"github.com/benjamw/gogame/model"
)

// AddChat adds the given message from the given player to the given room
func AddChat(ctx context.Context, roomID, playerID, message string) (chat model.Chat, myerr error) {
	rID, myerr := strconv.ParseInt(roomID, 10, 64)
	if myerr != nil {
		return
	}

	var room model.Room
	if _, myerr = db.LoadInt(ctx, rID, &room); myerr != nil {
		if rID != 0 {
			return
		}

		// lobby wasn't found. make the lobby
		room = model.Room{
			ID:   0,
			Name: config.SiteName + " Lobby",
		}
		if myerr = db.Save(ctx, &room); myerr != nil {
			return
		}
	}

	var player model.Player
	if _, myerr = db.LoadS(ctx, playerID, &player); myerr != nil {
		return
	}

	c := model.Chat{
		RoomKey:   room.GetKey(),
		PlayerKey: player.GetKey(),
		Message:   message,
	}

	if myerr = db.Save(ctx, &c); myerr != nil {
		return
	}

	chat = c

	return
}

// GetChats gets all the chats for the given room
func GetChats(ctx context.Context, roomID string) (room model.Room, chats model.ChatList, myerr error) {
	rID, myerr := strconv.ParseInt(roomID, 10, 64)
	if myerr != nil {
		return
	}

	if _, myerr = db.LoadInt(ctx, rID, &room); myerr != nil {
		return
	}

	if myerr = chats.ByRoomID(ctx, room.GetKey().IntID()); myerr != nil {
		chats = model.ChatList{}

		return
	}

	return
}

// GetChatsAfter gets all the chats for the given room after the given time
func GetChatsAfter(ctx context.Context, roomID string, after time.Time) (room model.Room, chats model.ChatList, myerr error) {
	rID, myerr := strconv.ParseInt(roomID, 10, 64)
	if myerr != nil {
		return
	}

	if _, myerr = db.LoadInt(ctx, rID, &room); myerr != nil {
		return
	}

	if myerr = chats.ByRoomIDAfter(ctx, room.GetKey().IntID(), after); myerr != nil {
		chats = model.ChatList{}

		return
	}

	return
}

// GetMuted returns a list of all the players the given player has muted
func GetMuted(ctx context.Context, playerID string) (mutes model.MuteList, myerr error) {
	var playerKey *datastore.Key
	if playerKey, myerr = datastore.DecodeKey(playerID); myerr != nil {
		return
	}

	var ml model.MuteList
	if _, myerr = ml.ByPlayer(ctx, playerKey); myerr != nil {
		return
	}

	mutes = ml

	return
}

// Mute the given player
func Mute(ctx context.Context, playerID, mutedID string) (mute model.Mute, myerr error) {
	var playerKey, mutedKey *datastore.Key
	if playerKey, myerr = datastore.DecodeKey(playerID); myerr != nil {
		return
	}
	if mutedKey, myerr = datastore.DecodeKey(mutedID); myerr != nil {
		return
	}

	m := model.Mute{
		PlayerKey: playerKey,
		MutedKey:  mutedKey,
	}
	if myerr = db.Save(ctx, &m); myerr != nil {
		return
	}

	mute = m

	return
}

// Unmute the given player
func Unmute(ctx context.Context, playerID, mutedID string) (myerr error) {
	var playerKey, mutedKey *datastore.Key
	if playerKey, myerr = datastore.DecodeKey(playerID); myerr != nil {
		return
	}
	if mutedKey, myerr = datastore.DecodeKey(mutedID); myerr != nil {
		return
	}

	var ml model.MuteList

	ml.ByPlayer(ctx, playerKey)

	for _, v := range ml {
		if v.MutedKey == mutedKey {
			myerr = db.Delete(ctx, &v)
			break
		}
	}

	return
}
