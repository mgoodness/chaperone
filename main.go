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
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/hugo/watcher"
	"github.com/spf13/pflag"
)

var timeout int
var path, url string
var verbose, version bool

func init() {
	pflag.StringVarP(&path, "path", "p", "", "File or directory to watch for changes")
	pflag.IntVarP(&timeout, "timeout", "t", 30, "Minimum time (in seconds) between POSTs")
	pflag.StringVarP(&url, "url", "u", "", "POST to URL when file changes")
	pflag.BoolVar(&verbose, "verbose", false, "Verbose output")
	pflag.BoolVar(&version, "version", false, "Print version & exit")
	pflag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if path == "" {
		log.Fatal("Path not set.")
	}

	if url == "" {
		log.Fatal("URL not set.")
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if version {
		printVersionAndExit()
	}

	watcher, err := watcher.New(time.Duration(timeout) * time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err := watcher.Add(path); err != nil {
		log.Fatal(err)
	}
	log.Infof("Watching %s...", path)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for {
			select {
			case events := <-watcher.Events:
				for _, event := range events {
					log.Info(event)
				}
				post()
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()
	wg.Wait()
}
