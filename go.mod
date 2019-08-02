module github.com/bryanklewis/prometheus-eventhubs-adapter

go 1.12

require (
	github.com/Azure/azure-event-hubs-go/v2 v2.0.0
	github.com/gin-gonic/gin v1.4.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/snappy v0.0.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90
	github.com/prometheus/common v0.6.0
	github.com/prometheus/prometheus v2.5.0+incompatible
	github.com/rs/zerolog v1.14.3
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	google.golang.org/genproto v0.0.0-20190708153700-3bdd9d9f5532 // indirect
)
