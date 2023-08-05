package logging

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var structuredEnabled bool

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	structuredEnabled = strings.ToLower(strings.TrimSpace(os.Getenv("LOG_MODE"))) == "structured"
	color.NoColor = structuredEnabled
}

type ctxKey string

const (
	loggerKey ctxKey = "logger"

	LogSrcKey string = "log_src"
)

func WithLogSource(log *zerolog.Logger, src string) zerolog.Logger {
	if !structuredEnabled {
		src = fmt.Sprintf("[%s]", color.CyanString(src))
	}

	return log.With().Str(LogSrcKey, src).Logger()
}

var (
	defaultLogger     *zerolog.Logger
	defaultLoggerOnce sync.Once
)

func NewLogger(lvl string, pretty bool) *zerolog.Logger {
	var sink *zerolog.ConsoleWriter

	if pretty {
		sink = prettyLogger(os.Stderr)
	} else {
		sink = &zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: true,
		}
	}

	zlvl, err := zerolog.ParseLevel(lvl)
	if err != nil {
		zlvl = zerolog.TraceLevel
	}

	logger := zerolog.New(sink).Level(zlvl).With().Timestamp().Logger()

	return &logger
}

func prettyLogger(out io.Writer) *zerolog.ConsoleWriter {
	return &zerolog.ConsoleWriter{
		Out:     out,
		NoColor: false,
		// Parts - are fields, that are displayed not in structured way, but just sequentually
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			LogSrcKey,
			zerolog.MessageFieldName,
		},
		FieldsExclude: []string{LogSrcKey},
		FormatLevel:   devLevelFormatter,
		// FormatExtra: func(m map[string]any, _ *bytes.Buffer) error {
		// 	for k, v := range m {
		// 		fmt.Printf("-> %s == %v\n", k, v)
		// 	}
		// 	return nil
		// },
	}
}

func NewLoggerFromEnv() *zerolog.Logger {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvl = "trace"
	}
	structuredEnabled := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_MODE"))) == "structured"

	return NewLogger(lvl, !structuredEnabled)
}

func CtxWithLogger(ctx context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *zerolog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*zerolog.Logger); ok {
		return logger
	}

	return DefaultLogger()
}

func DefaultLogger() *zerolog.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLoggerFromEnv()
	})

	return defaultLogger
}

//---------------------------------------------------

func devLevelFormatter(i interface{}) string {
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
