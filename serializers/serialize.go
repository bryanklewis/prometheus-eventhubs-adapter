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

	"github.com/prometheus/common/model"
	//"github.com/bryanklewis/prometheus-eventhubs-adapter/log"
)

// Serializer is an interface defining functions that a serializer must satisfy.
type Serializer interface {
	// Serialize takes a single Prometheus sample and turns it into a byte buffer.
	Serialize(metric model.Sample) ([]byte, error)

	// SerializeBatch takes an array of Prometheus samples and serializes it into
	// a byte buffer.
	SerializeBatch(metrics model.Samples) ([]byte, error)

	// Azure Data Explorer injestion format.
	ADXFormat() string
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
