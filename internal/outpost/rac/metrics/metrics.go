package metrics

import (
	"net/http"

	"go.uber.org/zap"
	"goauthentik.io/internal/config"
	"goauthentik.io/internal/utils/sentry"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunServer() {
	m := mux.NewRouter()
	l := config.Get().Logger().Named("authentik.outpost.metrics")
	m.Use(sentry.SentryNoSampleMiddleware)
	m.HandleFunc("/outpost.goauthentik.io/ping", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(204)
	})
	m.Path("/metrics").Handler(promhttp.Handler())
	listen := config.Get().Listen.Metrics
	l.Info("Starting Metrics server", zap.String("listen", listen))
	err := http.ListenAndServe(listen, m)
	if err != nil {
		l.Warn("Failed to start metrics listener", zap.Error(err))
	}
}
