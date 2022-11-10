/*
Copyright © 2022 ZeYanLin <zeyanlin@outlook.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var stdLogger = newLogger()

type logger struct {
	Config zap.Config
	Log    *zap.Logger
}

func (l *logger) DefaultConfig() zap.Config {
	return zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: true,
		Sampling: &zap.SamplingConfig{
			Initial:    10, /* Log number in same level output per second. */
			Thereafter: 10, /* If greater than 10 then output. */
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "call",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.NanosDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
			// ConsoleSeparator: "\t",
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		// InitialFields:    map[string]any{},
	}
}

func (l *logger) Sync() error { return l.Log.Sync() }

/* newLogger returns logger with configured zap.Config. */
func newLogger() *logger {
	log := new(logger)
	log.Config = log.DefaultConfig()
	log.Log, _ = log.Config.Build()
	return log
}

/* SetLoggerLevel set log level from given level. */
func SetLoggerLevel(level string) {
	switch strings.ToLower(level) {
	case zap.DebugLevel.String():
		stdLogger.Config.Level.SetLevel(zap.DebugLevel)
	case zap.InfoLevel.String():
		stdLogger.Config.Level.SetLevel(zap.InfoLevel)
	case zap.WarnLevel.String():
		stdLogger.Config.Level.SetLevel(zap.WarnLevel)
	case zap.ErrorLevel.String():
		stdLogger.Config.Level.SetLevel(zap.ErrorLevel)
	case zap.DPanicLevel.String():
		stdLogger.Config.Level.SetLevel(zap.DPanicLevel)
	case zap.PanicLevel.String():
		stdLogger.Config.Level.SetLevel(zap.PanicLevel)
	case zap.FatalLevel.String():
		stdLogger.Config.Level.SetLevel(zap.FatalLevel)
	default:
		stdLogger.Log.Fatal(ErrInvalidArg.Error(), DefaultField(level))
	}
}

/* NewLogger return field Log of global variable stdLogger. */
func NewLogger() *zap.Logger {
	return stdLogger.Log
}

/* DefaultField calls zap.Any and gives key name "arg". */
func DefaultField(value any) zapcore.Field {
	return zap.Any("arg", value)
}

/* NewField calls zap.Any. */
func NewField(key string, value any) zapcore.Field {
	return zap.Any(key, value)
}
