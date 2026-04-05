// Copyright 2021-2026 huija
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zap

import (
	"os"

	"github.com/taouniverse/tao"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/**
import _ "github.com/taouniverse/tao-zap"
*/

// Z is the global config instance for tao-zap
var Z = &Config{}

// Factory is the global factory instance for managing zap.SugaredLogger
var Factory *tao.BaseFactory[*zap.SugaredLogger]

func init() {
	var err error
	Factory, err = tao.Register(ConfigKey, Z, NewZap)
	if err != nil {
		panic(err.Error())
	}
}

// NewZap creates a new Zap logger for factory pattern
func NewZap(name string, config InstanceConfig) (*zap.SugaredLogger, func() error, error) {
	var localSetups []zapSetup
	var localWS []zapcore.WriteSyncer
	var localCores []zapcore.Core

	for t, c := range config.Logs {
		switch t {
		case File:
			ws, setup := buildFileLog(c)
			localWS = append(localWS, ws)
			localSetups = append(localSetups, setup)
		case Console:
			ws, setup := buildConsoleLog(c)
			localWS = append(localWS, ws)
			localSetups = append(localSetups, setup)
		default:
			return nil, nil, tao.NewErrorWrapped("config: log type invalid", nil)
		}
	}

	for _, s := range localSetups {
		localCores = append(localCores, zapcore.NewCore(s.enc, s.out, s.level))
	}

	logger := zap.New(
		zapcore.NewTee(localCores...),
		zap.AddCaller(),
		zap.AddCallerSkip(config.CallDepth),
	).Sugar()

	err := tao.SetLogger(ConfigKey, logger)
	if err != nil {
		return nil, nil, err
	}
	err = tao.SetWriter(ConfigKey, zapcore.NewMultiWriteSyncer(localWS...))
	if !config.Coexist {
		_ = tao.DeleteLogger(tao.ConfigKey)
		_ = tao.DeleteWriter(tao.ConfigKey)
	}

	closer := func() error {
		return logger.Sync()
	}

	return logger, closer, err
}

type zapSetup struct {
	enc   zapcore.Encoder
	level zapcore.Level
	out   zapcore.WriteSyncer
}

func buildConsoleLog(c *LogConfig) (zapcore.WriteSyncer, zapSetup) {
	ws := zapcore.AddSync(os.Stdout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05.000")
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return ws, zapSetup{
		enc:   zapcore.NewConsoleEncoder(encoderConfig),
		level: c.Level,
		out:   ws,
	}
}

func buildFileLog(c *LogConfig) (zapcore.WriteSyncer, zapSetup) {
	sliceLogger := &lumberjack.Logger{
		Filename:   c.Store.Path,
		MaxSize:    c.Store.MaxSize,
		MaxBackups: c.Store.MaxBackups,
		MaxAge:     c.Store.MaxAge,
		Compress:   c.Store.Compress,
		LocalTime:  c.Store.LocalZone,
	}

	ws := zapcore.AddSync(sliceLogger)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return ws, zapSetup{
		enc:   zapcore.NewJSONEncoder(encoderConfig),
		level: c.Level,
		out:   ws,
	}
}
