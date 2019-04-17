// gomuks - A terminal Matrix client written in Go.
// Copyright (C) 2019 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/kennetanti/gomuks/config"
)

func TestNewConfig_Defaults(t *testing.T) {
	cfg := config.NewConfig("/tmp/gomuks-test-0", "/tmp/gomuks-test-0")
	assert.Equal(t, "/tmp/gomuks-test-0", cfg.Dir)
	assert.Equal(t, "/tmp/gomuks-test-0/history.db", cfg.HistoryPath)
	assert.Equal(t, "/tmp/gomuks-test-0/media", cfg.MediaDir)
}

func TestConfig_Load_NonexistentDoesntFail(t *testing.T) {
	cfg := config.NewConfig("/tmp/gomuks-test-1", "/tmp/gomuks-test-1")

	defer os.RemoveAll("/tmp/gomuks-test-1")

	cfg.Load()

	stat, err := os.Stat(cfg.MediaDir)
	assert.Nil(t, err)
	assert.True(t, stat.IsDir())

	/* FIXME
	stat, err = os.Stat(cfg.HistoryDir)
	assert.Nil(t, err)
	assert.True(t, stat.IsDir())*/
}

func TestConfig_Load_DirectoryFails(t *testing.T) {
	os.MkdirAll("/tmp/gomuks-test-2/config.yaml", 0700)
	cfg := config.NewConfig("/tmp/gomuks-test-2", "/tmp/gomuks-test-2")

	defer os.RemoveAll("/tmp/gomuks-test-2")
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Load() didn't panic")
		}
	}()

	cfg.Load()
}

func TestConfig_Load_ExistingFileIsLoaded(t *testing.T) {
	os.MkdirAll("/tmp/gomuks-test-3", 0700)
	ioutil.WriteFile("/tmp/gomuks-test-3/config.yaml", []byte(`{
		"mxid": "foo",
		"homeserver": "bar",
		"history_path": "/tmp/gomuks-test-3/foo.db",
		"media_dir": "/tmp/gomuks-test-3/bar"
	}`), 0700)
	cfg := config.NewConfig("/tmp/gomuks-test-3", "/tmp/gomuks-test-3")

	defer os.RemoveAll("/tmp/gomuks-test-3")

	cfg.Load()

	assert.Equal(t, "foo", cfg.UserID)
	assert.Equal(t, "bar", cfg.HS)
	assert.Equal(t, "/tmp/gomuks-test-3/foo.db", cfg.HistoryPath)
	assert.Equal(t, "/tmp/gomuks-test-3/bar", cfg.MediaDir)

	stat, err := os.Stat(cfg.MediaDir)
	assert.Nil(t, err)
	assert.True(t, stat.IsDir())

	/* FIXME
	stat, err = os.Stat(cfg.HistoryDir)
	assert.Nil(t, err)
	assert.True(t, stat.IsDir())*/
}

func TestConfig_Load_InvalidExistingFilePanics(t *testing.T) {
	os.MkdirAll("/tmp/gomuks-test-4", 0700)
	ioutil.WriteFile("/tmp/gomuks-test-4/config.yaml", []byte(`this is not JSON.`), 0700)
	cfg := config.NewConfig("/tmp/gomuks-test-4", "/tmp/gomuks-test-4")

	defer os.RemoveAll("/tmp/gomuks-test-4")
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Load() didn't panic")
		}
	}()

	cfg.Load()
}

func TestConfig_Clear(t *testing.T) {
	cfg := config.NewConfig("/tmp/gomuks-test-5", "/tmp/gomuks-test-5")

	defer os.RemoveAll("/tmp/gomuks-test-5")

	cfg.Load()

	stat, err := os.Stat(cfg.MediaDir)
	assert.Nil(t, err)
	assert.True(t, stat.IsDir())

	/* FIXME
	stat, err = os.Stat(cfg.HistoryDir)
	assert.Nil(t, err)
	assert.True(t, stat.IsDir())*/

	cfg.Clear()

	stat, err = os.Stat(cfg.MediaDir)
	assert.True(t, os.IsNotExist(err))
	assert.Nil(t, stat)

	/* FIXME
	stat, err = os.Stat(cfg.HistoryDir)
	assert.True(t, os.IsNotExist(err))
	assert.Nil(t, stat)*/
}

func TestConfig_Save(t *testing.T) {
	cfg := config.NewConfig("/tmp/gomuks-test-6", "/tmp/gomuks-test-6")

	defer os.RemoveAll("/tmp/gomuks-test-6")

	cfg.Load()
	cfg.Save()

	dat, err := ioutil.ReadFile("/tmp/gomuks-test-6/config.yaml")
	assert.Nil(t, err)
	assert.Contains(t, string(dat), "/tmp/gomuks-test-6")
}
