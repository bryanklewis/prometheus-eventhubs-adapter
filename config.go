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

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/bryanklewis/prometheus-eventhubs-adapter/hub"
	"github.com/bryanklewis/prometheus-eventhubs-adapter/log"
	"github.com/bryanklewis/prometheus-eventhubs-adapter/serializers"
)

// config represents settings for the application
type config struct {
	readTimeout   time.Duration
	writeTimeout  time.Duration
	listenAddress string
	writePath     string
	telemetryPath string
	filterType    int
	filterBy      string
	logLevel      string
	writeHub      hub.EventHubConfig
}

var (
	adapterConfig = &config{}
)

// parseFlags parses the adapter configuration flags
func parseFlags() {
	// Main application
	flag.DurationVar(&adapterConfig.readTimeout, "read_timeout", 5*time.Second, "HTTP server time from when the connection is accepted to when the request body is fully read.")
	viper.SetDefault("read_timeout", 5*time.Second)

	flag.DurationVar(&adapterConfig.writeTimeout, "write_timeout", 10*time.Second, "HTTP server time from the end of the request header read to the end of the response write.")
	viper.SetDefault("write_timeout", 10*time.Second)

	flag.StringVar(&adapterConfig.listenAddress, "listen_address", ":9201", "Address to listen on for web endpoints.")
	viper.SetDefault("listen_address", ":9201")

	flag.StringVar(&adapterConfig.writePath, "write_path", "/write", "Path for write requests.")
	viper.SetDefault("write_path", "/write")

	flag.StringVar(&adapterConfig.telemetryPath, "telemetry_path", "/metrics", "Path for telemetry scraps.")
	viper.SetDefault("telemetry_path", "/metrics")

	flag.IntVar(&adapterConfig.filterType, "filter_type", 0, "0 (none), 1 (whitelist), 2 (blacklist).")
	viper.SetDefault("filter_type", 0)

	flag.StringVar(&adapterConfig.filterBy, "filter_by", "", "Names of metric(s) to filter by, comma-separated.")
	viper.SetDefault("filter_by", "")

	flag.StringVar(&adapterConfig.logLevel, "log_level", "info", "The log level to use [ \"error\", \"warn\", \"info\", \"debug\", \"none\" ].")
	viper.SetDefault("log_level", "info")

	// Event Hub Writer
	flag.StringVar(&adapterConfig.writeHub.Namespace, "write_namespace", "", "Namespace of the Event Hub instance.")

	flag.StringVar(&adapterConfig.writeHub.Hub, "write_hub", "", "Name of the Event Hub.")

	flag.StringVar(&adapterConfig.writeHub.KeyName, "write_keyname", "", "Name of the Event Hub key.")

	flag.StringVar(&adapterConfig.writeHub.KeyValue, "write_keyvalue", "", "Secret for the corresponding \"write_keyname\".")

	flag.StringVar(&adapterConfig.writeHub.ConnString, "write_connstring", "", "Connection string from the Azure portal.")

	flag.StringVar(&adapterConfig.writeHub.TenantID, "write_tenantid", "", "Azure Active Directory Tenant ID.")

	flag.StringVar(&adapterConfig.writeHub.ClientID, "write_clientid", "", "Azure Active Directory Client ID or Application ID.")

	flag.StringVar(&adapterConfig.writeHub.ClientSecret, "write_clientsecret", "", "Secret for the corresponding \"write_clientid\".")

	flag.StringVar(&adapterConfig.writeHub.CertPath, "write_certpath", "", "Path to the certificate file.")

	flag.StringVar(&adapterConfig.writeHub.CertPassword, "write_certpassword", "", "Password for the certificate.")

	flag.BoolVar(&adapterConfig.writeHub.Batch, "write_batch", true, "Send batch events or single events.")
	viper.SetDefault("write_batch", true)

	flag.StringVar(&adapterConfig.writeHub.ADXMapping, "write_adxmapping", "promMap", "Azure Data Explorer data injestion mapping name.")
	viper.SetDefault("write_adxmapping", "promMap")

	// Valid values can be found in serializers.NewSerializer
	flag.StringVar(&adapterConfig.writeHub.Serializer.DataFormat, "write_serializer", "json", "Serializer to use when sending events [ \"json\", \"avro-json\" ].")
	viper.SetDefault("write_serializer", "json")
}

// initConfig initializes configuration setup
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
		viper.AddConfigPath("/etc/" + AppName)
	}

	// Config file search: current working directory
	if workingpath, err := os.Getwd(); err != nil {
		log.Debug().Err(err).Msg("failed to detect working directory")
	} else {
		viper.AddConfigPath(workingpath)
	}

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore since its optional
			log.Debug().Msg("configuration file not detected (optional)")
		} else {
			// Config file was found but error was produced
			log.Error().Err(err).Msg("Error loading config file")
		}
	}
	log.Debug().Msg("configuration file detected (optional)")

	// Config environment vars
	// Will check for a environment variable with a name matching the key
	// uppercased and prefixed with the EnvPrefix if set.
	viper.SetEnvPrefix("adap")
	viper.AutomaticEnv()

	// Config commandline flags and defaults
	parseFlags()

	// Viper uses "pflag", add standard library "flag" to "pflag"
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.ErrorObj(err).Msg("Failed to bind pflags to config")
	}

	// Auto-config
	configLogging()
}

// configLogging configures logging
func configLogging() {
	// Set global logging preference
	if err := log.SetLevel(viper.GetString("log_level")); err != nil {
		log.ErrorObj(err).Msg("Invalid log level provided")
	}

	// Show config when debugging
	debugConfig := viper.AllSettings()
	log.Debug().Fields(debugConfig).Msg("show config")
}

// getWriterConfig returns the configuration for an Event Hub Writer
func getWriterConfig() *hub.EventHubConfig {
	return &hub.EventHubConfig{
		Namespace:    viper.GetString("write_namespace"),
		Hub:          viper.GetString("write_hub"),
		KeyName:      viper.GetString("write_keyname"),
		KeyValue:     viper.GetString("write_keyvalue"),
		ConnString:   viper.GetString("write_connstring"),
		TenantID:     viper.GetString("write_tenantid"),
		ClientID:     viper.GetString("write_clientid"),
		ClientSecret: viper.GetString("write_clientsecret"),
		CertPath:     viper.GetString("write_certpath"),
		CertPassword: viper.GetString("write_certpassword"),
		Batch:        viper.GetBool("write_batch"),
		ADXMapping:   viper.GetString("write_adxmapping"),
		Serializer:   serializers.SerializerConfig{DataFormat: viper.GetString("write_serializer")},
	}
}
