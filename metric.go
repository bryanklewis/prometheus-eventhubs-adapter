package main

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
  Copyright 2017 The Prometheus Authors, Apache License 2.0
*/

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bryanklewis/prometheus-eventhubs-adapter/util"
)

var (
	adapterInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eventhubs_adapter_info",
			Help: "Remote storage adapter for Azure Event Hubs metadata.",
		},
		[]string{"version", "commit", "build"},
	)
	httpRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all http requests",
		},
	)
	receivedSamples = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "samples_received_total",
			Help: "Total number of received samples.",
		},
	)
	sentSamples = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "samples_sent_total",
			Help: "Total number of processed samples sent to remote storage.",
		},
		[]string{"remote"},
	)
	failedSamples = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "samples_failed_total",
			Help: "Total number of processed samples which failed on send to remote storage.",
		},
		[]string{"remote"},
	)
	sentBatchDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "batch_send_duration_seconds",
			Help:    "Duration of sample batch send calls to the remote storage.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"remote"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_ms",
			Help:    "Duration of HTTP request in milliseconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)

	writeThroughput = util.NewThroughputCalc(time.Second)
)

func init() {
	prometheus.MustRegister(adapterInfo)
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(receivedSamples)
	prometheus.MustRegister(sentSamples)
	prometheus.MustRegister(failedSamples)
	prometheus.MustRegister(sentBatchDuration)
	prometheus.MustRegister(httpRequestDuration)
	writeThroughput.Start()
}
