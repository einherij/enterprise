package servers

import (
	"net/http"
	"net/http/pprof"

	"github.com/einherij/enterprise/httputils"
)

type PProfConfig struct {
	Port string `mapstructure:"port"`
}

func NewPProfServer(cfg PProfConfig) (*httputils.Server, error) {
	return httputils.NewServer("pprof", cfg.Port, newPProfMux())
}

func newPProfMux() *http.ServeMux {
	pprofMux := http.NewServeMux()
	pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
	pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	pprofMux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	pprofMux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	pprofMux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	pprofMux.Handle("/debug/pprof/block", pprof.Handler("block"))
	return pprofMux
}
