package commandmiddleware

import (
	"context"
	"reflect"

	"github.com/ali-mahdavi-dev/shikposh-framework/infrastructure/logging"
)

// Logging creates a middleware that logs command execution
func Logging() Middleware {
	return func(next CommandHandlerFunc) CommandHandlerFunc {
		return func(ctx context.Context, cmd any) error {
			cmdName := reflect.TypeOf(cmd).String()

			logging.Info("Command received").
				WithAny("command_name", cmdName).
				WithAny("command", cmd).
				Log()

			err := next(ctx, cmd)
			if err != nil {
				logging.Error("Command failed").
					WithAny("command_name", cmdName).
					WithAny("command", cmd).
					WithError(err).
					Log()
				return err
			}

			logging.Info("Command handled successfully").
				WithAny("command_name", cmdName).
				WithAny("command", cmd).
				Log()

			return nil
		}
	}
}
