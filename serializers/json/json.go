package json

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
	"encoding/json"
	"time"

	"github.com/prometheus/common/model"

	"github.com/bryanklewis/prometheus-eventhubs-adapter/kusto"
)

const (
	defaultMetricName model.LabelValue = "no_name"
)

// Serializer represents a serializer instance
type Serializer struct {
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
	serialized, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}

	return serialized, nil
}

func (s *Serializer) createObject(sample model.Sample) map[string]interface{} {
	var metricName model.LabelValue
	var hasName bool
	metricName, hasName = sample.Metric[model.MetricNameLabel]
	numLabels := len(sample.Metric) - 1
	if !hasName {
		numLabels = len(sample.Metric)
		metricName = defaultMetricName
	}

	// Remove sample name from labels set
	labels := make(map[string]string, numLabels)
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
