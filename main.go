package main

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
  Copyright 2017 The Prometheus Authors, Apache License 2.0
  Copyright 2019 Timescale, Inc., Apache License 2.0
  Copyright 2018 Telefónica, Apache License 2.0
  Copyright 2018 gin-contrib, MIT License
*/

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"

	"github.com/bryanklewis/prometheus-eventhubs-adapter/log"
)

const (
	// AppName is the application name. Value is static and will not change.
	AppName = "prometheus-eventhubs-adapter"
)

// Build information. Populated at compile-time using -ldflags "-X main.BUILD=value"
var (
	Version string
	Commit  string
	Build   string
)

func main() {
	log.Info().Str("version", Version).Str("commit", Commit).Str("build", Build).Msgf("%s starting", AppName)
	newConfig()

	/*writer := hub.NewClient(hubCfg)
	log.Info().Str("write-encoding", hubCfg.Serializer.ADXFormat()).Msg("created event hub writer")

	// Set GIN_MODE
	if e := log.Debug(); e.Enabled() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Gin router
	router := gin.New()

	// Global handler
	router.Use(logHandler([]string{appCfg.TelemetryPath}), gin.Recovery())

	// Route handlers
	router.POST(appCfg.WritePath, timeHandler("write"), writeHandler(writer))
	router.GET(appCfg.TelemetryPath, gin.WrapH(promhttp.Handler()))

	// HTTP server
	srv := &http.Server{
		Addr:         appCfg.ListenAddr,
		Handler:      router,
		ReadTimeout:  appCfg.RemoteTimeout,
		WriteTimeout: appCfg.RemoteTimeout,
	}

	log.Info().Msgf("listening and serving HTTP on %s", srv.Addr)
	go func() {
		// serve connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server listen error")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	// Send incomming quit signals to channel
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal on the channel
	<-quit

	log.Info().Msg("received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), (appCfg.RemoteTimeout + (3 * time.Second)))
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}

	// Close event hub client
	if err := writer.Hub.Close(ctx); err != nil {
		log.Error().Err(err).Msg("close event hub error")
	}*/

	log.Info().Str("version", Version).Str("commit", Commit).Str("build", Build).Msgf("%s exiting", AppName)
}

type writer interface {
	Write(ctx context.Context, samples model.Samples) error
	Name() string
}

// logHandler initializes a gin logging middleware.
//
// Used by the global router.Use() to generate a combined HTTP access and error log.
// An array of routes (example: []string{"/metrics", "/skip"}) can be passed to
// exclude a route from logging.
func logHandler(skipPaths ...[]string) gin.HandlerFunc {
	var newSkipPaths []string
	if len(skipPaths) > 0 {
		newSkipPaths = skipPaths[0]
	}

	var skip map[string]struct{}
	if length := len(newSkipPaths); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range newSkipPaths {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()

		if _, ok := skip[path]; !ok {
			latency := time.Since(start).Nanoseconds() / int64(time.Millisecond)

			msg := "Request"
			if len(c.Errors) > 0 {
				msg = c.Errors.String()
			}

			requestLogger := log.Logger.With().
				Int("status", c.Writer.Status()).
				Str("method", c.Request.Method).
				Str("path", path).
				Str("ip", c.ClientIP()).
				Int64("latency_ms", latency).
				Str("user-agent", c.Request.UserAgent()).
				Logger()

			switch {
			case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
				{
					requestLogger.Warn().Msg(msg)
				}
			case c.Writer.Status() >= http.StatusInternalServerError:
				{
					requestLogger.Error().Msg(msg)
				}
			default:
				requestLogger.Debug().Msg(msg)
			}
		}
	}
}

// timeHandler initializes a gin middleware to track HTTP request time
//
// To allow tracking of different routes, timeHandler is intentionally not set
// in the global gin router.Use(). Instead, each route exposed with router.VERB
// should list this middleware first, then the desired application handler.
// Uses Prometheus histogram to track time.
func timeHandler(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		elapsedMs := time.Since(start).Nanoseconds() / int64(time.Millisecond)
		httpRequestDuration.WithLabelValues(path).Observe(float64(elapsedMs))
	}
}

// writeHandler send to Event Hubs
func writeHandler(writer writer) func(c *gin.Context) {
	return func(c *gin.Context) {
		httpRequestsTotal.Add(float64(1))

		compressed, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			log.ErrorObj(err).Msg("read request body failed")
			return
		}

		reqBuf, err := snappy.Decode(nil, compressed)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			log.ErrorObj(err).Msg("decompress request body failed")
			return
		}

		var req prompb.WriteRequest
		if err := proto.Unmarshal(reqBuf, &req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			log.ErrorObj(err).Msg("unmarshal request body failed")
			return
		}

		samples := protoToSamples(&req)
		receivedSamples.Add(float64(len(samples)))

		//ctx, cancel := context.WithCancel(c)
		//defer cancel()

		/*if err := sendSamples(ctx, writer, samples); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			log.ErrorObj(err).Int("num_samples", len(samples)).Msg("Error sending samples to remote storage")
			return
		}*/

		counter, err := sentSamples.GetMetricWithLabelValues(writer.Name())
		if err != nil {
			log.ErrorObj(err).Str("labelValue", writer.Name()).Msg("Couldn't get a counter")
		}
		writeThroughput.SetCurrent(getCounterValue(counter))

		select {
		case d := <-writeThroughput.Values:
			log.Info().Float64("samples_per_sec", d).Msg("Samples write throughput")
		default:
		}
	}
}

// getCounterValue returns a prometheus counter value
func getCounterValue(counter prometheus.Counter) float64 {
	dtoMetric := &io_prometheus_client.Metric{}
	if err := counter.Write(dtoMetric); err != nil {
		log.ErrorObj(err).Msg("Error reading counter value")
	}
	return dtoMetric.GetCounter().GetValue()
}

// protoToSamples converts a Prometheus protobuf WriteRequest to Prometheus Samples
func protoToSamples(req *prompb.WriteRequest) model.Samples {
	var samples model.Samples
	for _, ts := range req.Timeseries {
		metric := make(model.Metric, len(ts.Labels))
		for _, l := range ts.Labels {
			metric[model.LabelName(l.Name)] = model.LabelValue(l.Value)
		}

		for _, s := range ts.Samples {
			samples = append(samples, &model.Sample{
				Metric:    metric,
				Value:     model.SampleValue(s.Value),
				Timestamp: model.Time(s.Timestamp),
			})
		}
	}

	return samples
}

/*
func sendSamples(ctx context.Context, w writer, samples model.Samples) error {
	begin := time.Now()

	err := w.Write(ctx, samples)

	duration := time.Since(begin).Seconds()
	if err != nil {
		failedSamples.WithLabelValues(w.Name()).Add(float64(len(samples)))
		return err
	}

	sentSamples.WithLabelValues(w.Name()).Add(float64(len(samples)))
	sentBatchDuration.WithLabelValues(w.Name()).Observe(duration)

	return nil
}
*/
