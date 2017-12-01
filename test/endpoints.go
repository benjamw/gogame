package test

import (
	"context"
	"net/http"

	"github.com/davecgh/go-spew/spew"

	gttp "github.com/benjamw/gogame/http"
)

func init() {
	gttp.R.PathPrefix("/").Extensions("html").Handler(&gttp.BlankHandler{handleHTML})
}

func handleHTML(ctx context.Context, w http.ResponseWriter, r *http.Request) (errReply error) {
	spew.Fdump(w, "--- BASE GAME HTML FILE REGEX HANDLER ---")
	spew.Fdump(w, r.URL.Path)

	return
}
