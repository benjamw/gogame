package session

import (
	"net/http"
)

type Cookier interface {
	Serialize() (string, error)
	Deserialize(string) error
	ToCookie(http.ResponseWriter, string) error
	FromCookie(*http.Request, string) (bool, error)
}
