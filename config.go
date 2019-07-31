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
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/bryanklewis/prometheus-eventhubs-adapter/log"
)

// initConfig initializes all configuration settings
func initConfig() {
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

	// Config environment vars
	// Will check for a environment variable with a name matching the key
	// uppercased and prefixed with the EnvPrefix if set.
	viper.SetEnvPrefix("adap")
	viper.AutomaticEnv()

	// Config commandline flags and defaults
	flag.Duration("connection_timeout", 5*time.Second, "The timeout to use when sending samples to the remote storage.")
	viper.SetDefault("connection_timeout", 5*time.Second)

	flag.String("listen_address", ":9201", "Address to listen on for web endpoints.")
	viper.SetDefault("listen_address", ":9201")

	flag.String("write_path", "/write", "Path for write requests.")
	viper.SetDefault("write_path", "/write")

	flag.String("telemetry_path", "/metrics", "Path for telemetry scraps.")
	viper.SetDefault("telemetry_path", "/metrics")

	flag.String("log_level", "info", "The log level to use [ \"error\", \"warn\", \"info\", \"debug\" ].")
	viper.SetDefault("log_level", "info")

	flag.String("send_encoding", "json", "Encoding to use when sending events [ \"json\", \"avro-json\" ].")
	viper.SetDefault("send_encoding", "json")

	flag.Bool("send_batch", true, "Send batch events or single events.")
	viper.SetDefault("send_batch", true)

	// Viper uses "pflag", add standard library "flag" to "pflag"
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}
