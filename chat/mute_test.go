package chat

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"google.golang.org/appengine/datastore"

	"github.com/benjamw/gogame/model"
	"github.com/benjamw/golibs/db"
	"github.com/benjamw/golibs/test"
)

// MODEL TESTS

// TestMain in chat_test.go

func TestMuteEntityType(t *testing.T) {
	var m model.Chat
	if "Chat" != m.EntityType() {
		t.Fatalf("Chat.EntityType() returned '%s', wanted 'Chat'", m.EntityType())
	}
}

func TestMutePreSave(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var c model.Chat

	// test with no player
	if err = c.PreSave(ctx); err == nil {
		t.Fatal("Mute.PreSave did not throw an error for missing Parent Player Key when one should not exist.")
	}

	player := createRandPlayer(ctx, t)

	// test with no player
	c.PlayerKey = player.GetKey()
	if err = c.PreSave(ctx); err != nil {
		spew.Dump(err)
		t.Fatal("Mute.PreSave threw an error for missing Parent Player Key when one should exist.")
	}

	key := c.GetKey()
	if key == nil {
		t.Fatal("Mute.PreSave did not create a datastore key.")
	}
}

// CONTROLLER TESTS

// TODO: add these
func TestMuteController(t *testing.T) {
	t.Fatal("chat controller tests not complete")
}

func TestGetMuted(t *testing.T) {

}

func TestMute(t *testing.T) {

}

func TestUnmute(t *testing.T) {

}

// HELPER FUNCTIONS

func createFullMute(ctx context.Context, t *testing.T, playerKey, mutedKey *datastore.Key) model.Mute {
	file, line, funct := test.GetCaller()

	thing := model.Mute{
		PlayerKey: playerKey,
		MutedKey:  mutedKey,
	}

	if err := db.Save(ctx, &thing); err != nil {
		t.Fatalf("Could not save the test Mute. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
	}

	// perform a Get to force the key to be applied so it's available in queries
	var get model.Mute // for use in the forcing Get, not actual data
	err := datastore.Get(ctx, thing.GetKey(), &get)
	if err != nil {
		t.Fatalf("Could not get the test Mute. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
	}

	return thing
}

func createMute(ctx context.Context, t *testing.T, player, muted model.Player) model.Mute {
	return createFullMute(ctx, t, player.GetKey(), muted.GetKey())
}

func createRandMute(ctx context.Context, t *testing.T) model.Mute {
	player := createRandPlayer(ctx, t)
	muted := createRandPlayer(ctx, t)

	return createMute(ctx, t, player, muted)
}
