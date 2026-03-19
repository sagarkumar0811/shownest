package errors

import "context"

type contextKey struct{}

// WithError attaches an AppError to the context.
func WithError(ctx context.Context, err *AppError) context.Context {
	return context.WithValue(ctx, contextKey{}, err)
}

// FromContext retrieves the AppError stored in the context, if any.
func FromContext(ctx context.Context) (*AppError, bool) {
	err, ok := ctx.Value(contextKey{}).(*AppError)
	return err, ok && err != nil
}
