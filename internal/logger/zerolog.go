package logger

import (
	"os"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

// NewZeroLogr returns a Zerolog logger wrapped in a go-logr/logr interface.
func NewZeroLogr() logr.Logger {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	zl := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return zerologr.New(&zl)
}
