package logger

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	colorjson "github.com/TylerBrock/colorjson"
	logrusloki "github.com/schoentoon/logrus-loki"
	logrus "github.com/sirupsen/logrus"
	viper "github.com/spf13/viper"
	config "movies/utils/config"
)

var client *logrus.Logger
var clientOnce sync.Once

// deflog: return default logger
func deflog() *logrus.Logger {
	clientOnce.Do(func() {
		client = logrus.New()

		// Set log level if verbose mode is enabled
		if viper.GetBool("verbose") {
			client.SetLevel(logrus.DebugLevel)
		}
		
		if config.Loki().Enabled() {
			hook, err := logrusloki.NewLokiDefaults(config.Loki().URL())
			if err != nil {
				client.Fatal(err)
			}
			client.AddHook(hook)
		}

	})
	return client
}

/*============================================================================*/
/*=====*                            Context                             *=====*/
/*============================================================================*/

type ctxKey string

const loggerCtxKey ctxKey = "logger"

type loggerCtx struct {
	labels map[string]any
	mutex  sync.RWMutex
}

func AddCtxLabel(ctx context.Context, key string, value any) context.Context {
	// Ensure we have a loggerCtx
	ctxValue := ctx.Value(loggerCtxKey)
	data, ok := ctxValue.(*loggerCtx)
	if !ok {
		data = &loggerCtx{
			labels: make(map[string]any),
			mutex:  sync.RWMutex{},
		}
	}

	// Add label
	key = normalizeKey(key)
	data.mutex.Lock()
	data.labels[key] = value
	data.mutex.Unlock()

	// Return new context
	return context.WithValue(ctx, loggerCtxKey, data)
}

/*============================================================================*/
/*=====*                             Export                             *=====*/
/*============================================================================*/

func Debug(format string, args ...any) {
	deflog().Debugf(format, args...)
}

func JSON[T any](body T) {
	f := colorjson.NewFormatter()
	f.Indent = 2

	unmarshal := func(v []byte) (any, error) {
		var a map[string]any
		if err := json.Unmarshal(v, &a); err != nil {
			var b []any
			if err := json.Unmarshal(v, &b); err != nil {
				return nil, err
			}
			return b, nil
		}
		return a, nil
	}

	var err error
	var parsed any
	switch v := any(body).(type) {
	case []byte:
		parsed, err = unmarshal(v)
	case string:
		parsed, err = unmarshal([]byte(v))
	default:
		inrec, err2 := json.Marshal(body)
		if err2 != nil {
			Debug("Failed to marshal JSON: %v", err)
			return
		}
		parsed, err = unmarshal(inrec)
	}
	if err != nil {
		Debug("Failed to marshal JSON: %v", err)
	}

	if s, err := f.Marshal(parsed); err != nil {
		Debug("Failed to marshal JSON: %v", err)
	} else {
		Debug("%s", s)
	}
}

func Info(ctx context.Context, format string, args ...any) {
	logCtx(ctx, deflog(), logrus.InfoLevel, format, args...)
}

func Warn(ctx context.Context, format string, args ...any) {
	logCtx(ctx, deflog(), logrus.WarnLevel, format, args...)
}

func Error(ctx context.Context, format string, args ...any) {
	logCtx(ctx, deflog(), logrus.ErrorLevel, format, args...)
}

/*============================================================================*/
/*=====*                           Labelling                            *=====*/
/*============================================================================*/

type entry struct {
	logger *logrus.Entry
}

func (e entry) Debug(format string, args ...any) {
	e.logger.Debugf(format, args...)
}

func (e entry) Info(ctx context.Context, format string, args ...any) {
	logCtx(ctx, e.logger, logrus.InfoLevel, format, args...)
}

func (e entry) Warn(ctx context.Context, format string, args ...any) {
	logCtx(ctx, e.logger, logrus.WarnLevel, format, args...)
}

func (e entry) Error(ctx context.Context, format string, args ...any) {
	logCtx(ctx, e.logger, logrus.ErrorLevel, format, args...)
}

func (e entry) With(key string, value any) *entry {
	key = normalizeKey(key)
	return &entry{logger: e.logger.WithField(key, value)}
}

func With(key string, value any) interface {
	Debug(string, ...any)
	Info(context.Context, string, ...any)
	Warn(context.Context, string, ...any)
	Error(context.Context, string, ...any)
	With(key string, value any) *entry
} {
	key = normalizeKey(key)
	return &entry{logger: deflog().WithField(key, value)}
}

/*============================================================================*/
/*=====*                              Utils                             *=====*/
/*============================================================================*/

func logCtx(
	ctx context.Context,
	output interface {
		WithField(string, any) *logrus.Entry
		Logf(logrus.Level, string, ...any)
	},
	level logrus.Level,
	format string,
	args ...any,
) {
	// Insert context data
	value := ctx.Value(loggerCtxKey)
	data, ok := value.(*loggerCtx)
	if !ok {
		output.Logf(level, format, args...)
	} else if len(data.labels) == 0 {
		output.Logf(level, format, args...)
	} else {
		// Lock data for reading
		data.mutex.RLock()
		defer data.mutex.RUnlock()

		// Iterate over labels
		var entry *logrus.Entry
		for k, v := range data.labels {
			if entry == nil {
				entry = output.WithField(k, v)
			} else {
				entry = entry.WithField(k, v)
			}
		}
		entry.Logf(level, format, args...)
	}
}

// normalizeKey: normalize key for logging
func normalizeKey(key string) string {
	key = strings.ToLower(key)
	key = strings.TrimSpace(key)
	key = strings.ReplaceAll(key, " ", "_")
	key = strings.ReplaceAll(key, "-", "_")
	return key
}
