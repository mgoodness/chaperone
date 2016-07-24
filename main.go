// Copyright © 2016 Michael Goodness <mgoodness@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/mgoodness/chaperone/log"
	"github.com/mgoodness/chaperone/process"
	v "github.com/mgoodness/chaperone/version"
	"github.com/spf13/pflag"
)

var dir, exe string
var verbose, version bool

func init() {
	pflag.StringVarP(&dir, "dir", "d", dir, "Directory to monitor for changes")
	pflag.StringVarP(&exe, "exe", "e", exe, "Process name to send SIGHUP")
	pflag.BoolVar(&verbose, "verbose", false, "Verbose output")
	pflag.BoolVar(&version, "version", false, "Print version & exit")
	pflag.Parse()

	if verbose {
		log.SetVerbose()
	}

	if dir == "" {
		log.Fatal("Directory not set.")
	}

	if exe == "" {
		log.Fatal("Process name not set.")
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if version {
		v.PrintVersionAndExit()
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() error {
		for {
			select {
			case event := <-watcher.Events:
				log.Debug(event)
				if err := process.SendSIGHUP(exe); err != nil {
					log.Fatal(err)
				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Watching %s...", dir)
	<-done
}
