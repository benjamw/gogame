package player

import (
	"context"
	"net/http"
	"time"

	"github.com/benjamw/golibs/db"
	"google.golang.org/appengine/datastore"

	"strings"

	"github.com/benjamw/gogame/game"
	gttp "github.com/benjamw/gogame/http"
	"github.com/benjamw/gogame/model"
	"github.com/benjamw/gogame/session"
)

func init() {
	gttp.R.Path("/register").
		Methods("POST").
		Handler(&gttp.JSONHandler{handleRegister})

	gttp.R.Path("/login").
		Methods("POST").
		Handler(&gttp.JSONHandler{handleLogin})

	gttp.R.Path("/logout").
		Methods("GET").
		Handler(&gttp.PlayerJSONHandler{handleLogout})

	gttp.R.Path("/update").
		Methods("PUT").
		Handler(&gttp.PlayerJSONHandler{handleUpdate})

	gttp.R.Path("/delete/{token:[a-zA-Z0-9]+}").
		Methods("DELETE").
		Handler(&gttp.PlayerJSONHandler{handleDelete})

	gttp.R.Path("/delete").
		Methods("POST").
		Handler(&gttp.PlayerJSONHandler{handlePreDelete})

	gttp.R.Path("/ping").
		Handler(&gttp.PlayerJSONHandler{handlePing})
}

type Reply struct {
	gttp.Response
	ID         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	CookieName string `json:"cookie_name,omitempty"`
	CookieData string `json:"cookie_data,omitempty"`
}

func (r *Reply) Set(p model.Player) {
	r.ID = p.GetKey().Encode()
	r.Username = p.Username
	r.Email = p.Email
}

func (r *Reply) SetCookie(c string) {
	r.CookieName = gttp.PlayerCookieName
	r.CookieData = c
}

func handleRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	username := r.FormValue("username")
	if username == "" {
		errReply = &gttp.MissingRequiredError{FormElement: "username"}
		return
	}
	email := r.FormValue("email")
	if email == "" {
		errReply = &gttp.MissingRequiredError{FormElement: "email"}
		return
	}
	pass := r.FormValue("password")
	if pass == "" {
		errReply = &gttp.MissingRequiredError{FormElement: "password"}
		return
	}

	plyr, errReply := Register(ctx, username, email, pass)
	if errReply != nil {
		return
	}

	reply := Reply{}
	reply.Success = true
	reply.Set(plyr)

	replyRaw = reply

	return
}

func handleLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	email := r.FormValue("email")
	if email == "" {
		errReply = &gttp.MissingRequiredError{FormElement: "email"}
		return
	}

	pass := r.FormValue("password")
	if pass == "" {
		errReply = &gttp.MissingRequiredError{FormElement: "password"}
		return
	}

	plr, sess, errReply := Login(ctx, email, pass)
	if errReply != nil {
		return
	}

	c, errReply := sess.Serialize()
	if errReply != nil {
		return
	}

	setCookie(w, sess)

	reply := Reply{}
	reply.Success = true
	reply.Set(plr)
	reply.SetCookie(c)

	replyRaw = reply

	return
}

func handleLogout(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	killCookie(w, s)

	reply := Reply{}
	reply.Success = true

	replyRaw = reply

	return
}

func handleUpdate(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	plyrKey, errReply := datastore.DecodeKey(s.PlayerID)
	if errReply != nil {
		return
	}

	var old model.Player
	_, errReply = db.Load(ctx, plyrKey, &old)
	if errReply != nil {
		return
	}

	oldPass := r.FormValue("password")
	if oldPass == "" {
		errReply = &gttp.MissingRequiredError{FormElement: "password"}
		return
	}

	email := r.FormValue("email")
	pass := r.FormValue("new_password")

	plyr, errReply := Update(ctx, old.Email, oldPass, email, pass)
	if errReply != nil {
		return
	}

	reply := Reply{}
	reply.Success = true
	reply.Set(plyr)

	replyRaw = reply

	return
}

type tokenReply struct {
	Token string `json:"token"`
}

func handlePreDelete(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	pass := r.FormValue("password")
	if pass == "" {
		errReply = &gttp.MissingRequiredError{FormElement: "password"}
		return
	}

	var token string
	token, errReply = GetDeleteToken(ctx, s.PlayerID, pass)
	if errReply != nil {
		return
	}

	reply := tokenReply{
		Token: token,
	}

	replyRaw = reply

	return
}

func handleDelete(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	token := gttp.GetURLValue(r, "token")
	token = token + strings.Repeat("-", 64)
	token = token[:64]

	errReply = Delete(ctx, s.PlayerID, token)
	if errReply != nil {
		return
	}

	killCookie(w, s)

	reply := Reply{}
	reply.Success = true

	replyRaw = reply

	return
}

type pingReply struct {
	gttp.Response
	PlayerID string    `json:"player_id"`
	Now      time.Time `json:"timestamp"`
}

func handlePing(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	reply := pingReply{
		Response: gttp.Response{
			Success: true,
		},
		PlayerID: s.PlayerID,
		Now:      game.Now(ctx),
	}

	replyRaw = reply

	return
}
