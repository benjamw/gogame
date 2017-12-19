package config

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	hasGameDefaults bool

	/***********************************************************
	  The Root* variables should not be set in a local.go file
	  as they are set dynamically elsewhere in the code
	 ***********************************************************/

	// Root is the root path to the project (e.g.- /var/www/path/to/vdot)
	Root string

	// RootURL is the root URL to the project (e.g.- localhost:8080)
	RootURL string

	/* END Root* VARIABLES */

	/*** Game Settings ***/

	// SiteName is the name of this website
	SiteName string

	// RoutePriority is the default priority value for mux routes
	// lower ( > 0 ) is better
	RoutePriority int

	/*** Email Settings ***/

	// FromEmail the email address to send from
	FromEmail string

	// TestToEmail the email address to send emails to when testing
	// This email address needs to be whitelisted in the MailGun Sandbox subdomain
	TestToEmail string

	/*** MailGun Settings ***/

	// MailGunDomain is the domain to use for MailGun
	MailGunDomain string

	// MailGunAPIKey is the private API key for MailGun
	MailGunAPIKey string

	// MailGunPubKey is the public key for MailGun
	MailGunPubKey string

	/*** Cloud Storage Settings ***/

	// StorageBucket stores the name of the Bucket used in CloudStorage
	StorageBucket string

	// StorageBaseURL is ???
	StorageBaseURL string

	/*** Forgot Password Token Settings ***/

	// FPTokenExpiry is the time in days to expire forgot password tokens
	FPTokenExpiry int

	/*** BCrypt password settings ***/

	// BcryptCost is the cost of the bcrypt hashing function
	BcryptCost int

	/*** Cookie settings ***/

	// CookieCryptKey is the encryption key for the session cookie
	// must be 32 bytes
	CookieCryptKey []byte

	// CookieSignatureKey is the encryption signature for the session cookie
	CookieSignatureKey []byte

	/*** Login Redirect Settings ***/

	// FrontLogin contains the relative path to the user login page
	FrontLogin string

	// AdminLogin contains the relative path to the vendor admin login page
	AdminLogin string

	/*** Encryption settings ***/

	// EncryptionKey is the encryption key for general encryption on the server
	EncryptionKey []byte
)

func init() {
	SetGameConfig()
}

func SetGameConfig() {
	if hasGameDefaults {
		return
	}
	hasGameDefaults = true

	RoutePriority = 9999

	FPTokenExpiry = 1

	BcryptCost = bcrypt.DefaultCost

	FrontLogin = "/login"
	AdminLogin = "/admin/#/login"
}

// GenRoot generates and sets the Root path value for the config data
func GenRoot() {
	if Root == "" {
		cwd, err := os.Getwd()
		if err != nil {
			Root = "."
		}

		cwParts := strings.SplitAfter(cwd, string(os.PathSeparator)+"gogame")

		Root = filepath.Clean(cwParts[0])
	}
}

// SetRoot sets the Root path with the given data
func SetRoot(root string) {
	Root = filepath.Clean(root)
}
