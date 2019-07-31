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

import (
	"time"

	"github.com/bryanklewis/prometheus-eventhubs-adapter/log"
)

// Defaults
const (
	// The timeout to use when sending samples to the remote storage.
	defaultTimeout = 5 * time.Second
	// Address to listen on for web endpoints.
	defaultListen = ":9201"
	// Path for write requests.
	defaultWritePath = "/write"
	// Path for telemetry scraps.
	defaultTelemetry = "/metrics"
	// Encoding to use when sending events [ \"json\", \"avro-json\" ].
	defaultEncoding = "json"
)

// Config for general application options
type Config struct {
	RemoteTimeout time.Duration
	ListenAddr    string
	WritePath     string
	TelemetryPath string
}

// initConfig initializes all configuration settings
func initConfig() {
	//
}

func parseSendTimeout(value string) time.Duration {
	tm, err := time.ParseDuration(value)
	if err != nil {
		log.ErrorObj(err).Str("timeout-value", value).Msg("Invalid timeout duration provided. Using default timeout.")
		return defaultTimeout
	}

	return tm
}
