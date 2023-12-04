package webservers

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/suite"
)

type MetricServerSuite struct {
	suite.Suite

	ctrl *gomock.Controller
}

func TestMetricServerSuite(t *testing.T) {
	suite.Run(t, new(MetricServerSuite))
}

func (s *MetricServerSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
}

func (s *MetricServerSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *MetricServerSuite) TestNewMetricsHandler() {
	metricsConf := MetricsConfig{
		Port: "0",
	}
	serv, err := NewMetricServer(metricsConf)
	s.NoError(err)
	s.NotEmpty(serv)
}

func (s *MetricServerSuite) TestMetricsHandler() {
	metricsConf := MetricsConfig{
		Port: "0",
	}
	serv, err := NewMetricServer(metricsConf)
	s.NoError(err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go serv.Run(ctx)

	endpoint := "http://" + serv.Addr().String() + "/metrics"
	resp, err := http.Get(endpoint)
	s.NoError(err)
	decoder := expfmt.NewDecoder(resp.Body, expfmt.FmtText)
	mf := &dto.MetricFamily{}
	s.NoError(decoder.Decode(mf))
}
