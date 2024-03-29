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
	"github.com/stretchr/testify/assert"
	"github.com/taouniverse/tao"
	"testing"
)

func TestTao(t *testing.T) {
	err := tao.SetConfigPath("./test.json")
	assert.Nil(t, err)

	tao.Debug("logger debug")
	tao.Debugf("logger %s", "debugf")
	tao.Info("logger info")
	tao.Infof("logger %s", "infof")
	tao.Warn("logger warn")
	tao.Warnf("logger %s", "warnf")
	tao.Error("logger error")
	tao.Errorf("logger %s", "errorf")
	//tao.Panic("logger panic")
	//tao.Panicf("logger %s", "panicf")
	//tao.Fatal("logger fatal")
	//tao.Fatalf("logger %s", "fatalf")

	err = tao.Run(nil, nil)
	assert.Nil(t, err)
}
