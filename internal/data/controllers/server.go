package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/redselig/currencier/internal/domain/usecase"
	"github.com/redselig/currencier/internal/util"
)

const (
	ErrID = "must be id in query"
)

type HTTPServer struct {
	logger     usecase.Logger
	server     *http.Server
	currencier usecase.Currencier
}

func NewHttpServer(addr string, logger usecase.Logger, currencier usecase.Currencier) *HTTPServer {
	server := &http.Server{Addr: addr}
	return &HTTPServer{
		server:     server,
		logger:     logger,
		currencier: currencier,
	}
}

func (s *HTTPServer) Serve() error {
	s.logger.Log(context.Background(), "starting http server on address [%v]", s.server.Addr)

	router := mux.NewRouter()

	router.HandleFunc("/", s.getCurrencies).Methods(http.MethodGet)
	router.HandleFunc("/currencies", s.getCurrencies).Methods(http.MethodGet)
	router.HandleFunc("/lazycurrencies", s.getLazyCurrencies).Methods(http.MethodGet)

	router.HandleFunc("/currency/{id}", s.getCurrency).Methods(http.MethodGet) //todo: should be /currencies/{id}

	handler := s.accessLogMiddleware(router)
	handler = s.panicMiddleware(handler)
	s.server.Handler = handler

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return errors.Wrapf(err, "can't start listen address [%v]", s.server.Addr)
	}
	return nil
}

func (s *HTTPServer) getCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok || id == "" {
		s.httpError(r.Context(), w, ErrID, http.StatusBadRequest)
		return
	}
	c, err := s.currencier.GetCurrencyBuID(r.Context(), id)
	if err != nil {
		s.httpError(r.Context(), w, err.Error(), http.StatusBadRequest)
		return
	}
	s.httpAnswer(w, c, http.StatusOK)
}

func (s *HTTPServer) getCurrencies(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()

	limit, ok := vars["limit"]
	if !ok || len(limit) != 1 {
		limit = []string{"10"}
	}
	offset, ok := vars["offset"]
	if !ok || len(offset) != 1 {
		offset = []string{"0"}
	}
	iLimit, err := strconv.Atoi(limit[0])
	if err != nil {
		s.httpError(r.Context(), w, err.Error(), http.StatusBadRequest)
		return
	}
	iOffset, err := strconv.Atoi(offset[0])
	if err != nil {
		s.httpError(r.Context(), w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := s.currencier.GetCurrenciesPage(r.Context(), iLimit, iOffset)
	if err != nil {
		s.httpError(r.Context(), w, err.Error(), http.StatusBadRequest)
		return
	}
	s.httpAnswer(w, c, http.StatusOK)
}

func (s *HTTPServer) getLazyCurrencies(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()

	limit, ok := vars["limit"]
	if !ok || len(limit) != 1 {
		limit = []string{"10"}
	}
	lastID, ok := vars["lastid"]
	if !ok || len(lastID) != 1 {
		lastID = []string{""}
	}
	iLimit, err := strconv.Atoi(limit[0])
	if err != nil {
		s.httpError(r.Context(), w, err.Error(), http.StatusBadRequest)
	}

	c, err := s.currencier.GetCurrenciesLazy(r.Context(), iLimit, lastID[0])
	if err != nil {
		s.httpError(r.Context(), w, err.Error(), http.StatusBadRequest)
		return
	}
	s.httpAnswer(w, c, http.StatusOK)
}

func (s *HTTPServer) StopServe() {
	ctx := context.Background()
	s.logger.Log(ctx, "stopping http server")
	defer s.logger.Log(ctx, "http server stopped")
	if s.server == nil {
		s.logger.Log(ctx, "http server is nil")
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Log(ctx, "can't stop http server with error: %v", err)
	}
}

func (s *HTTPServer) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := util.SetRequestID(r.Context())

		next.ServeHTTP(w, r.WithContext(ctx))

		latency := time.Since(start)
		s.logRequest(ctx, r.RemoteAddr, start.Format(util.LayoutISO), r.Method, r.URL.Path, latency)
	})
}

func (s *HTTPServer) panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.httpError(r.Context(), w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *HTTPServer) logRequest(ctx context.Context, remoteAddr, start, method, path string, latency time.Duration) {
	s.logger.Log(ctx, "%s [%s] %s %s [%s]", remoteAddr, start, method, path, latency)
}

func (s *HTTPServer) httpError(ctx context.Context, w http.ResponseWriter, error string, code int) {
	s.logger.Log(ctx, error)
	http.Error(w, error, code)
}

func (s *HTTPServer) httpAnswer(w http.ResponseWriter, msg interface{}, code int) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		code = http.StatusInternalServerError
	}
	w.WriteHeader(code)
	w.Write(jmsg) //nolint:errcheck
}
