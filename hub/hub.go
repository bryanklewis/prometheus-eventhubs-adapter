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
	//"context"
	//"time"

	//"github.com/bryanklewis/prometheus-eventhubs-adapter/log"
	//"github.com/bryanklewis/prometheus-eventhubs-adapter/serializer"

	eventhub "github.com/Azure/azure-event-hubs-go/v2"
	//"github.com/prometheus/common/model"
)

// EventHubConfig for an Event Hub
type EventHubConfig struct {
	Namespace    string
	Hub          string
	KeyName      string
	KeyValue     string
	ConnString   string
	TenantID     string
	ClientID     string
	ClientSecret string
	CertPath     string
	CertPassword string
	Batch        bool
	Serializer   string
}

// Client sends Prometheus samples to Event Hubs
type Client struct {
	Hub *eventhub.Hub
}
