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

	"github.com/bryanklewis/prometheus-eventhubs-adapter/hub"
	"github.com/bryanklewis/prometheus-eventhubs-adapter/log"
)

// config represents settings for the application
type config struct {
	connectionTimeout time.Duration
	listenAddress     string
	writePath         string
	telemetryPath     string
	logLevel          string
	writeHub          hub.EventHubConfig
}

var (
	adapterConfig = &config{}
)

// parseFlags parses the adapter configuration flags
func parseFlags() {
	// Main application
	flag.DurationVar(&adapterConfig.connectionTimeout, "connection_timeout", 5*time.Second, "The timeout to use when sending samples to the remote storage.")
	viper.SetDefault("connection_timeout", 5*time.Second)

	flag.StringVar(&adapterConfig.listenAddress, "listen_address", ":9201", "Address to listen on for web endpoints.")
	viper.SetDefault("listen_address", ":9201")

	flag.StringVar(&adapterConfig.writePath, "write_path", "/write", "Path for write requests.")
	viper.SetDefault("write_path", "/write")

	flag.StringVar(&adapterConfig.telemetryPath, "telemetry_path", "/metrics", "Path for telemetry scraps.")
	viper.SetDefault("telemetry_path", "/metrics")

	flag.StringVar(&adapterConfig.logLevel, "log_level", "info", "The log level to use [ \"error\", \"warn\", \"info\", \"debug\" ].")
	viper.SetDefault("log_level", "info")

	// Event Hub Writer
	flag.StringVar(&adapterConfig.writeHub.Namespace, "write_namespace", "", "Namespace of the Event Hub instance.")
	viper.RegisterAlias("write_namespace", "EVENTHUB_NAMESPACE")

	flag.StringVar(&adapterConfig.writeHub.Hub, "write_hub", "", "Name of the Event Hub.")
	viper.RegisterAlias("write_hub", "EVENTHUB_NAME")

	flag.StringVar(&adapterConfig.writeHub.KeyName, "write_keyname", "", "Name of the Event Hub key.")
	viper.RegisterAlias("write_keyname", "EVENTHUB_KEY_NAME")

	flag.StringVar(&adapterConfig.writeHub.KeyValue, "write_keyvalue", "", "Secret for the corresponding \"write_keyname\".")
	viper.RegisterAlias("write_keyvalue", "EVENTHUB_KEY_VALUE")

	flag.StringVar(&adapterConfig.writeHub.ConnString, "write_connstring", "", "Connection string from the Azure portal.")
	viper.RegisterAlias("write_connstring", "EVENTHUB_CONNECTION_STRING")

	flag.StringVar(&adapterConfig.writeHub.TenantID, "write_tenantid", "", "Azure Active Directory Tenant ID.")
	viper.RegisterAlias("write_tenantid", "AZURE_TENANT_ID")

	flag.StringVar(&adapterConfig.writeHub.ClientID, "write_clientid", "", "Azure Active Directory Client ID or Application ID.")
	viper.RegisterAlias("write_clientid", "AZURE_CLIENT_ID")

	flag.StringVar(&adapterConfig.writeHub.ClientSecret, "write_clientsecret", "", "Secret for the corresponding \"write_clientid\".")
	viper.RegisterAlias("write_clientsecret", "AZURE_CLIENT_SECRET")

	flag.StringVar(&adapterConfig.writeHub.CertPath, "write_certpath", "", "Path to the certificate file.")
	viper.RegisterAlias("write_certpath", "AZURE_CERTIFICATE_PATH")

	flag.StringVar(&adapterConfig.writeHub.CertPassword, "write_certpassword", "", "Password for the certificate.")
	viper.RegisterAlias("write_certpassword", "AZURE_CERTIFICATE_PASSWORD")

	flag.BoolVar(&adapterConfig.writeHub.Batch, "write_batch", true, "Send batch events or single events.")
	viper.SetDefault("write_batch", true)

	// Valid values can be found in serializers.NewSerializer
	flag.StringVar(&adapterConfig.writeHub.Serializer.DataFormat, "write_serializer", "json", "Serializer to use when sending events [ \"json\", \"avro-json\" ].")
	viper.SetDefault("write_serializer", "json")
}

// newConfig initializes configuration setup
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
	viper.BindPFlags(pflag.CommandLine)

	// Show config when debugging
	debugConfig := viper.AllSettings()
	log.Debug().Fields(debugConfig).Msg("show config")
}
