package contextkeys

import "context"

// ContextKeys gets and sets context values from a context.
type contextKeys struct {
	xAccountID struct{}
}

var (
	keys = contextKeys{
		xAccountID: struct{}{},
	}
)

// SetXAccountID inserts the account ID value into context.
func SetXAccountID(ctx context.Context, value string) context.Context {
	return keys.setXAccountID(ctx, value)
}

// GetXAccountID returns the account ID from the context.
func GetXAccountID(ctx context.Context) (string, bool) {
	return keys.getXAccountIDFromContext(ctx)
}

func (c contextKeys) getXAccountIDFromContext(ctx context.Context) (string, bool) {
	if nil != ctx {
		xAccountID, ok := ctx.Value(keys.xAccountID).(string)
		return xAccountID, ok
	}

	return "", false
}

func (c contextKeys) setXAccountID(ctx context.Context, value string) context.Context {
	if nil == ctx {
		ctx = context.Background()
	}

	return context.WithValue(ctx, c.xAccountID, value)
}
