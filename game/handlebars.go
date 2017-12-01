package game

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/aymerick/raymond"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/benjamw/gogame/config"
)

var (
	helpersRegistered = false
)

func init() {
	if config.Root == "" {
		if appengine.IsDevAppServer() {
			config.GenRoot()
			config.SetRoot(config.Root + "/app")
		} else {
			config.Root = "."
		}
	}
	config.GenRoot()

	registerHelpers()
}

// Templatize reads a file and converts it to a handlebars template.
func Templatize(filename string) (tpl *raymond.Template, myerr error) {
	contents, myerr := ioutil.ReadFile(filename)
	if myerr != nil {
		return
	}

	t, myerr := raymond.Parse(string(contents))
	if myerr != nil {
		return
	}

	tpl = t

	return
}

func registerHelpers() {
	if helpersRegistered {
		return
	}

	raymond.RegisterHelper("Key", func(key *datastore.Key) string {
		return key.Encode()
	})

	raymond.RegisterHelper("ParentKey", func(key *datastore.Key) string {
		return key.Parent().Encode()
	})

	raymond.RegisterHelper("DateLocal", func(date string) string {
		t, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", date)
		return strconv.FormatInt(t.Unix(), 10)
	})

	raymond.RegisterHelper("Spew", func(i interface{}) string {
		return fmt.Sprintf("<pre>%s</pre>", spew.Sdump(i))
	})

	helpersRegistered = true
}
