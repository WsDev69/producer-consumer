package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Prometheus struct {
	reg *prometheus.Registry
}

func New(reg *prometheus.Registry, cs ...prometheus.Collector) *Prometheus {
	for i := range cs {
		reg.MustRegister(cs[i])
	}
	return &Prometheus{reg: reg}
}

func (p *Prometheus) Serve(addr string) error {
	http.Handle("/metrics", promhttp.HandlerFor(
		p.reg,
		promhttp.HandlerOpts{},
	))
	return http.ListenAndServe(addr, nil)
}
