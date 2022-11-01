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
	"context"
	"github.com/taouniverse/tao"
	"go.uber.org/zap/zapcore"
)

// ConfigKey for this repo
const ConfigKey = "zap"

// LogType of zap log
type LogType string

const (
	// Console log
	Console LogType = "console"
	// File log
	File LogType = "file"
)

// Config implements tao.Config
type Config struct {
	Logs      map[LogType]*config `json:"logs"`
	CallDepth int                 `json:"call_depth"`
	Coexist   bool                `json:"coexist"`
	RunAfters []string            `json:"run_after,omitempty"`
}

// config of log unit
type config struct {
	Level zapcore.Level `json:"level"`
	Store *store        `json:"store,omitempty"`
}

// store config for File log
type store struct {
	// default value is invalid
	Path       string `json:"path"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`

	// default value is valid
	Compress  bool `json:"compress"`
	LocalZone bool `json:"local_zone"`
}

var defaultCommand = &config{
	Level: zapcore.DebugLevel,
}

var defaultFile = &config{
	Level: zapcore.DebugLevel,
	Store: &store{
		Path:       "./test.log",
		MaxSize:    1024, // 1024 mb
		MaxBackups: 7,    // seven files
		MaxAge:     30,   // one month
		Compress:   true, // compress rotated log files with gzip
		LocalZone:  true, // backup using local time zone
	},
}

var defaultZap = &Config{
	Logs: map[LogType]*config{
		Console: defaultCommand,
	},
	CallDepth: 1,
	Coexist:   false,
}

// Name of Config
func (z *Config) Name() string {
	return ConfigKey
}

// ValidSelf with some default values
func (z *Config) ValidSelf() {
	if z.Logs == nil {
		z.Logs = defaultZap.Logs
	}
	for k, v := range z.Logs {
		switch k {
		case Console:
			if v.Level < zapcore.DebugLevel || v.Level > zapcore.FatalLevel {
				v.Level = defaultCommand.Level
			}
		case File:
			if v.Level < zapcore.DebugLevel || v.Level > zapcore.FatalLevel {
				v.Level = defaultFile.Level
			}
			if v.Store == nil {
				v.Store = defaultFile.Store
			} else {
				if v.Store.Path == "" {
					v.Store.Path = defaultFile.Store.Path
				}
				if v.Store.MaxSize == 0 {
					v.Store.MaxSize = defaultFile.Store.MaxSize
				}
				if v.Store.MaxBackups == 0 {
					v.Store.MaxBackups = defaultFile.Store.MaxBackups
				}
				if v.Store.MaxAge == 0 {
					v.Store.MaxAge = defaultFile.Store.MaxAge
				}
			}
		default:
			delete(z.Logs, k)
		}
	}
	if z.CallDepth <= 0 {
		z.CallDepth = defaultZap.CallDepth
	}
}

// ToTask transform itself to Task
func (z *Config) ToTask() tao.Task {
	return tao.NewTask(
		ConfigKey,
		func(ctx context.Context, param tao.Parameter) (tao.Parameter, error) {
			// non-block check
			select {
			case <-ctx.Done():
				return param, tao.NewError(tao.ContextCanceled, "%s: context has been canceled", ConfigKey)
			default:
			}
			// nothing to do
			return param, nil
		})
}

// RunAfter defines pre task names
func (z *Config) RunAfter() []string {
	return z.RunAfters
}
