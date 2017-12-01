package game

import (
	"context"

	netcontext "golang.org/x/net/context"
)

func ConvertContext(ctx context.Context) (c netcontext.Context) {
	c = ctx

	return c
}

func ConvertOldContext(ctx netcontext.Context) (c context.Context) {
	c = ctx

	return c
}
