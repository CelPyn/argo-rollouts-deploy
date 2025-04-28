package main

import (
	"github.com/CelPyn/argo-rollouts-deploy/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
)

var (
	invokedByPath = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "req_count_by_path",
		Help: "The total number of incoming requests by path",
	}, []string{"path", "method"})
	successByPath = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "success_count_by_path",
		Help: "The total number of incoming requests by path with success status",
	}, []string{"path", "method", "code"})
	errByPath = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "err_count_by_path",
		Help: "The total number of incoming requests by path with error status",
	}, []string{"path", "method", "code"})
)

func main() {
	router := http.DefaultServeMux

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	registerRoutes(router)

	log.Println("starting server, listening on", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func registerRoutes(router *http.ServeMux) {
	log.Println("registering routes")
	router.Handle("GET /metrics", promhttp.Handler())

	router.HandleFunc("GET /health/liveness", api.Health(meter()))
	router.HandleFunc("GET /health/readiness", api.Health(meter()))

	router.HandleFunc("GET /", api.ColourHTML(meter()))
	router.HandleFunc("GET /json", api.ColourJSON(meter()))
}

func meter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invokedByPath.WithLabelValues(r.URL.Path, r.Method).Inc()

		statusCodeCtx := *api.StatusCodeCtx(r.Context())

		if statusCodeCtx < 400 {
			successByPath.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(int(statusCodeCtx))).Inc()
		} else {
			errByPath.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(int(statusCodeCtx))).Inc()
		}

		get := r.Header.Get("X-Forwarded-For")
		if get == "" {
			get = r.RemoteAddr
		}
		log.Println(r.Method, r.URL.Path, "-", get, r.UserAgent(), "-", statusCodeCtx)
	}
}
