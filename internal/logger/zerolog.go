package logger

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

// NewZeroLogr returns a Zerolog logger wrapped in a go-logr/logr interface.
func NewZeroLogr() logr.Logger {
	zl := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	return zerologr.New(&zl)
}
