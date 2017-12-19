package forgot

import (
	"context"
	"time"

	"github.com/benjamw/golibs/db"
	"github.com/benjamw/golibs/password"
	"github.com/benjamw/golibs/random"
	netcontext "golang.org/x/net/context"
	"google.golang.org/appengine/delay"

	"github.com/benjamw/gogame/config"
	"github.com/benjamw/gogame/game"
	"github.com/benjamw/gogame/mail"
	"github.com/benjamw/gogame/model"
)

// CreateToken creates a forgot password token for the given email address
func CreateToken(ctx context.Context, email string) (token model.ForgotToken, myerr error) {
	if myerr = ClearTokens(ctx, email); myerr != nil {
		return
	}

	pl := model.Player{}
	if myerr = pl.ByEmail(ctx, email); myerr != nil {
		return
	}

	ft := model.ForgotToken{
		PlayerKey: pl.GetKey(),
		Value:     random.Stringnt(64, random.ALPHANUMERIC),
		Expires:   game.Now(ctx).Add(time.Hour * time.Duration(24*config.FPTokenExpiry)),
	}
	if myerr = db.Save(ctx, &ft); myerr != nil {
		return
	}

	token = ft

	return
}

// TestToken tests the given token and returns the player associated with it if found
func TestToken(ctx context.Context, token string) (pl model.Player, myerr error) {
	var ft model.ForgotToken
	if myerr = ft.ByValue(ctx, token); myerr != nil {
		return
	}

	player := model.Player{}
	if _, err := db.Load(ctx, ft.PlayerKey, &player); err != nil {
		myerr = &db.UnfoundObjectError{
			EntityType: "Entity",
			Key:        "token",
			Value:      token,
			Err:        err,
		}
		return
	}

	pl = player

	return
}

// TokenChangePassword changes the password for the lookup associated to the given token and clears out the token
func TokenChangePassword(ctx context.Context, token string, pass string) (myerr error) {
	var pl model.Player
	if pl, myerr = TestToken(ctx, token); myerr != nil {
		return
	}

	pl.PasswordHash = password.Encode(pass)
	if myerr = db.Save(ctx, &pl); myerr != nil {
		return
	}

	new(model.ForgotToken).ClearExisting(ctx, pl.GetKey())

	return
}

// ClearTokens clears the tokens for the given email address
func ClearTokens(ctx context.Context, email string) (myerr error) {
	var pl model.Player

	if err := pl.ByEmail(ctx, email); err != nil {
		myerr = &db.UnfoundObjectError{
			EntityType: pl.EntityType(),
			Key:        "email",
			Value:      email,
		}

		return
	}

	var ft model.ForgotToken

	ft.ClearExisting(ctx, pl.GetKey())

	return
}

var SendForgotEmailDelay = delay.Func("forgot_email", sendForgotEmail)

// SendForgotEmail
func sendForgotEmail(ctx netcontext.Context, email, token string) error {
	ctx = game.ConvertOldContext(ctx)

	to := make([]string, 0)
	to = append(to, email)

	params := make(map[string]interface{}, 1)
	params["token"] = token

	return mail.FromTemplate(ctx, "forgot", to, params)
}
