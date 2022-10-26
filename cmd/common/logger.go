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

func (l *logger) Debug()      { l.Config.Level.SetLevel(zap.DebugLevel) }
func (l *logger) Info()       { l.Config.Level.SetLevel(zap.InfoLevel) }
func (l *logger) Warn()       { l.Config.Level.SetLevel(zap.WarnLevel) }
func (l *logger) Error()      { l.Config.Level.SetLevel(zap.ErrorLevel) }
func (l *logger) DPanic()     { l.Config.Level.SetLevel(zap.DPanicLevel) }
func (l *logger) Panic()      { l.Config.Level.SetLevel(zap.PanicLevel) }
func (l *logger) Fatal()      { l.Config.Level.SetLevel(zap.FatalLevel) }
func (l *logger) Sync() error { return l.Log.Sync() }

func newLogger() *logger {
	l := new(logger)
	l.Config = l.DefaultConfig()
	l.Log, _ = l.Config.Build()
	return l
}

func NewLogger() *logger {
	return stdLogger
}
