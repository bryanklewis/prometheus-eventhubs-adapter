module github.com/bryanklewis/prometheus-eventhubs-adapter

go 1.15

require (
	github.com/Azure/azure-amqp-common-go/v3 v3.1.0
	github.com/Azure/azure-event-hubs-go/v3 v3.3.4
	github.com/Azure/go-autorest/autorest v0.11.13
	github.com/gin-gonic/gin v1.7.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/snappy v0.0.2
	github.com/linkedin/goavro/v2 v2.10.0
	github.com/prometheus/client_golang v1.8.0
	github.com/prometheus/common v0.15.0
	github.com/prometheus/prometheus v2.5.0+incompatible
	github.com/rs/zerolog v1.20.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
)
