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

package process

import (
	"errors"
	"syscall"

	log "github.com/mgoodness/chaperone/log"
	ps "github.com/mitchellh/go-ps"
)

func findByName(processName string) (int, error) {
	pid := 0
	err := errors.New("Process '" + processName + "' not found")
	ps, _ := ps.Processes()
	for i := range ps {
		if ps[i].Executable() == processName {
			pid = ps[i].Pid()
			err = nil
			log.Debugf("Found process '%s' with PID %d", processName, pid)
			break
		}
	}
	return pid, err
}

// SendSIGHUP sends a signal to the process
func SendSIGHUP(processName string) error {
	pid, err := findByName(processName)
	if err != nil {
		return err
	}
	syscall.Kill(pid, syscall.SIGHUP)
	log.Debugf("SIGHUP sent to process %d", pid)
	return nil
}
