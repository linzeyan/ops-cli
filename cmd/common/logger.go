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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var stdLogger = newLogger()

type logger struct {
	Config zap.Config
	Log    *zap.Logger
}

func (*logger) DefaultEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

func (l *logger) DefaultConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    l.DefaultEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func (l *logger) Sync() error { return l.Log.Sync() }

func newLogger() *logger {
	log := new(logger)
	log.Config = log.DefaultConfig()
	log.Log, _ = log.Config.Build()
	return log
}

func SetLoggerLevel(level string) {
	switch level {
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
	}
}

func NewLogger() *zap.Logger {
	return stdLogger.Log
}
