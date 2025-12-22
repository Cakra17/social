package prom

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusClient struct {
	counterVec []*prometheus.CounterVec
	registry *prometheus.Registry
}

func NewPrometheusService() *PrometheusClient {
	var cv []*prometheus.CounterVec
	HttpRequestTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_total",
		Help: "Total number of requests proceed by the server",
	}, []string{"path", "status"})

	HttpRequestErrorTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_error_total",
		Help: "Total number of error returned by the server",
	}, []string{"path", "status"})

	customRegistry := prometheus.NewRegistry()

	cv = append(cv, HttpRequestTotal, HttpRequestErrorTotal)
	return &PrometheusClient{
		counterVec: cv,
		registry: customRegistry,
	}
}

func (p *PrometheusClient) Register() {
	if len(p.counterVec) < 1 {
		log.Fatalf("counterVec is empty")
	}

	for _, v := range p.counterVec {
		p.registry.MustRegister(v)
	}
}

func (p *PrometheusClient) Handler() http.HandlerFunc {
	h := promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{})
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

func (p *PrometheusClient) RequestMetricMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		wrap := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(wrap, r)
		status := wrap.Status()
		if status < 400 {
			p.counterVec[0].WithLabelValues(path, strconv.Itoa(status)).Inc()
		} else {
			p.counterVec[1].WithLabelValues(path, strconv.Itoa(status)).Inc()
		}
	}
	return http.HandlerFunc(fn)
}