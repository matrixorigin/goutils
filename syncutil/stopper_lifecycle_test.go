// Copyright 2014 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.
//
//
//
// Copyright 2017-2019 Lei Ni (nilei81@gmail.com) and other Dragonboat authors.
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

package syncutil

import (
	"runtime"
	"sync/atomic"
	"testing"
)

func TestStopJoinsWorkerGoexit(t *testing.T) {
	s := NewStopper()
	exited := make(chan struct{})
	s.RunWorker(func() { defer close(exited); runtime.Goexit() })
	<-exited
	s.Stop()
}

func TestStoppingWorkerCanRejectFurtherAdmission(t *testing.T) {
	s := NewStopper()
	exited := make(chan struct{})
	var admitted int32
	s.RunWorker(func() {
		defer close(exited)
		<-s.ShouldStop()
		s.RunWorker(func() { atomic.AddInt32(&admitted, 1) })
	})
	s.Stop()
	select {
	case <-exited:
	default:
		t.Fatal("Stop returned before its admitted worker joined")
	}
	if atomic.LoadInt32(&admitted) != 0 {
		t.Fatal("work admitted after stop signal")
	}
}
