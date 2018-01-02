package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/appengine/log"

	"github.com/benjamw/gogame/config"
	"github.com/benjamw/gogame/session"
)

const (
	// PlayerCookieName is the name of the cookie that players will get
	PlayerCookieName = "bue"
)

// BlankHandler handles endpoints with no built-in response
type BlankHandler struct {
	H func(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

func (h BlankHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := PrepHandler(r)

	err := h.H(ctx, w, r)
	if err != nil {
		ReplyData(ctx, w, nil, err)
		return
	}
}

// HTMLHandler handles endpoints with an HTML response
type HTMLHandler struct {
	H func(ctx context.Context, w http.ResponseWriter, r *http.Request) ([]byte, error)
}

func (h HTMLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := PrepHandler(r)

	html, err := h.H(ctx, w, r)

	ReplyHTML(ctx, w, html, err)
	return
}

// JSONHandler handles endpoints with a JSON response
type JSONHandler struct {
	H func(ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error)
}

func (h JSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := PrepHandler(r)
	ctx = context.WithValue(ctx, "w", w)

	data, err := h.H(ctx, w, r)

	ReplyData(ctx, w, data, err)
	return
}

// PlayerBlankHandler requires a Player login and handles endpoints with no built-in response
type PlayerBlankHandler struct {
	H func(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) error
}

func (h PlayerBlankHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, s, err := prepSession(r, PlayerCookieName)
	if _, ok := err.(*session.InvalidSessionError); ok {
		http.Redirect(w, r, config.FrontLogin, http.StatusFound)
		return
	}
	if err != nil {
		ReplyData(ctx, w, nil, err)
		return
	}

	err = h.H(ctx, s, w, r)
	if err != nil {
		ReplyData(ctx, w, nil, err)
		return
	}
}

// PlayerHTMLHandler requires a Player login and handles endpoints with an HTML response
type PlayerHTMLHandler struct {
	H func(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) ([]byte, error)
}

func (h PlayerHTMLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, s, err := prepSession(r, PlayerCookieName)
	if _, ok := err.(*session.InvalidSessionError); ok {
		http.Redirect(w, r, config.FrontLogin, http.StatusFound)
		return
	}
	if err != nil {
		ReplyData(ctx, w, nil, err)
		return
	}

	html, err := h.H(ctx, s, w, r)

	ReplyHTML(ctx, w, html, err)
	return
}

// PlayerJSONHandler requires a Player login and handles endpoints with a JSON response
type PlayerJSONHandler struct {
	H func(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) (interface{}, error)
}

func (h PlayerJSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, s, err := prepSession(r, PlayerCookieName)
	ctx = context.WithValue(ctx, "w", w)
	// no redirect for JSON data
	if err != nil {
		ReplyData(ctx, w, nil, err)
		return
	}

	reply, err := h.H(ctx, s, w, r)

	ReplyData(ctx, w, reply, err)
	return
}

// Response is a generic response to be returned with errors or empty successes
type Response struct {
	Success      bool   `json:"success,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

// ReplyJSON formats and submits a JSON response
func ReplyJSON(ctx context.Context, w http.ResponseWriter, data interface{}) {
	replyBytes, err := json.Marshal(data)
	if err != nil {
		log.Errorf(ctx, "JSON marshalling failed: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(replyBytes)
	return
}

// ReplyErr formats and submits an error response
func ReplyErr(ctx context.Context, w http.ResponseWriter, err error) {
	if e, ok := err.(*RedirectError); ok {
		w.Header().Set("Location", e.URL)
	}

	errCode, errMessage := processError(ctx, err)
	http.Error(w, errMessage, errCode)
	return
}

// ReplyHTML formats and submits an HTML response
func ReplyHTML(ctx context.Context, w http.ResponseWriter, replyBytes []byte, err error) {
	if err != nil {
		ReplyErr(ctx, w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(replyBytes)
	return
}

// ReplyData determines the appropriate response for the given
// data and formats and submits that response
func ReplyData(ctx context.Context, w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		ReplyErr(ctx, w, err)
		return
	}

	if d, ok := data.([]byte); ok {
		w.Write(d)
	} else if d, ok := data.(string); ok {
		w.Write([]byte(d))
	} else {
		ReplyJSON(ctx, w, data)
	}
	return
}

// PrepHandler preps the context as well as the various dynamic config items
func PrepHandler(r *http.Request) context.Context {
	ctx := buildContext(r)
	updateConfig(r)
	return ctx
}

// prepSession should be split so that each module has it's own prepSession that gets invoked
// by the various module handlers
func prepSession(r *http.Request, cookieName string) (ctx context.Context, s session.Data, err error) {
	ctx = PrepHandler(r)

	if cookieName != PlayerCookieName {
		err = fmt.Errorf("unknown cookie encountered: '%s'", cookieName)
		return
	}

	var found bool
	found, err = s.FromCookie(r, PlayerCookieName)
	if !found || !s.IsPlayer {
		err = &session.InvalidSessionError{}
	}

	return
}
