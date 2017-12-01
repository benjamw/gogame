package player

import (
	"context"

	"github.com/benjamw/golibs/hooks"

	"github.com/benjamw/gogame/model"
	"github.com/benjamw/gogame/session"
)

func init() {
	hooks.Register("PreRegister", &PreRegistrationListener{})
	hooks.Register("Register", &RegistrationListener{})
	hooks.Register("Login", &LoginListener{})
	hooks.Register("Logout", &LogoutListener{})
	hooks.Register("Update", &UpdateListener{})
	hooks.Register("Delete", &DeleteListener{})

	// now that the hooks are registered, add the listeners (in listeners.go)
	listen()
}

// PreRegistrationListener is a hook that runs before registration
type PreRegistrationListener struct {
	// H is the function that gets processed by the Doer
	// Parameters:
	//	The requested username
	//	The requested email
	//	The requested password in plain text
	H func(context.Context, string, string, string) (bool, error)
}

// Do satisfies the hook.Doer interface
func (h *PreRegistrationListener) Do(ctx context.Context, p ...interface{}) (bool, error) {
	if 3 < len(p) {
		panic("too many parameters passed to registration doer")
	}

	var ok bool

	var username string
	if username, ok = p[1].(string); !ok {
		panic("second parameter of pre-registration doer is of invalid type")
	}

	var email string
	if email, ok = p[1].(string); !ok {
		panic("third parameter of pre-registration doer is of invalid type")
	}

	var pass string
	if pass, ok = p[2].(string); !ok {
		panic("forth parameter of pre-registration doer is of invalid type")
	}

	return h.H(ctx, username, email, pass)
}

// RegistrationListener is a hook that runs on successful registration
type RegistrationListener struct {
	// H is the function that gets processed by the Doer
	// Parameters:
	//	The newly registered Player model
	//	The newly registered player's password in plain text
	H func(context.Context, model.Player, string) (bool, error)
}

// Do satisfies the hook.Doer interface
func (h *RegistrationListener) Do(ctx context.Context, p ...interface{}) (bool, error) {
	if 2 < len(p) {
		panic("too many parameters passed to registration doer")
	}

	var ok bool

	var plyr model.Player
	if plyr, ok = p[0].(model.Player); !ok {
		panic("second parameter of registration doer is of invalid type")
	}

	var pass string
	if pass, ok = p[1].(string); !ok {
		panic("third parameter of registration doer is of invalid type")
	}

	return h.H(ctx, plyr, pass)
}

// LoginListener is a hook that runs on successful login
type LoginListener struct {
	// H is the function that gets processed by the Doer
	// Parameters:
	//	The logged in session data
	H func(ctx context.Context, s session.Cookier) (bool, error)
}

// Do satisfies the hook.Doer interface
func (h *LoginListener) Do(ctx context.Context, p ...interface{}) (bool, error) {
	if 1 < len(p) {
		panic("too many parameters passed to login doer")
	}

	var ok bool

	var s session.Cookier
	if s, ok = p[0].(session.Cookier); !ok {
		// if this panics, make sure the parameter passed is a pointer
		// because the session.Cookier interface methods require a pointer receiver
		panic("second parameter of login doer is of invalid type")
	}

	return h.H(ctx, s)
}

// LogoutListener is a hook that runs on successful logout
type LogoutListener struct {
	// H is the function that gets processed by the Doer
	// Parameters:
	//	The logged in session data
	H func(ctx context.Context, s session.Cookier) (bool, error)
}

// Do satisfies the hook.Doer interface
func (h *LogoutListener) Do(ctx context.Context, p ...interface{}) (bool, error) {
	if 1 < len(p) {
		panic("too many parameters passed to logout doer")
	}

	var ok bool

	var s session.Cookier
	if s, ok = p[0].(session.Cookier); !ok {
		// if this panics, make sure the parameter passed is a pointer
		// because the session.Cookier interface methods require a pointer receiver
		panic("second parameter of logout doer is of invalid type")
	}

	return h.H(ctx, s)
}

// UpdateListener is a hook that runs on successful profile update
type UpdateListener struct {
	// H is the function that gets processed by the Doer
	// Parameters:
	//	The previous Player model data (old data)
	//	The updated Player model data (new data)
	//	The current (new) player's password in plain text
	H func(context.Context, model.Player, model.Player, string) (bool, error)
}

// Do satisfies the hook.Doer interface
func (h *UpdateListener) Do(ctx context.Context, p ...interface{}) (bool, error) {
	if 3 < len(p) {
		panic("too many parameters passed to update doer")
	}

	var ok bool

	var old model.Player
	if old, ok = p[0].(model.Player); !ok {
		panic("second parameter of update doer is of invalid type")
	}

	var plyr model.Player
	if plyr, ok = p[1].(model.Player); !ok {
		panic("third parameter of update doer is of invalid type")
	}

	var pass string
	if pass, ok = p[2].(string); !ok {
		panic("fourth parameter of update doer is of invalid type")
	}

	return h.H(ctx, old, plyr, pass)
}

// DeleteListener is a hook that runs on successful profile deletion
type DeleteListener struct {
	// H is the function that gets processed by the Doer
	// Parameters:
	//	The Player model data (old data)
	H func(context.Context, model.Player) (bool, error)
}

// Do satisfies the hook.Doer interface
func (h *DeleteListener) Do(ctx context.Context, p ...interface{}) (bool, error) {
	if 1 < len(p) {
		panic("too many parameters passed to delete doer")
	}

	var ok bool

	var old model.Player
	if old, ok = p[0].(model.Player); !ok {
		panic("second parameter of delete doer is of invalid type")
	}

	return h.H(ctx, old)
}
