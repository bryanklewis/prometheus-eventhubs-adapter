package kusto

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
	"fmt"
	"strings"
)

// DataFormat is an enum for Azure Data Explorer (Kusto) data injestion formats.
//
// The possible formats are found in the Azure documentation.
// See [ https://docs.microsoft.com/en-us/azure/kusto/management/data-ingestion/#supported-data-formats ]
type DataFormat uint8

const (
	// CSVFormat defines the csv data format.
	CSVFormat DataFormat = iota
	// JSONFormat defines the json data format.
	JSONFormat
	// AVROFormat defines the avro data format.
	AVROFormat
	// NoFormat defines an absent format.
	NoFormat
)

func (km DataFormat) String() string {
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

// ParseDataFormat converts a Kusto Data Format string into a kusto.DataFormat value.
// returns an error if the input string does not match known values.
func ParseDataFormat(formatStr string) (DataFormat, error) {
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
