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

// InstanceConfig 单实例配置
type InstanceConfig struct {
	Logs      map[LogType]*LogConfig `json:"logs"`
	CallDepth int                    `json:"call_depth"`
	Coexist   bool                   `json:"coexist"`
}

// LogConfig of log unit
type LogConfig struct {
	Level zapcore.Level `json:"level"`
	Store *Store        `json:"store,omitempty"`
}

// Store config for File log
type Store struct {
	Path       string `json:"path"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
	LocalZone  bool   `json:"local_zone"`
}

// Config 总配置，实现 tao.MultiConfig 接口
type Config struct {
	tao.BaseMultiConfig[InstanceConfig]
	RunAfters []string `json:"run_after,omitempty" yaml:"run_after,omitempty"`
}

var defaultCommand = &LogConfig{
	Level: zapcore.DebugLevel,
}

var defaultFile = &LogConfig{
	Level: zapcore.DebugLevel,
	Store: &Store{
		Path:       "./test.log",
		MaxSize:    1024,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
		LocalZone:  true,
	},
}

var defaultInstance = &InstanceConfig{
	Logs: map[LogType]*LogConfig{
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
	for i := range z.Instances {
		instance := &z.Instances[i].Cfg
		if instance.Logs == nil {
			instance.Logs = defaultInstance.Logs
		}
		for k, v := range instance.Logs {
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
				delete(instance.Logs, k)
			}
		}
		if instance.CallDepth <= 0 {
			instance.CallDepth = defaultInstance.CallDepth
		}
	}
	if z.RunAfters == nil {
		z.RunAfters = []string{}
	}
}

// ToTask transform itself to Task
func (z *Config) ToTask() tao.Task {
	return tao.NewTask(
		ConfigKey,
		func(ctx context.Context, param tao.Parameter) (tao.Parameter, error) {
			select {
			case <-ctx.Done():
				return param, tao.NewError(tao.ContextCanceled, "%s: context has been canceled", ConfigKey)
			default:
			}
			return param, nil
		})
}

// RunAfter defines pre task names
func (z *Config) RunAfter() []string {
	return z.RunAfters
}
