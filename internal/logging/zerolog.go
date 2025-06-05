package logging

import (
	"fmt"
	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type filteredLevelWriter struct {
	target zerolog.LevelWriter
	level  zerolog.Level
}

var _ zerolog.LevelWriter = &filteredLevelWriter{}

func (l filteredLevelWriter) Write(p []byte) (n int, err error) {
	return l.target.Write(p)
}

func (l filteredLevelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level >= l.level {
		return l.target.WriteLevel(level, p)
	}
	return len(p), nil
}

// newFilteredLevelWriter creates a LevelWriter that filters log entries below the specified level before writing them.
func newFilteredLevelWriter(target zerolog.LevelWriter, level zerolog.Level) zerolog.LevelWriter {
	return &filteredLevelWriter{
		target: target,
		level:  level,
	}
}

type consoleWriterOption func(w *zerolog.ConsoleWriter)

func NewZeroLogLogger(settings *configuration.ApplicationConfiguration) *zerolog.Logger {

	levelString := settings.LogLevel
	level, err := zerolog.ParseLevel(levelString)
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}

	stdOutLogger := zerolog.NewConsoleWriter(
		withTimeFormat(time.RFC3339),
		withLevelFormat(fiveLetterLevelFormat),
		withErrorValueFormat(escapeLineBreaks))

	writer := zerolog.MultiLevelWriter(
		newFilteredLevelWriter(zerolog.MultiLevelWriter(stdOutLogger), level))

	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger
}

func withTimeFormat(timeFormat string) consoleWriterOption {
	return func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = timeFormat
	}
}

func withLevelFormat(formatter func(level zerolog.Level) string) consoleWriterOption {
	return func(w *zerolog.ConsoleWriter) {
		w.FormatLevel = func(i interface{}) string {
			level, ok := i.(zerolog.Level)
			if ok {
				return formatter(level)
			}
			return fmt.Sprintf("%v", i)
		}
	}
}

func withErrorValueFormat(format func(value string) string) consoleWriterOption {
	return func(w *zerolog.ConsoleWriter) {
		w.FormatErrFieldValue = func(i interface{}) string {
			return format(fmt.Sprintf("%v", i))
		}
	}
}

func fiveLetterLevelFormat(level zerolog.Level) string {
	return fmt.Sprintf("%-5s", level)
}

func escapeLineBreaks(s string) string {
	return strings.Replace(s, "\n", "\\n", -1)
}
