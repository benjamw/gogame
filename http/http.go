package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/benjamw/golibs/random"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/benjamw/gogame/config"
	"github.com/benjamw/gogame/game"
)

func buildContext(r *http.Request) context.Context {
	ctx := appengine.NewContext(r)
	ctx = game.SetNow(ctx)
	return ctx
}

func updateConfig(r *http.Request) {
	if config.RootURL == "" {
		config.RootURL = r.Host
	}

	if config.Root == "" {
		config.Root = "."
	}
}

func defaultError(err error, errID int32) (httpCode int, errMessage string) {
	// Any error types we don't specifically look out for
	// default to serving an HTTP 500
	if appengine.IsDevAppServer() {
		errMessage = err.Error()
		httpCode = 555 // 555 to differentiate from a legit 500 when testing
	} else {
		errMessage = fmt.Sprintf("%s (%d)", http.StatusText(http.StatusInternalServerError), errID)
		httpCode = http.StatusInternalServerError
	}

	return
}

func processError(ctx context.Context, err error) (httpCode int, errMessage string) {
	errID := random.Int31()
	errMessage = fmt.Sprintf("%v (%d)", err, errID)

	log.Errorf(ctx, "error (%d): %v", errID, err)

	if e, ok := err.(game.Error); ok {
		httpCode = e.Code()
	} else {
		httpCode, errMessage = defaultError(err, errID)
	}

	return
}

// ParseURL strips the prefix off the url and performs some minor validity tests
func ParseURL(url string, base string) (string, error) {
	var remainder string
	if "/" == base[len(base)-1:] && !strings.Contains(url, base) {
		// assume the URL came in without a trailing slash, but one was expected
		remainder = ""
	} else {
		remainder = strings.Split(url, base)[1]
	}

	// remove any query params
	remainder = strings.Split(remainder, "?")[0]

	if len([]rune(remainder)) == 0 {
		err := errors.New("URL path is missing")
		return "", err
	}

	return remainder, nil
}

// GetURLValue checks the form for the value and if not found there
// will grab any mux var with the same name as the value
// Used for API URLs like /object/read/12345 which returns the object with ID 12345
// where those URLs can also be like /object/read?id=12345
func GetURLValue(r *http.Request, key string) string {
	value := r.FormValue(key)
	if value == "" {
		vars := Vars(r)
		if val, ok := vars[key]; ok {
			value = val
		}
	}

	return value
}

// FormMultiValue is used to allow a client to submit multiple values as
// either a comma-separated string, or an array of inputs, or any combination of both
func FormMultiValue(r *http.Request, key string, sep string) []string {
	r.FormValue("") // force form parsing

	if sep == "" {
		sep = ","
	}

	ret := make([]string, 0)
	value := r.Form[key]
	for _, v := range value {
		vs := strings.Split(v, sep)
		for k := range vs {
			vs[k] = strings.TrimSpace(vs[k])
		}
		ret = append(ret, vs...)
	}

	return ret
}
