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
func Register(ctx context.Context, username, email, pass string) (p model.Player, e error) {
	hooks.Do("PreRegister", ctx, username, email, pass)

	p = model.Player{}
	e = p.ByEmail(ctx, email)
	if _, ok := e.(*db.UnfoundObjectError); !ok {
		// check the login
		if p, _, err := Login(ctx, email, pass); err == nil {
			return p, nil
		}

		if e == nil {
			e = &game.AccountExistsError{
				Email: email,
			}
		}

		p = model.Player{}

		return
	}

	e = password.Validate(pass)
	if e != nil {
		p = model.Player{}

		return
	}

	p.Username = username
	p.Email = email
	p.PasswordHash = password.Encode(pass)

	e = db.Save(ctx, &p)
	if e != nil {
		p = model.Player{}

		return
	}

	hooks.Do("Register", ctx, p, pass)

	return
}

// Login the user with the given credentials
func Login(ctx context.Context, email, pass string) (p model.Player, s session.Data, e error) {
	p = model.Player{}
	e = p.ByEmail(ctx, email)
	if e != nil {
		e = &game.InvalidCredentialsError{}
	}

	if !password.Compare(p.PasswordHash, pass) {
		e = &game.InvalidCredentialsError{}
		return
	}

	s.IsPlayer = true
	s.PlayerID = p.GetKey().Encode()

	hooks.Do("Login", ctx, &s)

	return
}

// Update an existing user with the given information
func Update(ctx context.Context, oldEmail, oldPass, newEmail, newPass string) (model.Player, error) {
	// test password
	p, _, myerr := Login(ctx, oldEmail, oldPass)
	if myerr != nil {
		return model.Player{}, myerr
	}

	if oldEmail != newEmail {
		myerr = p.ByEmail(ctx, newEmail)
		if _, ok := myerr.(*db.UnfoundObjectError); !ok {
			if myerr == nil {
				myerr = &game.AccountExistsError{
					Email: newEmail,
				}
			}

			return model.Player{}, myerr
		}
	}

	old := p

	p.Email = newEmail

	if newPass != "" {
		p.PasswordHash = password.Encode(newPass)
	}

	myerr = db.Save(ctx, &p)
	if myerr != nil {
		return model.Player{}, myerr
	}

	hooks.Do("Update", ctx, old, p, newPass)

	return p, myerr
}

// GetDeleteToken creates a token for use when deleting an account
func GetDeleteToken(ctx context.Context, plyrID, pass string) (token string, myerr error) {
	var old model.Player
	_, myerr = db.LoadS(ctx, plyrID, &old)
	if myerr != nil {
		return
	}
	pk := old.GetKey()

	// test password
	_, _, myerr = Login(ctx, old.Email, pass)
	if myerr != nil {
		return
	}

	dt := model.DeleteToken{
		Token: model.Token{
			PlayerKey: pk,
			Token:     random.Stringnt(64, random.ALPHANUMERIC),
		},
	}
	dt.ClearExisting(ctx, dt.PlayerKey)
	myerr = db.Save(ctx, &dt)
	if myerr != nil {
		return
	}

	token = dt.Token.Token

	return
}

// Delete the user with the given token
func Delete(ctx context.Context, plyrID, token string) (myerr error) {
	var old model.Player
	_, myerr = db.LoadS(ctx, plyrID, &old)
	if myerr != nil {
		return
	}

	dt := model.DeleteToken{}
	myerr = dt.ByToken(ctx, token)
	if myerr != nil {
		return
	}

	if dt.PlayerKey != old.GetKey() {
		myerr = &game.InvalidCredentialsError{}
	}

	dt.ClearExisting(ctx, dt.PlayerKey)

	p := old
	myerr = db.Delete(ctx, &p)
	if myerr != nil {
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
