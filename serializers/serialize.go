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
	//"math"
	//"time"
	"fmt"
	"strings"

	"github.com/prometheus/common/model"
	//"github.com/bryanklewis/prometheus-eventhubs-adapter/log"
)

// KustoFormat is an enum for Azure Data Explorer (Kusto) data injestion formats.
//
// The possible formats are found in the Azure documentation.
// See [ https://docs.microsoft.com/en-us/azure/kusto/management/data-ingestion/#supported-data-formats ]
type KustoFormat uint8

const (
	// CSVFormat defines the csv data format.
	CSVFormat KustoFormat = iota
	// JSONFormat defines the json data format.
	JSONFormat
	// AVROFormat defines the avro data format.
	AVROFormat
	// NoFormat defines an absent format.
	NoFormat
)

func (km KustoFormat) String() string {
	switch km {
	case CSVFormat:
		return "csv"
	case JSONFormat:
		return "json"
	case AVROFormat:
		return "avro"
	default:
		return ""
	}
}

// ParseKustoFormat converts a Kusto Data Format string into a KustoFormat value.
// returns an error if the input string does not match known values.
func ParseKustoFormat(formatStr string) (KustoFormat, error) {
	switch strings.ToLower(formatStr) {
	case "csv":
		return CSVFormat, nil
	case "json":
		return JSONFormat, nil
	case "avro":
		return AVROFormat, nil
	default:
		return NoFormat, fmt.Errorf("Unknown Kusto Data Format: '%s'", strings.ToLower(formatStr))
	}
}

// Serializer is an interface defining functions that a serializer must satisfy.
type Serializer interface {
	// Serialize takes a single Prometheus sample and turns it into a byte buffer.
	Serialize(metric model.Sample) ([]byte, error)

	// SerializeBatch takes an array of Prometheus samples and serializes it into
	// a byte buffer.
	SerializeBatch(metrics model.Samples) ([]byte, error)

	// ADXFormat Azure Data Explorer data injestion format.
	ADXFormat() KustoFormat
}

// SerializerConfig is a struct that covers the data types needed for all serializer types,
// and can be used to instantiate _any_ of the serializers.
type SerializerConfig struct {
	// Dataformat can be one of the serializer types listed in NewSerializer.
	DataFormat string

	// Additional fields can be defined here to pass to a specific serializer.
}

/*
// NewSerializer provides a Serializer based on the given config.
func NewSerializer(cfg *SerializerConfig) (Serializer, error) {
	var err error
	var serializer Serializer
	switch cfg.DataFormat {
	case "json":
		serializer, err = newJsonSerializer()
	case "avro-json":
		serializer, err = newAvroJsonSerializer()
	default:
		err = fmt.Errorf("Invalid data format: %s", config.DataFormat)
	}
	return serializer, err
}


func newJsonSerializer() (Serializer, error) {
	return json.NewSerializer()
}
*/
