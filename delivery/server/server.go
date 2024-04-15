package server

import (
	"context"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox-crawler/v2/adapter/storage"
)

type Server struct {
	log     *zerolog.Logger
	srv     *http.Server
	storage *storage.Storage

	stopScraping chan struct{}

	cfg Config
}

func New(cfg Config, s *storage.Storage, l zerolog.Logger) *Server {
	l = l.With().Str("cmp", "server").Logger()

	return &Server{
		log:          &l,
		cfg:          cfg,
		storage:      s,
		stopScraping: make(chan struct{}),
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.srv = &http.Server{
		Addr:              ":" + s.cfg.Port,
		ReadHeaderTimeout: 1 * time.Second,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	if s.cfg.MetricsEnabled {
		http.Handle("/metrics/", promhttp.Handler())
	}

	go func() {
		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.log.Fatal().Err(err).Msg("ListenAndServe error")
		}
	}()

	if s.cfg.MetricsEnabled {
		go s.startMetricsScrapper()
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.cfg.MetricsEnabled {
		s.stopScraping <- struct{}{}
	}
	return s.srv.Shutdown(ctx)
}
