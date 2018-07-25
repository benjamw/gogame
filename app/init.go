package app

import (
	"net/http"

	"github.com/benjamw/gogame/config"
	gttp "github.com/benjamw/gogame/http"

	// instead of adding all the endpoints here
	// add them to the initializer so `app` can be imported
	// into various child modules without causing possible
	// cyclic dependencies
	_ "github.com/benjamw/gogame/initializer"
)

func init() {
	config.SetRoot(".")

	http.Handle("/", gttp.R)
}
