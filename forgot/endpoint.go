package forgot

import (
	"context"
	"net/http"
	"strings"

	gttp "github.com/benjamw/gogame/http"
	"github.com/benjamw/golibs/db"
)

func init() {
	gttp.R.Path("/forgot").
		Methods("POST").
		Handler(&gttp.JSONHandler{handleForgotPassword})

	gttp.R.Path("/change_password/{token:[a-zA-Z0-9]+}").
		Methods("POST").
		Handler(&gttp.JSONHandler{handleChangePassword})
}

func handleForgotPassword(ctx context.Context, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	email := r.FormValue("email")
	if email == "" {
		errReply = &gttp.MissingRequiredError{FormElement: "email"}
		return
	}

	// this will also delete any existing tokens
	token, errReply := CreateToken(ctx, email)
	if errReply != nil {
		if _, ok := errReply.(*db.UnfoundObjectError); ok {
			// don't give hackers any info on whether or not this email address exists
			errReply = nil
			replyRaw = gttp.Response{
				Success: true,
			}
		}

		return
	}

	if errReply = SendForgotEmailDelay.Call(ctx, email, token.Token); errReply != nil {
		return
	}

	replyRaw = gttp.Response{
		Success: true,
	}

	return
}

func handleChangePassword(ctx context.Context, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	token := gttp.GetURLValue(r, "token")
	token = token + strings.Repeat("-", 64)
	token = token[:64]

	password := r.FormValue("password")
	if password == "" {
		errReply = &gttp.MissingRequiredError{FormElement: "password"}
		return
	}

	if errReply = TokenChangePassword(ctx, token, password); errReply != nil {
		return
	}

	reply := gttp.Response{
		Success: true,
	}

	replyRaw = reply

	return
}
