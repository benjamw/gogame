package player

import (
	"context"
	"net/http"

	"github.com/benjamw/golibs/db"
	"github.com/benjamw/golibs/hooks"
	"github.com/benjamw/golibs/password"

	"github.com/benjamw/gogame/game"
	gttp "github.com/benjamw/gogame/http"
	"github.com/benjamw/gogame/model"
	"github.com/benjamw/gogame/session"
	"github.com/benjamw/golibs/random"
)

// Register a new user with the given credentials
func Register(ctx context.Context, username, email, pass string) (p model.Player, myerr error) {
	hooks.Do("PreRegister", ctx, username, email, pass)

	p = model.Player{}
	myerr = p.ByEmail(ctx, email)
	if _, ok := myerr.(*db.UnfoundObjectError); !ok {
		// check the login
		if p, _, myerr = Login(ctx, email, pass); myerr == nil {
			return
		}

		if myerr == nil {
			myerr = &game.AccountExistsError{
				Email: email,
			}
		}

		p = model.Player{}

		return
	}

	if myerr = password.Validate(pass); myerr != nil {
		p = model.Player{}

		return
	}

	p.Username = username
	p.Email = email
	p.PasswordHash = password.Encode(pass)

	if myerr = db.Save(ctx, &p); myerr != nil {
		p = model.Player{}

		return
	}

	hooks.Do("Register", ctx, p, pass)

	return
}

// Login the user with the given credentials
func Login(ctx context.Context, email, pass string) (p model.Player, s session.Data, myerr error) {
	p = model.Player{}
	if myerr = p.ByEmail(ctx, email); myerr != nil {
		myerr = &game.InvalidCredentialsError{}
	}

	if !password.Compare(p.PasswordHash, pass) {
		myerr = &game.InvalidCredentialsError{}

		return
	}

	s.IsPlayer = true
	s.PlayerID = p.GetKey().Encode()

	hooks.Do("Login", ctx, &s)

	return
}

// Update an existing user with the given information
func Update(ctx context.Context, oldEmail, oldPass, newEmail, newPass string) (p model.Player, myerr error) {
	// test password
	p, _, myerr = Login(ctx, oldEmail, oldPass)
	if myerr != nil {
		p = model.Player{}

		return
	}

	if oldEmail != newEmail {
		myerr = p.ByEmail(ctx, newEmail)
		if _, ok := myerr.(*db.UnfoundObjectError); !ok {
			if myerr == nil {
				myerr = &game.AccountExistsError{
					Email: newEmail,
				}
			}

			p = model.Player{}

			return
		}
	}

	old := p

	p.Email = newEmail

	if newPass != "" {
		p.PasswordHash = password.Encode(newPass)
	}

	if myerr = db.Save(ctx, &p); myerr != nil {
		p = model.Player{}

		return
	}

	hooks.Do("Update", ctx, old, p, newPass)

	return
}

// GetDeleteToken creates a token for use when deleting an account
func GetDeleteToken(ctx context.Context, plyrID, pass string) (token string, myerr error) {
	var old model.Player
	if _, myerr = db.LoadS(ctx, plyrID, &old); myerr != nil {
		return
	}
	pk := old.GetKey()

	// test password
	if _, _, myerr = Login(ctx, old.Email, pass); myerr != nil {
		return
	}

	dt := model.DeleteToken{
		PlayerKey: pk,
		Value:     random.Stringnt(64, random.ALPHANUMERIC),
	}
	dt.ClearExisting(ctx, dt.PlayerKey)
	if myerr = db.Save(ctx, &dt); myerr != nil {
		return
	}

	token = dt.Value

	return
}

// Delete the user with the given token
func Delete(ctx context.Context, plyrID, token string) (myerr error) {
	var old model.Player
	if _, myerr = db.LoadS(ctx, plyrID, &old); myerr != nil {
		return
	}

	dt := model.DeleteToken{}
	if myerr = dt.ByValue(ctx, token); myerr != nil {
		return
	}

	if dt.PlayerKey != old.GetKey() {
		myerr = &game.InvalidCredentialsError{}
	}

	dt.ClearExisting(ctx, dt.PlayerKey)

	p := old
	if myerr = db.Delete(ctx, &p); myerr != nil {
		return
	}

	hooks.Do("Delete", ctx, old)

	return
}

func setCookie(w http.ResponseWriter, s session.Data) {
	s.ToCookie(w, gttp.PlayerCookieName)
}

func killCookie(w http.ResponseWriter, s session.Data) {
	s.KillCookie(w, gttp.PlayerCookieName)
}
