// Copyright 2022 huija
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
	"github.com/taouniverse/tao"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

/**
import _ "github.com/taouniverse/tao-zap"
*/

// Z config of zap
var Z = new(Config)

func init() {
	err := tao.Register(ConfigKey, Z, setup)
	if err != nil {
		panic(err.Error())
	}
}

// Logger based on zap
var Logger *zap.SugaredLogger

type zapSetup struct {
	enc   zapcore.Encoder
	level zapcore.Level
	out   zapcore.WriteSyncer
}

var setups []*zapSetup       // core's setup
var ws []zapcore.WriteSyncer // core's root writers
var cores []zapcore.Core     // cores

// setup with zap config
func setup() (err error) {
	for t, c := range Z.Logs {
		switch t {
		case File:
			setupFileLog(c)
		case Console:
			setupConsoleLog(c)
		default:
			return tao.NewErrorWrapped("config: log type invalid", nil)
		}
	}

	// setup cores
	for _, s := range setups {
		cores = append(cores, zapcore.NewCore(s.enc, s.out, s.level))
	}

	// final logger
	Logger = zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(Z.CallDepth),
	).Sugar()

	// add zap to tao
	err = tao.SetLogger(ConfigKey, Logger)
	if err != nil {
		return
	}
	return tao.SetWriter(ConfigKey, zapcore.NewMultiWriteSyncer(ws...))
}

func setupConsoleLog(c *config) {
	// writerSync
	ws = append(ws, zapcore.AddSync(os.Stdout))
	// encoderConfig
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05.000")
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// setup
	setups = append(setups, &zapSetup{
		enc:   zapcore.NewConsoleEncoder(encoderConfig),
		level: c.Level,
		out:   ws[len(ws)-1],
	})
}

func setupFileLog(c *config) {
	sliceLogger := &lumberjack.Logger{
		Filename:   c.Store.Path,
		MaxSize:    c.Store.MaxSize,
		MaxBackups: c.Store.MaxBackups,
		MaxAge:     c.Store.MaxAge,
		Compress:   c.Store.Compress,
		LocalTime:  c.Store.LocalZone,
	}

	// writerSync
	ws = append(ws, zapcore.AddSync(sliceLogger))
	// encoderConfig
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// setup
	setups = append(setups, &zapSetup{
		enc:   zapcore.NewJSONEncoder(encoderConfig),
		level: c.Level,
		out:   ws[len(ws)-1],
	})
}
