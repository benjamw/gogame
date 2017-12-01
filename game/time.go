package game

import (
	"context"
	"time"
)

var timeNowContextKey = "holds a local version of time.Now"

// Now returns the Now time that is stored in context, or the real Now if not found or not valid
func Now(ctx context.Context) time.Time {
	if now, ok := ctx.Value(&timeNowContextKey).(time.Time); ok {
		return now
	}

	return time.Now()
}

// SetNow sets the given time to be "Now" in the rest of the script
// If no value is passed as "Now", it defaults to the actual Now time
// Only one now value is expected and used
func SetNow(ctx context.Context, now ...time.Time) context.Context {
	// the now slice is here to allow for an empty now value
	// that defaults to the actual now time
	if len(now) == 0 {
		now = make([]time.Time, 1)
		now[0] = time.Now()
	}

	return context.WithValue(ctx, &timeNowContextKey, now[0])
}
