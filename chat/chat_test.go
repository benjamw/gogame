package chat

import (
	"context"
	"os"
	"testing"
	"time"

	"google.golang.org/appengine/datastore"

	"github.com/benjamw/gogame/model"
	"github.com/benjamw/golibs/db"
	"github.com/benjamw/golibs/password"
	"github.com/benjamw/golibs/random"
	"github.com/benjamw/golibs/test"
)

func TestMain(m *testing.M) {
	test.InitCtx()
	runVal := m.Run()
	test.ReleaseCtx()
	os.Exit(runVal)
}

// MODEL TESTS

func TestChatEntityType(t *testing.T) {
	var m model.Chat
	if "Chat" != m.EntityType() {
		t.Fatalf("Chat.EntityType() returned '%s', wanted 'Chat'", m.EntityType())
	}
}

func TestChatPreSave(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var c model.Chat

	// test with no room
	if err = c.PreSave(ctx); err == nil {
		t.Fatal("Chat.PreSave did not throw an error for missing Parent Room Key when one should not exist.")
	}

	room := createRandRoom(ctx, t)

	// test with no player
	c.PlayerKey = room.GetKey()
	if err = c.PreSave(ctx); err != nil {
		t.Fatal("Chat.PreSave threw an error for missing Parent Room Key when one should exist.")
	}

	key := c.GetKey()
	if key == nil {
		t.Fatal("Chat.PreSave did not create a datastore key.")
	}
}

// CONTROLLER TESTS

// TODO: add these
func TestController(t *testing.T) {
	t.Fatal("chat controller tests not complete")
}

func TestAddChat(t *testing.T) {

}

func TestGetChats(t *testing.T) {

}

func TestGetChatsAfter(t *testing.T) {

}

// HELPER FUNCTIONS

func createFullPlayer(ctx context.Context, t *testing.T, username string, email string, passwrd string) model.Player {
	file, line, funct := test.GetCaller()

	thing := model.Player{
		Username:     username,
		Email:        email,
		PasswordHash: password.Encode(passwrd),
	}
	if err := db.Save(ctx, &thing); err != nil {
		t.Fatalf("Could not save the test Player. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
	}

	// perform a Get to force the key to be applied so it's available in queries
	var get model.Player // for use in the forcing Get, not actual data
	err := datastore.Get(ctx, thing.GetKey(), &get)
	if err != nil {
		t.Fatalf("Could not get the test Player. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
	}

	return thing
}

func createFullRoom(ctx context.Context, t *testing.T, id int64, name string) model.Room {
	file, line, funct := test.GetCaller()

	thing := model.Room{
		ID:   id,
		Name: name,
	}
	if err := db.Save(ctx, &thing); err != nil {
		t.Fatalf("Could not save the test Room. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
	}

	// perform a Get to force the key to be applied so it's available in queries
	var get model.Room // for use in the forcing Get, not actual data
	err := datastore.Get(ctx, thing.GetKey(), &get)
	if err != nil {
		t.Fatalf("Could not get the test Room. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
	}

	return thing
}

func createFullChat(ctx context.Context, t *testing.T, roomKey, playerKey *datastore.Key, message string, created time.Time) model.Chat {
	file, line, funct := test.GetCaller()

	thing := model.Chat{
		RoomKey:   roomKey,
		PlayerKey: playerKey,
		Message:   message,
		Created:   created,
	}

	if err := db.Save(ctx, &thing); err != nil {
		t.Fatalf("Could not save the test Chat. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
	}

	// perform a Get to force the key to be applied so it's available in queries
	var get model.Chat // for use in the forcing Get, not actual data
	err := datastore.Get(ctx, thing.GetKey(), &get)
	if err != nil {
		t.Fatalf("Could not get the test Chat. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
	}

	return thing
}

func createRandPlayer(ctx context.Context, t *testing.T) model.Player {
	username := random.Stringn(64)
	email := random.Email()
	passwrd := random.Stringn(64)

	return createFullPlayer(ctx, t, username, email, passwrd)
}

func createRandRoom(ctx context.Context, t *testing.T) model.Room {
	return createFullRoom(ctx, t, random.Int63(), random.Stringn(10))
}

func createChat(ctx context.Context, t *testing.T, room model.Room, player model.Player, message string) model.Chat {
	return createFullChat(ctx, t, room.GetKey(), player.GetKey(), message, time.Now())
}

func createRandChat(ctx context.Context, t *testing.T) model.Chat {
	room := createRandRoom(ctx, t)
	player := createRandPlayer(ctx, t)
	message := random.String()

	return createChat(ctx, t, room, player, message)
}
