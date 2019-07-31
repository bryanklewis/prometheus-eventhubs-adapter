package util

/*
  Copyright 2019 Micron Technology, Inc.

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
/*
  This work contains copyrighted material, see NOTICE
  for additional information.
  ---------------------------------------------------
  Copyright 2019 Timescale, Inc., Apache License 2.0
*/

import (
	"sync"
	"time"
)

// ThroughputCalc runs on scheduled interval to calculate the throughput per second and sends results to a channel
type ThroughputCalc struct {
	tickInterval time.Duration
	previous     float64
	current      chan float64
	Values       chan float64
	running      bool
	lock         sync.Mutex
}

// NewThroughputCalc creates a new throughput calculation
func NewThroughputCalc(interval time.Duration) *ThroughputCalc {
	return &ThroughputCalc{tickInterval: interval, current: make(chan float64, 1), Values: make(chan float64, 1)}
}

// SetCurrent sends the throughput value to the current channel
func (dt *ThroughputCalc) SetCurrent(value float64) {
	select {
	case dt.current <- value:
	default:
	}
}

// Start creates a thread to run the scheduled throughput interval
func (dt *ThroughputCalc) Start() {
	dt.lock.Lock()
	defer dt.lock.Unlock()
	if !dt.running {
		dt.running = true
		ticker := time.NewTicker(dt.tickInterval)
		go func() {
			for range ticker.C {
				if !dt.running {
					return
				}
				current := <-dt.current
				diff := current - dt.previous
				dt.previous = current
				select {
				case dt.Values <- diff / dt.tickInterval.Seconds():
				default:
				}
			}
		}()
	}
}
