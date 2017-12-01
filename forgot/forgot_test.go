package forgot

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
	email  = "benjam@iohelix.net"                          // set this to a valid email address for the send email test
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
	var ft model.ForgotToken
	if "ForgotToken" != ft.EntityType() {
		t.Fatal("ForgotToken.EntityType returned incorrect EntityType.")
	}
}

func TestPreSave(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var ft model.ForgotToken

	// test with no player
	if err = ft.PreSave(ctx); err == nil {
		t.Fatal("ForgotToken.PreSave did not throw an error for missing Parent Player Key when one should not exist.")
	}

	player := createPlayer(ctx, t)

	// test proper
	ft.PlayerKey = player.GetKey()
	if err = ft.PreSave(ctx); err != nil {
		t.Fatal("ForgotToken.PreSave threw an error for missing Parent Player Key when one should exist.")
	}

	key := ft.GetKey()
	if key == nil {
		t.Fatal("ForgotToken.PreSave did not create a datastore key.")
	}
}

func TestPostLoad(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var ft model.ForgotToken

	testToken := createRandForgotToken(ctx, t)

	// test with no key
	err = ft.PostLoad(ctx)
	if _, ok := err.(*db.MissingKeyError); !ok {
		t.Fatalf("ForgotToken.PostLoad threw the wrong error for missing key. Error: %v", err)
	}

	ft.SetKey(testToken.GetKey())

	// test proper
	err = ft.PostLoad(ctx)
	if err != nil {
		t.Fatalf("ForgotToken.PostLoad threw an error: %v", err)
	}
	if ft.PlayerKey == nil {
		t.Fatal("ForgotToken.PostLoad did not attach the Parent Player Key.")
	}
}

func TestClearExisting(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var ft model.ForgotToken

	token := random.Stringn(64)

	player := createPlayer(ctx, t)

	// create the proper tokens
	numTokens := 4
	for i := 0; i < numTokens; i++ {
		createForgotToken(ctx, t, player, token)
	}

	// create some other tokens
	player2 := createRandPlayer(ctx, t)
	numOtherTokens := 3
	for i := 0; i < numOtherTokens; i++ {
		createForgotToken(ctx, t, player2, token)
	}

	// do a manual selection here because .ByToken fails if more than one found
	ftList := make([]model.ForgotToken, 0)
	keys, err := datastore.NewQuery(ft.EntityType()).
		Filter("Token =", token).
		KeysOnly().
		GetAll(ctx, &ftList)
	if err != nil {
		t.Fatalf("TestClearExisting first query threw an error. Error: %v", err)
	}
	if len(keys) != numTokens+numOtherTokens {
		t.Fatal("TestClearExisting did not create the proper number of ForgotTokens.")
	}

	// test proper
	num, err := ft.ClearExisting(ctx, player.GetKey())
	if err != nil {
		t.Fatalf("ForgotToken.ClearExisting threw an error. Error: %v", err)
	}
	if num != numTokens {
		t.Fatalf("ForgotToken.ClearExisting reported the wrong number of tokens deleted: %d. Wanted: %d", num, numTokens)
	}

	// run a get on all the keys to force an update
	var getForgotToken model.ForgotToken // for use in the forcing Get, not actual data
	for _, v := range keys {
		datastore.Get(ctx, v, &getForgotToken)
	}

	// do a manual selection here because .ByToken fails if more than one found
	ftList = make([]model.ForgotToken, 0)
	keys, err = datastore.NewQuery(ft.EntityType()).
		Filter("Token =", token).
		KeysOnly().
		GetAll(ctx, &ftList)
	if err != nil {
		t.Fatalf("TestClearExisting second query threw an error. Error: %v", err)
	}
	if len(keys) != numOtherTokens {
		t.Fatal("ForgotToken.ClearExisting did not delete the correct number of ForgotTokens.")
	}
}

func TestClearExpired(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var ft model.ForgotToken

	player := createPlayer(ctx, t)

	// create the proper tokens
	numTokens := 4
	for i := 0; i < numTokens; i++ {
		createPastToken(ctx, t, player)
	}

	// create some other tokens
	numOtherTokens := 3
	for i := 0; i < numOtherTokens; i++ {
		createRandForgotTokenForPlayer(ctx, t, player)
	}

	// do a manual selection here because .ByToken fails if more than one found
	ftList := make([]model.ForgotToken, 0)
	keys, err := datastore.NewQuery(ft.EntityType()).
		Ancestor(player.GetKey()).
		KeysOnly().
		GetAll(ctx, &ftList)
	if err != nil {
		t.Fatalf("TestClearExpired first query threw an error. Error: %v", err)
	}
	if len(keys) != numTokens+numOtherTokens {
		t.Fatalf("TestClearExpired did not create the correct number of ForgotTokens: %d. Wanted: %d.", len(keys), numTokens+numOtherTokens)
	}

	// test proper
	num, err := ft.ClearExpired(ctx)
	if err != nil {
		t.Fatalf("ForgotToken.ClearExpired threw an error. Error: %v", err)
	}
	if num != numTokens {
		t.Fatalf("ForgotToken.ClearExpired reported the wrong number of tokens deleted: %d. Wanted: %d", num, numTokens)
	}

	// run a get on all the keys to force an update
	var getForgotToken model.ForgotToken // for use in the forcing Get, not actual data
	for _, v := range keys {
		datastore.Get(ctx, v, &getForgotToken)
	}

	// do a manual selection here because .ByToken fails if more than one found
	ftList = make([]model.ForgotToken, 0)
	keys, err = datastore.NewQuery(ft.EntityType()).
		Ancestor(player.GetKey()).
		KeysOnly().
		GetAll(ctx, &ftList)
	if err != nil {
		t.Fatalf("TestClearExpired second query threw an error. Error: %v", err)
	}
	if len(keys) != numOtherTokens {
		t.Fatalf("ForgotToken.ClearExpired did not delete the correct number of ForgotTokens: %d. Wanted: %d.", numTokens+numOtherTokens-len(keys), numTokens)
	}
}

func TestByEmail(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var ft model.ForgotToken

	// test with no player
	err = ft.ByEmail(ctx, email)
	if e, ok := err.(*db.UnfoundObjectError); !ok || e.EntityType != new(model.Player).EntityType() {
		t.Fatalf("ForgotToken.ByEmail threw the wrong error with no player: %v", err)
	}

	player := createPlayer(ctx, t)

	// test with no tokens
	err = ft.ByEmail(ctx, email)
	if e, ok := err.(*db.UnfoundObjectError); !ok || e.EntityType != new(model.ForgotToken).EntityType() {
		t.Fatalf("ForgotToken.ByEmail threw the wrong error: %v", err)
	}

	createPastToken(ctx, t, player)

	// test with expired token
	err = ft.ByEmail(ctx, email)
	if e, ok := err.(*db.UnfoundObjectError); !ok || e.EntityType != new(model.ForgotToken).EntityType() {
		t.Fatalf("ForgotToken.ByEmail threw the wrong error with expired token: %v", err)
	}

	testToken := createRandForgotTokenForPlayer(ctx, t, player)

	// test proper
	err = ft.ByEmail(ctx, email)
	if err != nil {
		t.Fatalf("ForgotToken.ByEmail threw an error: %v", err)
	}
	if testToken.Token != ft.Token {
		t.Fatal("ForgotToken.ByEmail returned the wrong ForgotToken.")
	}

	createRandForgotTokenForPlayer(ctx, t, player)

	// test with too many tokens
	var ft2 model.ForgotToken
	err = ft2.ByEmail(ctx, email)
	if _, ok := err.(*game.MultipleObjectError); !ok {
		t.Fatalf("ForgotToken.ByEmail returned the wrong error when two matched ForgotTokens existed. Got: %T '%v'", err, err)
	}
}

func TestByToken(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var ft model.ForgotToken

	token := random.Stringn(64)

	// test with no tokens
	err = ft.ByToken(ctx, token)
	if err == nil {
		t.Fatal("ForgotToken.ByToken returned a ForgotToken when none should exist.")
	}
	if _, ok := err.(*db.UnfoundObjectError); !ok {
		t.Fatalf("ForgotToken.ByToken did not throw an unfound object error when none should exist: Type: %T; Error: %v", err, err)
	}

	player := createPlayer(ctx, t)

	createPastToken(ctx, t, player)

	// test with expired token
	err = ft.ByToken(ctx, email)
	if err == nil {
		t.Fatal("ForgotToken.ByToken returned a ForgotToken when none should exist (Expired).")
	}
	if _, ok := err.(*db.UnfoundObjectError); !ok {
		t.Fatalf("ForgotToken.ByToken did not throw an unfound object error with an expired token: Type: %T; Error: %v", err, err)
	}

	createForgotToken(ctx, t, player, token)

	// test proper
	err = ft.ByToken(ctx, token)
	if err != nil {
		t.Fatalf("ForgotToken.ByToken threw an error: %v", err)
	}
	if token != ft.Token {
		t.Fatal("ForgotToken.ByToken returned the wrong ForgotToken.")
	}

	createForgotToken(ctx, t, player, token)

	// test with too many tokens
	var ft2 model.ForgotToken
	err = ft2.ByToken(ctx, token)
	if _, ok := err.(*game.MultipleObjectError); !ok {
		t.Fatalf("ForgotToken.ByToken returned the wrong error when two matched ForgotTokens existed. Got: %T '%v'", err, err)
	}
}

// CONTROLLER TESTS

func TestCreateToken(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var ft model.ForgotToken

	createPlayer(ctx, t)

	// make sure it's empty
	err = ft.ByEmail(ctx, email)
	if e, ok := err.(*db.UnfoundObjectError); !ok || e.EntityType != new(model.ForgotToken).EntityType() {
		t.Fatalf("TestCreateToken threw an unexpected error. Error: %v", err)
	}

	// test proper
	token, err := CreateToken(ctx, email)
	if err != nil {
		t.Fatalf("CreateToken threw an error. Error: %v", err)
	}

	// force availability
	var getForgotToken model.ForgotToken
	datastore.Get(ctx, token.GetKey(), &getForgotToken)

	// make sure it's full
	err = ft.ByEmail(ctx, email)
	if err != nil {
		t.Fatalf("ForgotToken.ByEmail threw an error. Error: %v", err)
	}

	// make sure it's correct
	if ft.Token != token.Token {
		t.Fatal("ForgotToken.ByEmail somehow returned an incorrect token.")
	}
}

func TestTestToken(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error

	token := random.Stringn(63)

	player1 := createPlayer(ctx, t)

	token1 := createForgotToken(ctx, t, player1, token)

	err = db.Delete(ctx, &player1)
	if err != nil {
		t.Fatalf("Unable to delete test Player. Error: %v", err)
	}

	// force availability
	var getPlayer model.Player
	datastore.Get(ctx, player1.GetKey(), &getPlayer)

	// test with deleted player
	lu, err := TestToken(ctx, token)
	if err == nil {
		t.Fatal("TestToken did not throw an error when the Player was deleted.")
	}

	err = db.Delete(ctx, &token1)
	if err != nil {
		t.Fatalf("Unable to delete test ForgotToken. Error: %v", err)
	}

	// force availability
	var getForgotToken model.ForgotToken
	datastore.Get(ctx, token1.GetKey(), &getForgotToken)

	// make sure it's empty
	lu, err = TestToken(ctx, token)
	if err == nil {
		t.Fatal("TestToken did not return an error when no tokens exist.")
	}
	if lu.Email != "" {
		t.Fatal("TestToken returned a Player when none should have been found.")
	}

	player := createPlayer(ctx, t)

	createForgotToken(ctx, t, player, token)

	// test proper
	lu, err = TestToken(ctx, token)
	if err != nil {
		t.Fatalf("TestToken threw an error. Error: %v", err)
	}
	if lu.Email != player.Email {
		t.Fatal("TestToken returned the wrong Player.")
	}

	createForgotToken(ctx, t, player, token)

	// test with too many tokens
	lu, err = TestToken(ctx, token)
	if err == nil {
		t.Fatal("TestToken did not throw an error when multiple tokens exist.")
	}
	if _, ok := err.(*game.MultipleObjectError); !ok {
		t.Fatalf("TestToken threw the wrong error when multiple tokens exist. Got: %T '%v'", err, err)
	}
}

func TestTokenChangePassword(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var ft model.ForgotToken

	player := createPlayer(ctx, t)
	newPass := random.Stringn(10)

	ft = createPastToken(ctx, t, player)

	// test with expired token
	err = TokenChangePassword(ctx, ft.Token, newPass)
	if err == nil {
		t.Fatal("TokenChangePassword did not throw an error with an expired token.")
	}

	err = datastore.Get(ctx, player.GetKey(), &player)
	if err != nil {
		t.Fatalf("TestTokenChangePassword could not get the test Player. Error: %v", err)
	}
	if password.Compare(player.PasswordHash, newPass) {
		t.Fatal("TokenChangePassword changed the password with an expired token.")
	}

	// run a manual query to prevent .ByToken from clearing it out
	ftList := make([]model.ForgotToken, 0)
	keys, err := datastore.NewQuery(ft.EntityType()).
		Filter("Token =", ft.Token).
		GetAll(ctx, &ftList)
	if err != nil {
		t.Fatalf("TestTokenChangePassword first query threw an error. Error: %v", err)
	}
	if len(keys) != 0 {
		t.Fatal("TokenChangePassword did not clear out the expired token.")
	}

	player = createPlayer(ctx, t)
	newPass = random.Stringn(10)

	createPastToken(ctx, t, player)
	ft = createRandForgotTokenForPlayer(ctx, t, player)

	// test with one valid and one expired token
	err = TokenChangePassword(ctx, ft.Token, newPass)
	if err != nil {
		t.Fatalf("TokenChangePassword threw an error. Error: %v", err)
	}

	_, err = db.Load(ctx, player.GetKey(), &player)
	if err != nil {
		t.Fatalf("TestTokenChangePassword was unable to retrieve the test player. Error: %v", err)
	}
	if !password.Compare(player.PasswordHash, newPass) {
		t.Fatal("TokenChangePassword did not change the password.")
	}

	// run a manual query to prevent .ByToken from clearing it out
	ftList = make([]model.ForgotToken, 0)
	keys, err = datastore.NewQuery(ft.EntityType()).
		Filter("Token =", ft.Token).
		GetAll(ctx, &ftList)
	if err != nil {
		t.Fatalf("TestTokenChangePassword second query threw an error. Error: %v", err)
	}
	if len(keys) != 0 {
		t.Fatal("TokenChangePassword did not clear out the expired token nor the valid token.")
	}
}

func TestClearTokens(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var ft model.ForgotToken

	player := createPlayer(ctx, t)

	// create the proper tokens
	numTokens := 4
	for i := 0; i < numTokens; i++ {
		createRandForgotTokenForPlayer(ctx, t, player)
	}

	player2 := createRandPlayer(ctx, t)

	// create some other tokens
	numOtherTokens := 3
	for i := 0; i < numOtherTokens; i++ {
		createRandForgotTokenForPlayer(ctx, t, player2)
	}

	// test proper
	err = ClearTokens(ctx, player.Email)
	if err != nil {
		t.Fatalf("ClearTokens threw an error. Error: %v", err)
	}

	err = ft.ByEmail(ctx, player.Email)
	if e, ok := err.(*db.UnfoundObjectError); !ok || e.EntityType != new(model.ForgotToken).EntityType() {
		t.Fatalf("TestClearTokens ByEmail threw the wrong error. Error: %v", err)
	}
}

func TestSendForgotEmail(t *testing.T) {
	defer test.ResetDB()
	ctx := test.GetCtx()
	var err error
	var ft model.ForgotToken

	// setup the config values
	config.GenRoot()
	config.RootURL = "UNIT_TESTING"

	token := random.Stringn(10)

	player := createPlayer(ctx, t)
	ft = createForgotToken(ctx, t, player, token)

	err = sendForgotEmail(ctx, email, ft)
	if err != nil {
		t.Fatalf("sendForgotEmail threw an error. Error: %v", err)
	}

	// manual checking for the email is required
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

func createFullForgotToken(ctx context.Context, t *testing.T, playerKey *datastore.Key, expires time.Time, token string) model.ForgotToken {
	file, line, funct := test.GetCaller()

	thing := model.ForgotToken{
		PlayerKey: playerKey,
		Expires:   expires,
		Token:     token,
	}

	if err := db.Save(ctx, &thing); err != nil {
		t.Fatalf("Could not save the test ForgotToken. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
	}

	// perform a Get to force the key to be applied so it's available in queries
	var get model.ForgotToken // for use in the forcing Get, not actual data
	err := datastore.Get(ctx, thing.GetKey(), &get)
	if err != nil {
		t.Fatalf("Could not get the test ForgotToken. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
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

func createForgotToken(ctx context.Context, t *testing.T, player model.Player, token string) model.ForgotToken {
	return createFullForgotToken(ctx, t, player.GetKey(), expiry, token)
}

func createPastToken(ctx context.Context, t *testing.T, player model.Player) model.ForgotToken {
	expires := expiry.AddDate(-1, 0, 0) // 1 year in the past
	token := random.Stringn(64)

	return createFullForgotToken(ctx, t, player.GetKey(), expires, token)
}

func createRandForgotTokenForPlayer(ctx context.Context, t *testing.T, player model.Player) model.ForgotToken {
	token := random.Stringnt(64, random.ALPHANUMERIC)

	return createFullForgotToken(ctx, t, player.GetKey(), expiry, token)
}

func createRandForgotToken(ctx context.Context, t *testing.T) model.ForgotToken {
	player := createPlayer(ctx, t)
	token := random.Stringnt(64, random.ALPHANUMERIC)

	return createFullForgotToken(ctx, t, player.GetKey(), expiry, token)
}
