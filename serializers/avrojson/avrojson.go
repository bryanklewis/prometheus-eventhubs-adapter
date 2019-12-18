package avrojson

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
	"time"

	"github.com/linkedin/goavro/v2"
	"github.com/prometheus/common/model"

	"github.com/bryanklewis/prometheus-eventhubs-adapter/kusto"
)

const (
	// SCHEMA is the avro schema used for serialization
	SCHEMA = `{
		"namespace": "io.prometheus",
		"type": "record",
		"name": "Metric",
		"doc:" : "A basic schema for representing Prometheus metrics",
		"fields": [
			{"name": "timestamp", "type": "string"},
			{"name": "value", "type": "double"},
			{"name": "name", "type": "string"},
			{"name": "labels", "type": { "type": "map", "values": "string"} }
		]
	}`
)

// Serializer represents a serializer instance
type Serializer struct {
	Codec *goavro.Codec
}

// ADXFormat Azure Data Explorer injestion data format.
//
// Implements the serializers.Serializer interface
func (s *Serializer) ADXFormat() kusto.DataFormat {
	return kusto.JSONFormat
}

// Serialize takes a single Prometheus sample and turns it into a byte buffer.
//
// Implements the serializers.Serializer interface
func (s *Serializer) Serialize(sample model.Sample) ([]byte, error) {
	m := s.createObject(sample)
	return s.Codec.TextualFromNative(nil, m)
}

func (s *Serializer) createObject(sample model.Sample) map[string]interface{} {
	metricName := sample.Metric[model.MetricNameLabel]

	// Remove sample name from labels set
	labels := make(map[string]string, (len(sample.Metric) - 1))
	for label, value := range sample.Metric {
		if label != model.MetricNameLabel {
			labels[string(label)] = string(value)
		}
	}

	m := map[string]interface{}{
		"timestamp": time.Unix(0, sample.Timestamp.UnixNano()).UTC().Format(time.RFC3339),
		"value":     float64(sample.Value),
		"name":      string(metricName),
		"labels":    labels,
	}

	return m
}
