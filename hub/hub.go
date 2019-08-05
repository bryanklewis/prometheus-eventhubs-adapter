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
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-amqp-common-go/v2/aad"
	"github.com/Azure/azure-amqp-common-go/v2/sas"
	eventhub "github.com/Azure/azure-event-hubs-go/v2"
	"github.com/Azure/go-autorest/autorest/azure"

	//"github.com/prometheus/common/model"

	"github.com/bryanklewis/prometheus-eventhubs-adapter/log"
	"github.com/bryanklewis/prometheus-eventhubs-adapter/serializers"
)

// EventHubConfig for an Event Hub
type EventHubConfig struct {
	Namespace     string
	Hub           string
	KeyName       string
	KeyValue      string
	ConnString    string
	TenantID      string
	ClientID      string
	ClientSecret  string
	CertPath      string
	CertPassword  string
	Batch         bool
	BatchMaxBytes int
	ADXMapping    string
	Serializer    serializers.SerializerConfig
}

// EventHubClient sends Prometheus samples to Event Hubs
type EventHubClient struct {
	Hub           *eventhub.Hub
	runtimeInfo   *eventhub.HubRuntimeInformation
	batch         bool
	batchMaxBytes int
	adxMapping    string
	serializer    serializers.Serializer
}

// NewClient creates a new event hub client
func NewClient(cfg *EventHubConfig) (*EventHubClient, error) {
	hb, err := newHubFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	rt, err := hb.GetRuntimeInformation(ctx)
	if err != nil {
		return nil, err
	}

	ser, err := serializers.NewSerializer(&cfg.Serializer)
	if err != nil {
		return nil, err
	}

	client := &EventHubClient{
		Hub:           hb,
		runtimeInfo:   rt,
		adxMapping:    cfg.ADXMapping,
		batch:         cfg.Batch,
		batchMaxBytes: cfg.BatchMaxBytes,
		serializer:    ser,
	}

	return client, nil
}

// Name identifies the client path
func (c *EventHubClient) Name() string {
	return c.runtimeInfo.Path
}

// newHubFromConfig returns an event hub instance creation function based on the configuration options provided
//
// Based on (github.com/Azure/azure-event-hubs-go/v2) NewHubWithNamespaceNameAndEnvironment(),
// but uses a local config instead of environment variables
func newHubFromConfig(cfg *EventHubConfig) (*eventhub.Hub, error) {
	if cfg.ConnString != "" {
		return eventhub.NewHubFromConnectionString(cfg.ConnString)
	}

	if cfg.Namespace != "" && cfg.Hub != "" {
		if cfg.KeyName != "" && cfg.KeyValue != "" {
			provider, sasErr := sas.NewTokenProvider(sas.TokenProviderWithKey(cfg.KeyName, cfg.KeyValue))
			if sasErr == nil {
				return eventhub.NewHub(cfg.Namespace, cfg.Hub, provider)
			}
			log.ErrorObj(sasErr).Msg("failure creating SAS token provider")
		}

		if cfg.TenantID != "" && cfg.ClientID != "" {
			if cfg.ClientSecret != "" {
				provider, aadErr := aad.NewJWTProvider(jwtProviderFromConfig(*cfg))
				if aadErr == nil {
					return eventhub.NewHub(cfg.Namespace, cfg.Hub, provider)
				}
				log.ErrorObj(aadErr).Msg("failure creating AAD token provider with client secret")
			}

			if cfg.CertPath != "" && cfg.CertPassword != "" {
				provider, aadErr := aad.NewJWTProvider(jwtProviderFromConfig(*cfg))
				if aadErr == nil {
					return eventhub.NewHub(cfg.Namespace, cfg.Hub, provider)
				}
				log.ErrorObj(aadErr).Msg("failure creating AAD token provider with certificate")
			}
		}
	}

	return nil, errors.New("unable to determine event hub creation; missing configuration parameter")
}

// jwtProviderFromConfig provides an aad.JWTProviderOption using provided configuration
//
// Based on (github.com/Azure/azure-amqp-common-go/v2/aad) JWTProviderWithEnvironmentVars(),
// but uses a local config instead of environment variables
func jwtProviderFromConfig(cfg EventHubConfig) aad.JWTProviderOption {
	return func(config *aad.TokenProviderConfiguration) error {
		config.TenantID = cfg.TenantID
		config.ClientID = cfg.ClientID
		config.ClientSecret = cfg.ClientSecret
		config.CertificatePath = cfg.CertPath
		config.CertificatePassword = cfg.CertPassword

		config.Env = &azure.PublicCloud

		return nil
	}
}
