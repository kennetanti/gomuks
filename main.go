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

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/kennetanti/gomuks/debug"
	"github.com/kennetanti/gomuks/interface"
	"github.com/kennetanti/gomuks/ui"
)

var MainUIProvider ifc.UIProvider = ui.NewGomuksUI

func main() {
	defer debug.Recover()

	enableDebug := len(os.Getenv("DEBUG")) > 0
	debug.RecoverPrettyPanic = !enableDebug

	configDir, err := UserConfigDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get config directory:", err)
		os.Exit(3)
	}
	cacheDir, err := UserCacheDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get cache directory:", err)
		os.Exit(3)
	}

	gmx := NewGomuks(MainUIProvider, configDir, cacheDir)
	gmx.Start()

	// We use os.Exit() everywhere, so exiting by returning from Start() shouldn't happen.
	time.Sleep(5 * time.Second)
	fmt.Println("Unexpected exit by return from gmx.Start().")
	os.Exit(2)
}

func UserCacheDir() (dir string, err error) {
	dir = os.Getenv("GOMUKS_CACHE_HOME")
	if dir == "" {
		dir, err = os.UserCacheDir()
		dir = filepath.Join(dir, "gomuks")
	}
	return
}

func UserConfigDir() (dir string, err error) {
	dir = os.Getenv("GOMUKS_CONFIG_HOME")
	if dir != "" {
		return
	}
	if runtime.GOOS == "windows" {
		dir = os.Getenv("AppData")
		if dir == "" {
			err = errors.New("%AppData% is not defined")
		}
	} else {
		dir = os.Getenv("XDG_CONFIG_HOME")
		if dir == "" {
			dir = os.Getenv("HOME")
			if dir == "" {
				err = errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
			} else if runtime.GOOS == "darwin" {
				dir = filepath.Join(dir, "Library", "Application Support")
			} else {
				dir = filepath.Join(dir, ".config")
			}
		}
	}
	dir = filepath.Join(dir, "gomuks")
	return
}
