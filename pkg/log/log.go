package log

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

func New() *zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Str("role", filepath.Base(os.Args[0])).Logger()
	return &log
}
