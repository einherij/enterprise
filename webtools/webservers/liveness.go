package webservers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/einherij/enterprise/webtools"
)

type LivenessConfig struct {
	LivenessPort string `mapstructure:"liveness_port"`
}

func NewLivenessServer(cfg LivenessConfig) (*webtools.Server, error) {
	srv, err := webtools.NewServer(
		"liveness_probe",
		cfg.LivenessPort,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/health") {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}),
	)
	if err != nil {
		return srv, fmt.Errorf("error creating liveness probe server: %w", err)
	}
	return srv, nil
}
