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
	"os"
	"path/filepath"
	"runtime"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

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

// newConfig initializes all configuration settings
func newConfig() {
	// Config file
	viper.SetConfigName(AppName)
	viper.SetConfigType("toml")

	// Config file search: adapter executable directory
	if ex, err := os.Executable(); err != nil {
		log.Debug().Err(err).Msg("failed to detect executable directory")
	} else {
		viper.AddConfigPath(filepath.Dir(ex))
	}

	// Config file search: Unix-like system configuration directory
	if runtime.GOOS != "windows" {
		viper.AddConfigPath("/etc/" + AppName + "/")
	}

	// Config file search: current working directory
	if workingpath, err := os.Getwd(); err != nil {
		log.Debug().Err(err).Msg("failed to detect working directory")
	} else {
		viper.AddConfigPath(workingpath)
	}

	// Config file search: OS-specific home directory
	if homepath, err := homedir.Dir(); err != nil {
		log.Debug().Err(err).Msg("failed to detect home directory")
	} else {
		viper.AddConfigPath(homepath)
	}

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore since its optional
			log.Debug().Msg("configuration file not detected (optional)")
		} else {
			// Config file was found but error was produced
			log.Panic().Err(err).Msg("Fatal error loading config file \n")
		}
	}
	log.Debug().Msg("configuration file detected (optional)")
}

func parseSendTimeout(value string) time.Duration {
	tm, err := time.ParseDuration(value)
	if err != nil {
		log.ErrorObj(err).Str("timeout-value", value).Msg("Invalid timeout duration provided. Using default timeout.")
		return defaultTimeout
	}

	return tm
}
