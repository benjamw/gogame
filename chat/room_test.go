package chat

import (
	"testing"

	"github.com/benjamw/gogame/model"
)

// MODEL TESTS

func TestRoomEntityType(t *testing.T) {
	var m model.Chat
	if "Chat" != m.EntityType() {
		t.Fatalf("Chat.EntityType() returned '%s', wanted 'Chat'", m.EntityType())
	}
}

// CONTROLLER TESTS

// TODO: add these
func TestRoomController(t *testing.T) {
	t.Fatal("chat controller tests not complete")
}

// HELPER FUNCTIONS
