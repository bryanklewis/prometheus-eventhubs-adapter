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

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bryanklewis/prometheus-eventhubs-adapter/log"
)

// Defaults
const (
	// The timeout to use when sending samples to the remote storage.
	defaultTimeout = 10 * time.Second
	// Address to listen on for web endpoints.
	defaultListen = ":9201"
	// Path for write requests.
	defaultWritePath = "/write"
	// Path for telemetry scraps.
	defaultTelemetry = "/metrics"
	// The log level to use [ \"error\", \"info\", \"debug\" ].
	defaultLogLevel = zerolog.InfoLevel
	// Encoding to use when sending events [ \"json\", \"avro-json\" ].
	defaultEncoding = "json"

	appName      string = "prometheus-eventhubs-adapter"
)


// Config for general application options
type Config struct {
	RemoteTimeout time.Duration
	ListenAddr    string
	WritePath     string
	TelemetryPath string
}

// GetConfig processes all configuration settings
func GetConfig() (*Config, *hub.Config) {
	appCfg := &Config{}

	if value, ok := os.LookupEnv("ADAPTER_SEND_TIMEOUT"); ok {
		appCfg.RemoteTimeout = parseSendTimeout(value)
	} else {
		appCfg.RemoteTimeout = defaultTimeout
	}

	if value, ok := os.LookupEnv("ADAPTER_LISTEN_ADDRESS"); ok {
		appCfg.ListenAddr = value
	} else {
		appCfg.ListenAddr = defaultListen
	}

	if value, ok := os.LookupEnv("ADAPTER_WRITE_PATH"); ok {
		appCfg.WritePath = value
	} else {
		appCfg.WritePath = defaultWritePath
	}

	if value, ok := os.LookupEnv("ADAPTER_TELEMETRY_PATH"); ok {
		appCfg.TelemetryPath = value
	} else {
		appCfg.TelemetryPath = defaultTelemetry
	}

	if value, ok := os.LookupEnv("LOG_LEVEL"); ok {
		zerolog.SetGlobalLevel(parseLogLevel(value))
	} else {
		zerolog.SetGlobalLevel(defaultLogLevel)
	}

	hubCfg := &hub.Config{}

	var serialErr error
	if value, ok := os.LookupEnv("WRITE_ENCODING"); ok {
		hubCfg.Serializer, serialErr = parseEncoding(value)
		if serialErr != nil {
			log.Fatal().Err(serialErr).Str("encoding", value).Msg("Failed to create serializer")
		}
	} else {
		hubCfg.Serializer, serialErr = parseEncoding(defaultEncoding)
		if serialErr != nil {
			log.Fatal().Err(serialErr).Str("encoding", value).Msg("Failed to create serializer")
		}
	}

	return appCfg, hubCfg
}

func parseSendTimeout(value string) time.Duration {
	tm, err := time.ParseDuration(value)
	if err != nil {
		log.ErrorObj(err).Str("timeout-value", value).Msg("Invalid timeout duration provided. Using default timeout.")
		return defaultTimeout
	}

	return tm
}

func parseLogLevel(value string) zerolog.Level {
	level, err := zerolog.ParseLevel(value)
	if err != nil {
		log.ErrorObj(err).Str("log-level-value", value).Msg("Invalid log level provided, using level `info`")
		return defaultLogLevel
	}

	return level
}

func parseEncoding(value string) (serialize.Serializer, error) {
	e := strings.ToLower(value)
	switch e {
	case "avro-json":
		return serialize.NewAvroJSONSerializer()
	case "json":
		return serialize.NewJSONSerializer()
	default:
		return nil, errors.New("encoding value provided is not valid")
	}
}
