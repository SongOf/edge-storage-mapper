/*
Copyright 2020 The KubeEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"k8s.io/klog/v2"
	"time"
)

// Timer is to call a function periodically.
type Timer struct {
	Function func()
	Duration time.Duration
	Times    int
	Shutdown chan string
}

// Start start a timer.
func (t *Timer) Start() {
	ticker := time.NewTicker(t.Duration)
	if t.Times > 0 {
		for i := 0; i < t.Times; i++ {
			select {
			case <-ticker.C:
				t.Function()
			case <-t.Shutdown:
				klog.Info("timer is exit")
				return
			}
		}
	} else {
		for {
			select {
			case <-ticker.C:
				t.Function()
			case <-t.Shutdown:
				klog.Info("timer is exit")
				return
			}
		}
	}
}

func (t *Timer) Terminated() {
	t.Shutdown <- "Done"
}
