package app

import (
	"net/http"

	gttp "github.com/benjamw/gogame/http"

	// instead of adding all the endpoints here
	// add them to the initializer so it can be imported
	// into various child modules without causing a possible
	// cyclic dependency
	_ "github.com/benjamw/gogame/initializer"
)

func init() {
	http.Handle("/", gttp.R)
}
