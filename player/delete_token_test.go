package player

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/benjamw/golibs/db"
	"github.com/benjamw/golibs/password"
	"github.com/benjamw/golibs/random"
	"github.com/benjamw/golibs/test"
	"google.golang.org/appengine/datastore"

	"github.com/benjamw/gogame/config"
	"github.com/benjamw/gogame/game"
	"github.com/benjamw/gogame/model"
)

var (
	email  = "bwelker@daz3d.com"                           // set this to a valid email address for the send email test
	expiry = time.Now().Add(time.Hour * time.Duration(24)) // 24 hours from now
)

func TestMain(m *testing.M) {
	test.InitCtx()

	config.GenRoot()
	config.RootURL = "UNIT_TESTING"

	runVal := m.Run()
	test.ReleaseCtx()
	os.Exit(runVal)
}

// MODEL TESTS

func TestEntityType(t *testing.T) {
	var dt model.DeleteToken
	if "DeleteToken" != dt.EntityType() {
		t.Fatal("DeleteToken.EntityType returned incorrect EntityType.")
	}
}

func TestPreSave(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var dt model.DeleteToken

	// test with no player
	if err = dt.PreSave(ctx); err == nil {
		t.Fatal("DeleteToken.PreSave did not throw an error for missing Parent Player Key when one should not exist.")
	}

	player := createPlayer(ctx, t)

	// test proper
	dt.PlayerKey = player.GetKey()
	if err = dt.PreSave(ctx); err != nil {
		t.Fatal("DeleteToken.PreSave threw an error for missing Parent Player Key when one should exist.")
	}

	key := dt.GetKey()
	if key == nil {
		t.Fatal("DeleteToken.PreSave did not create a datastore key.")
	}
}

func TestPostLoad(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var dt model.DeleteToken

	testToken := createRandDeleteToken(ctx, t)

	// test with no key
	err = dt.PostLoad(ctx)
	if _, ok := err.(*db.MissingKeyError); !ok {
		t.Fatalf("DeleteToken.PostLoad threw the wrong error for missing key. Error: %v", err)
	}

	dt.SetKey(testToken.GetKey())

	// test proper
	err = dt.PostLoad(ctx)
	if err != nil {
		t.Fatalf("DeleteToken.PostLoad threw an error: %v", err)
	}
	if dt.PlayerKey == nil {
		t.Fatal("DeleteToken.PostLoad did not attach the Parent Player Key.")
	}
}

func TestClearExisting(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var dt model.DeleteToken

	token := random.Stringn(64)

	player := createPlayer(ctx, t)

	// create the proper tokens
	numTokens := 4
	for i := 0; i < numTokens; i++ {
		createDeleteToken(ctx, t, player, token)
	}

	// create some other tokens
	player2 := createRandPlayer(ctx, t)
	numOtherTokens := 3
	for i := 0; i < numOtherTokens; i++ {
		createDeleteToken(ctx, t, player2, token)
	}

	// do a manual selection here because .ByValue fails if more than one found
	dtList := make([]model.DeleteToken, 0)
	keys, err := datastore.NewQuery(dt.EntityType()).
		Filter("Value =", token).
		KeysOnly().
		GetAll(ctx, &dtList)
	if err != nil {
		t.Fatalf("TestClearExisting first query threw an error. Error: %v", err)
	}
	if len(keys) != numTokens+numOtherTokens {
		t.Fatal("TestClearExisting did not create the proper number of DeleteTokens.")
	}

	// test proper
	num, err := dt.ClearExisting(ctx, player.GetKey())
	if err != nil {
		t.Fatalf("DeleteToken.ClearExisting threw an error. Error: %v", err)
	}
	if num != numTokens {
		t.Fatalf("DeleteToken.ClearExisting reported the wrong number of tokens deleted: %d. Wanted: %d", num, numTokens)
	}

	// run a get on all the keys to force an update
	var getDeleteToken model.DeleteToken // for use in the forcing Get, not actual data
	for _, v := range keys {
		datastore.Get(ctx, v, &getDeleteToken)
	}

	// do a manual selection here because .ByValue fails if more than one found
	dtList = make([]model.DeleteToken, 0)
	keys, err = datastore.NewQuery(dt.EntityType()).
		Filter("Value =", token).
		KeysOnly().
		GetAll(ctx, &dtList)
	if err != nil {
		t.Fatalf("TestClearExisting second query threw an error. Error: %v", err)
	}
	if len(keys) != numOtherTokens {
		t.Fatal("DeleteToken.ClearExisting did not delete the correct number of DeleteTokens.")
	}
}

func TestClearExpired(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var dt model.DeleteToken

	player := createPlayer(ctx, t)

	// create the proper tokens
	numTokens := 4
	for i := 0; i < numTokens; i++ {
		createPastToken(ctx, t, player)
	}

	// create some other tokens
	numOtherTokens := 3
	for i := 0; i < numOtherTokens; i++ {
		createRandDeleteTokenForPlayer(ctx, t, player)
	}

	// do a manual selection here because .ByValue fails if more than one found
	dtList := make([]model.DeleteToken, 0)
	keys, err := datastore.NewQuery(dt.EntityType()).
		Ancestor(player.GetKey()).
		KeysOnly().
		GetAll(ctx, &dtList)
	if err != nil {
		t.Fatalf("TestClearExpired first query threw an error. Error: %v", err)
	}
	if len(keys) != numTokens+numOtherTokens {
		t.Fatalf("TestClearExpired did not create the correct number of DeleteTokens: %d. Wanted: %d.", len(keys), numTokens+numOtherTokens)
	}

	// test proper
	num, err := dt.ClearExpired(ctx)
	if err != nil {
		t.Fatalf("DeleteToken.ClearExpired threw an error. Error: %v", err)
	}
	if num != numTokens {
		t.Fatalf("DeleteToken.ClearExpired reported the wrong number of tokens deleted: %d. Wanted: %d", num, numTokens)
	}

	// run a get on all the keys to force an update
	var getDeleteToken model.DeleteToken // for use in the forcing Get, not actual data
	for _, v := range keys {
		datastore.Get(ctx, v, &getDeleteToken)
	}

	// do a manual selection here because .ByValue fails if more than one found
	dtList = make([]model.DeleteToken, 0)
	keys, err = datastore.NewQuery(dt.EntityType()).
		Ancestor(player.GetKey()).
		KeysOnly().
		GetAll(ctx, &dtList)
	if err != nil {
		t.Fatalf("TestClearExpired second query threw an error. Error: %v", err)
	}
	if len(keys) != numOtherTokens {
		t.Fatalf("DeleteToken.ClearExpired did not delete the correct number of DeleteTokens: %d. Wanted: %d.", numTokens+numOtherTokens-len(keys), numTokens)
	}
}

func TestByToken(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var dt model.DeleteToken

	token := random.Stringn(64)

	// test with no tokens
	err = dt.ByValue(ctx, token)
	if err == nil {
		t.Fatal("DeleteToken.ByValue returned a DeleteToken when none should exist.")
	}
	if _, ok := err.(*db.UnfoundObjectError); !ok {
		t.Fatalf("DeleteToken.ByValue did not throw an unfound object error when none should exist: Type: %T; Error: %v", err, err)
	}

	player := createPlayer(ctx, t)

	createPastToken(ctx, t, player)

	// test with expired token
	err = dt.ByValue(ctx, email)
	if err == nil {
		t.Fatal("DeleteToken.ByValue returned a DeleteToken when none should exist (Expired).")
	}
	if _, ok := err.(*db.UnfoundObjectError); !ok {
		t.Fatalf("DeleteToken.ByValue did not throw an unfound object error with an expired token: Type: %T; Error: %v", err, err)
	}

	createDeleteToken(ctx, t, player, token)

	// test proper
	err = dt.ByValue(ctx, token)
	if err != nil {
		t.Fatalf("DeleteToken.ByValue threw an error: %v", err)
	}
	if token != dt.Value {
		t.Fatal("DeleteToken.ByValue returned the wrong DeleteToken.")
	}

	createDeleteToken(ctx, t, player, token)

	// test with too many tokens
	var dt2 model.DeleteToken
	err = dt2.ByValue(ctx, token)
	if _, ok := err.(*game.MultipleObjectError); !ok {
		t.Fatalf("DeleteToken.ByValue returned the wrong error when two matched DeleteTokens existed. Got: %T '%v'", err, err)
	}
}

// CONTROLLER TESTS

// TODO: add these
func TestDeleteTokenController(t *testing.T) {
	t.Fatal("delete_token_test not complete")
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

func createFullDeleteToken(ctx context.Context, t *testing.T, playerKey *datastore.Key, expires time.Time, token string) model.DeleteToken {
	file, line, funct := test.GetCaller()

	thing := model.DeleteToken{
		PlayerKey: playerKey,
		Expires:   expires,
		Value:     token,
	}

	if err := db.Save(ctx, &thing); err != nil {
		t.Fatalf("Could not save the test DeleteToken. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
	}

	// perform a Get to force the key to be applied so it's available in queries
	var get model.DeleteToken // for use in the forcing Get, not actual data
	err := datastore.Get(ctx, thing.GetKey(), &get)
	if err != nil {
		t.Fatalf("Could not get the test DeleteToken. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
	}

	return thing
}

func createPlayer(ctx context.Context, t *testing.T) model.Player {
	username := random.Stringn(64)
	passwrd := random.Stringn(64)

	return createFullPlayer(ctx, t, username, email, passwrd)
}

func createRandPlayer(ctx context.Context, t *testing.T) model.Player {
	username := random.Stringn(64)
	email := random.Email()
	passwrd := random.Stringn(64)

	return createFullPlayer(ctx, t, username, email, passwrd)
}

func createDeleteToken(ctx context.Context, t *testing.T, player model.Player, token string) model.DeleteToken {
	return createFullDeleteToken(ctx, t, player.GetKey(), expiry, token)
}

func createPastToken(ctx context.Context, t *testing.T, player model.Player) model.DeleteToken {
	expires := expiry.AddDate(-1, 0, 0) // 1 year in the past
	token := random.Stringn(64)

	return createFullDeleteToken(ctx, t, player.GetKey(), expires, token)
}

func createRandDeleteTokenForPlayer(ctx context.Context, t *testing.T, player model.Player) model.DeleteToken {
	token := random.Stringnt(64, random.ALPHANUMERIC)

	return createFullDeleteToken(ctx, t, player.GetKey(), expiry, token)
}

func createRandDeleteToken(ctx context.Context, t *testing.T) model.DeleteToken {
	player := createPlayer(ctx, t)
	token := random.Stringnt(64, random.ALPHANUMERIC)

	return createFullDeleteToken(ctx, t, player.GetKey(), expiry, token)
}
