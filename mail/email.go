package mail

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/aymerick/raymond"
	"github.com/mailgun/mailgun-go"
	"google.golang.org/appengine/urlfetch"

	"github.com/benjamw/gogame/config"
	"github.com/benjamw/gogame/game"
)

var (
	templates map[string]map[string]*raymond.Template
	tplLock   sync.RWMutex
	tplFiles  = []string{"subject.hbs", "text.hbs", "html.hbs"}
)

// FromTemplate reads the template from settings and sends the email
func FromTemplate(ctx context.Context, template string, to []string, params map[string]interface{}, mg mailgun.Mailgun) error {
	tpls, err := parseTemplates(template, params)
	if err != nil {
		return err
	}

	return Send(ctx, config.FromEmail, to, tpls["subject.hbs"], tpls["text.hbs"], tpls["html.hbs"], mg)
}

// Send an email
func Send(ctx context.Context, from string, to []string, subj string, plain string, html string, mg mailgun.Mailgun) (myerr error) {
	if mg == nil {
		mg = initMailgun(ctx)
	}

	message := mg.NewMessage(from, subj, plain, to...)

	if strings.TrimSpace(html) != "" {
		message.SetHtml(html)
	}

	msg, id, err := mg.Send(message)
	if err != nil {
		myerr = fmt.Errorf("could not send message to %v, ID %v: %v, %v", to, id, msg, err)
		return
	}

	return
}

func parseTemplates(template string, params map[string]interface{}) (tpls map[string]string, myerr error) {
	tplPath := config.Root + "/emails/" + template + "/"

	if params == nil {
		params = make(map[string]interface{}, 0)
	}

	// add some defaults to the param list
	params["ROOT"] = config.RootURL
	params["SiteName"] = config.SiteName

	t := make(map[string]string, 0)
	var result string
	if ts, ok := getTemplates(template); ok {
		for k, v := range ts {
			result, myerr = v.Exec(params)
			if myerr != nil {
				return
			}

			t[k] = result
		}

		tpls = t

		return
	}

	files, myerr := ioutil.ReadDir(tplPath)
	if myerr != nil {
		return
	}

	var tpl *raymond.Template
	f := make(map[string]string, 0)
	for _, file := range files {
		tpl, myerr = game.Templatize(tplPath + file.Name())
		if myerr != nil {
			return
		}
		putTemplate(template, strings.ToLower(file.Name()), tpl)

		result, myerr = tpl.Exec(params)
		if myerr != nil {
			return
		}

		f[strings.ToLower(file.Name())] = strings.TrimSpace(string(result))
	}

	for _, file := range tplFiles {
		if _, ok := f[file]; !ok {
			tpl, myerr = raymond.Parse("")
			if myerr != nil {
				return
			}

			putTemplate(template, file, tpl)

			f[file] = ""
		}

	}

	tpls = f

	return
}

func getTemplates(template string) (tpls map[string]*raymond.Template, ok bool) {
	tplLock.RLock()
	defer tplLock.RUnlock()

	tpls, ok = templates[template]

	return
}

func getTemplate(template string, file string) (tpl *raymond.Template, ok bool) {
	tplLock.RLock()
	defer tplLock.RUnlock()

	tpl, ok = templates[template][file]

	return
}

func putTemplate(template string, file string, tpl *raymond.Template) {
	tplLock.Lock()
	defer tplLock.Unlock()

	if templates == nil {
		templates = make(map[string]map[string]*raymond.Template, 0)
	}

	if templates[template] == nil {
		templates[template] = make(map[string]*raymond.Template, 0)
	}

	templates[template][file] = tpl
}

func initMailgun(ctx context.Context) mailgun.Mailgun {
	httpc := urlfetch.Client(ctx)

	mg := mailgun.NewMailgun(
		config.MailGunDomain,
		config.MailGunAPIKey,
		config.MailGunPubKey,
	)
	mg.SetClient(httpc)

	return mg
}
