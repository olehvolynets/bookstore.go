package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	outLogger zerolog.Logger
	errLogger zerolog.Logger
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if os.Getenv("ENV") == "production" {
		outLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		errLogger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		levelFormatter := devLevelFormatter(false)

		outLogger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, FormatLevel: levelFormatter}).
			With().Timestamp().Logger()
		errLogger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, FormatLevel: levelFormatter}).
			With().Timestamp().Logger()
	}
}

func Log() *zerolog.Event   { return outLogger.Log() }
func Debug() *zerolog.Event { return outLogger.Debug() }
func Trace() *zerolog.Event { return outLogger.Trace() }
func Info() *zerolog.Event  { return outLogger.Info() }
func Warn() *zerolog.Event  { return outLogger.Warn() }
func Error() *zerolog.Event { return errLogger.Error() }
func Fatal() *zerolog.Event { return errLogger.Fatal() }
func Panic() *zerolog.Event { return errLogger.Panic() }

func Print(msg string)                            { errLogger.Print(msg) }
func Printf(format string, values ...interface{}) { errLogger.Printf(format, values...) }

//---------------------------------------------------

func devLevelFormatter(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		var lvl string
		var col *color.Color
		if ll, ok := i.(string); ok {
			switch ll {
			case zerolog.LevelTraceValue:
				lvl = "TRC"
				col = color.New(color.FgGreen)
			case zerolog.LevelDebugValue:
				lvl = "DBG"
				col = color.New(color.FgMagenta)
			case zerolog.LevelInfoValue:
				lvl = "INF"
				col = color.New(color.FgBlue)
			case zerolog.LevelWarnValue:
				lvl = "WRN"
				col = color.New(color.FgYellow)
			case zerolog.LevelErrorValue:
				lvl = "ERR"
				col = color.New(color.FgRed, color.Underline)
			case zerolog.LevelFatalValue:
				lvl = "FTL"
				col = color.New(color.FgRed, color.Bold)
			case zerolog.LevelPanicValue:
				lvl = "PNC"
				col = color.New(color.FgRed, color.Bold)
			default:
				lvl = ll
				col = color.New(color.FgWhite, color.Bold)
			}
		} else {
			if i == nil {
				lvl = "???"
				col = color.New(color.FgWhite, color.Bold)
			} else {
				lvl = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
				col = color.New(color.FgWhite, color.Bold)
			}
		}

		return fmt.Sprintf("[%s]", col.Sprint(lvl))
	}
}
