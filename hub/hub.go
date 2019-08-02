package hub

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
	//"context"
	//"time"

	//"github.com/bryanklewis/prometheus-eventhubs-adapter/log"
	//"github.com/bryanklewis/prometheus-eventhubs-adapter/serializer"

	eventhub "github.com/Azure/azure-event-hubs-go/v2"
	"github.com/spf13/viper"
	//"github.com/prometheus/common/model"
)

// EventHubConfig for an Event Hub
type EventHubConfig struct {
	batch bool
	// Serializer to use when sending or recieving events
	serializer string
	//serializer serialize.Serializer
}

// Client sends Prometheus samples to Event Hubs
type Client struct {
	Hub *eventhub.Hub
	cfg *EventHubConfig
}

// ParseFlagsWriter gets flags specific to Event Hubs writer
func ParseFlagsWriter(cfg *EventHubConfig) *EventHubConfig {
	flag.BoolVar(&cfg.batch, "write_batch", true, "Send batch events or single events.")
	viper.SetDefault("write_batch", true)

	flag.StringVar(&cfg.serializer, "write_serializer", "json", "Serializer to use when sending events [ \"json\", \"json-avro\" ].")
	viper.SetDefault("write_serializer", "json")

	return cfg
}
