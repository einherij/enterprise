package webservers

import (
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/einherij/enterprise/webtools"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsConfig struct {
	Port string `mapstructure:"port"`
}

func NewMetricServer(cfg MetricsConfig) (*webtools.Server, error) {
	promHTTPHandler := promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
			ErrorLog: logrus.StandardLogger().WithError(errors.New("prometheus handler error")),
		}),
	)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promHTTPHandler)

	return webtools.NewServer("metrics", cfg.Port, mux)
}
