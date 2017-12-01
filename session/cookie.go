package session

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"net/http"
	"time"

	"github.com/benjamw/golibs/crypto"

	"github.com/benjamw/gogame/config"
)

const ( // reset iota
	playerFlag uint8 = 1 << iota // (0x01)
)

const (
	version uint8 = 1
)

// SignAndEncode the given message with the given keys
// Opposite of DecodeAndCheckSig
func SignAndEncode(msg bytes.Buffer, sigKey, cryptKey []byte) (string, error) {
	// Add HMAC
	signed := crypto.AddSignature(msg.Bytes(), sigKey)

	// Encrypt
	crypted, err := crypto.Encrypt(signed, cryptKey)
	if err != nil {
		return "", err
	}

	// Base64-encode
	return base64.URLEncoding.EncodeToString(crypted), nil
}

// DecodeAndCheckSig the given message with the given keys
// Opposite of SignAndEncode
func DecodeAndCheckSig(msg string, sigKey, cryptKey []byte) ([]byte, error) {
	// Base64-decode
	cryptBytes, err := base64.URLEncoding.DecodeString(msg)
	if err != nil {
		return []byte{}, err
	}

	// Decrypt
	rawBytes, err := crypto.Decrypt(cryptBytes, cryptKey)
	if err != nil {
		return []byte{}, err
	}

	// Check HMAC
	signedBytes, err := crypto.CheckSignature(rawBytes, sigKey)
	if err != nil {
		return []byte{}, err
	}

	return signedBytes, nil
}

// ToCookie writes the given data to a cookie with the given name
// The cookie is assumed to be a session cookie, and therefore will expire
// in 24 hours
func ToCookie(data, name string, w http.ResponseWriter) {
	expires := time.Now().Add(24 * time.Hour) // +1 day
	cookie := http.Cookie{Name: name, Value: data, Expires: expires, Path: "/"}
	http.SetCookie(w, &cookie)
}

// KillCookie deletes the given cookie
func KillCookie(name string, w http.ResponseWriter) {
	expires := time.Now().Add(-8760 * time.Hour) // -1 year
	cookie := http.Cookie{Name: name, Value: "", Expires: expires, Path: "/"}
	http.SetCookie(w, &cookie)
}

// Data is the session data for the current user's session
type Data struct {
	IsPlayer    bool
	IsSuperUser bool
	PlayerID    string
}

// Serialize the session data
func (s *Data) Serialize() (string, error) {
	var msg bytes.Buffer

	msg.WriteByte(version)
	var flagByte uint8
	if s.IsPlayer {
		flagByte |= playerFlag
		msg.WriteByte(flagByte)

		size := make([]byte, binary.MaxVarintLen64)
		n := binary.PutUvarint(size, uint64(len(s.PlayerID)))
		msg.Write(size[:n])

		msg.WriteString(s.PlayerID)
	}

	return SignAndEncode(msg, config.CookieSignatureKey, config.CookieCryptKey)
}

// Deserialize an encoded string to the session
func (s *Data) Deserialize(rawString string) error {
	signedBytes, err := DecodeAndCheckSig(rawString, config.CookieSignatureKey, config.CookieCryptKey)
	if err != nil {
		return err
	}

	encoded := bytes.NewBuffer(signedBytes)
	vers, err := encoded.ReadByte()
	if err != nil {
		return err
	}
	if vers != version {
		return errors.New("incorrect byte version")
	}

	flagByte, _ := encoded.ReadByte()
	strCount, _ := binary.ReadUvarint(encoded)

	var letters []byte
	var letter byte
	for ; strCount > 0; strCount-- {
		letter, _ = encoded.ReadByte()
		letters = append(letters, letter)
	}

	if flagByte&playerFlag != 0 {
		s.IsPlayer = true
		s.PlayerID = string(letters)
	}

	return nil
}

// ToCookie store the session as an encoded cookie with the given name
func (s *Data) ToCookie(w http.ResponseWriter, name string) error {
	data, err := s.Serialize()
	if err != nil {
		return err
	}

	ToCookie(data, name, w)

	return nil
}

// KillCookie kills the session cookie with the given name
func (s *Data) KillCookie(w http.ResponseWriter, name string) error {
	s = new(Data)

	KillCookie(name, w)

	return nil
}

// FromCookie pull the session from a cookie with the given name and decode
func (s *Data) FromCookie(r *http.Request, name string) (found bool, err error) {
	var c *http.Cookie
	c, err = r.Cookie(name)
	if err != nil {
		err = nil
		found = false
		return
	}

	err = s.Deserialize(c.Value)
	if err != nil {
		found = false
		return
	}

	found = true
	return
}
