package serializers

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
  Copyright (c) 2015-2019 InfluxData Inc., MIT License
  Copyright 2018 Telefónica, Apache License 2.0
*/

import (
	"fmt"
	"strings"

	"github.com/linkedin/goavro/v2"
	"github.com/prometheus/common/model"

	"github.com/bryanklewis/prometheus-eventhubs-adapter/kusto"
	"github.com/bryanklewis/prometheus-eventhubs-adapter/log"
	"github.com/bryanklewis/prometheus-eventhubs-adapter/serializers/avrojson"
	"github.com/bryanklewis/prometheus-eventhubs-adapter/serializers/json"
)

// Serializer is an interface defining functions that a serializer must satisfy.
type Serializer interface {
	// Serialize takes a single Prometheus sample and turns it into a byte buffer.
	Serialize(metric model.Sample) ([]byte, error)

	// ADXFormat Azure Data Explorer injestion data format.
	ADXFormat() kusto.DataFormat
}

// SerializerConfig is a struct that covers the data types needed for all serializer types,
// and can be used to instantiate _any_ of the serializers.
type SerializerConfig struct {
	// Dataformat can be one of the serializer types listed in serializers.NewSerializer.
	DataFormat string
}

// NewSerializer provides a Serializer based on the given config.
//
// Parses SerializerConfig.DataFormat string
func NewSerializer(cfg *SerializerConfig) (Serializer, error) {
	switch strings.ToLower(cfg.DataFormat) {
	case "json":
		return NewJSONSerializer()
	case "avro-json":
		return NewAvroJSONSerializer()
	default:
		err := fmt.Errorf("Invalid data format: %s", strings.ToLower(cfg.DataFormat))
		return nil, err
	}
}

// NewJSONSerializer provides a 'json' Serializer
func NewJSONSerializer() (Serializer, error) {
	return &json.Serializer{}, nil
}

// NewAvroJSONSerializer provides a 'avro-json' Serializer
func NewAvroJSONSerializer() (Serializer, error) {
	codec, err := goavro.NewCodec(avrojson.SCHEMA)
	if err != nil {
		log.ErrorObj(err).Msg("Failed to create avro codec")
		return nil, err
	}

	return &avrojson.Serializer{
		Codec: codec,
	}, nil
}
