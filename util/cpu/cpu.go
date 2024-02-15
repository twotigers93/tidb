// Copyright 2022 PingCAP, Inc.
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

package cpu

import (
	"runtime"
	"sync"
	"time"

	"github.com/twotigers93/tidb/util/mathutil"
	"go.uber.org/atomic"
)

var cpuUsage atomic.Float64

// If your kernel is lower than linux 4.7, you cannot get the cpu usage in the container.
var unsupported atomic.Bool

// GetCPUUsage returns the cpu usage of the current process.
func GetCPUUsage() (float64, bool) {
	return cpuUsage.Load(), unsupported.Load()
}

// Observer is used to observe the cpu usage of the current process.
type Observer struct {
	utime int64
	stime int64
	now   int64
	exit  chan struct{}
	cpu   mathutil.ExponentialMovingAverage
	wg    sync.WaitGroup
}

// NewCPUObserver returns a cpu observer.
func NewCPUObserver() *Observer {
	return &Observer{
		exit: make(chan struct{}),
		now:  time.Now().UnixNano(),
		cpu:  *mathutil.NewExponentialMovingAverage(0.95, 10),
	}
}

// Start starts the cpu observer.
func (c *Observer) Start() {
	if runtime.GOOS == "darwin" {
		return
	}
}

// Stop stops the cpu observer.
func (c *Observer) Stop() {
	close(c.exit)
	c.wg.Wait()
}
