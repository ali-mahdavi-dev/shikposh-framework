package commandmiddleware

import "context"

// CommandHandlerFunc is the base function type for command handlers
type CommandHandlerFunc func(ctx context.Context, cmd any) error

// Middleware is a function that wraps a CommandHandlerFunc with additional behavior
type Middleware func(next CommandHandlerFunc) CommandHandlerFunc

// Chain chains multiple middlewares together, applying them in order
// The first middleware in the slice will be the outermost layer
func Chain(middlewares ...Middleware) Middleware {
	return func(next CommandHandlerFunc) CommandHandlerFunc {
		// Apply middlewares in reverse order so the first one is outermost
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Apply applies a middleware to a handler function
func Apply(handler CommandHandlerFunc, middleware Middleware) CommandHandlerFunc {
	return middleware(handler)
}

// ApplyChain applies a chain of middlewares to a handler function
func ApplyChain(handler CommandHandlerFunc, middlewares ...Middleware) CommandHandlerFunc {
	return Chain(middlewares...)(handler)
}

