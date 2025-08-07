package router

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"awesomeProject/internal/router/handler"
	"go.uber.org/zap"
)

type Router struct {
	mux         *http.ServeMux
	log         *zap.Logger
	coinHandler *handler.Handler
	server      *http.Server
}

func NewRouter(coinHandler *handler.Handler, log *zap.Logger) *Router {
	return &Router{
		mux:         http.NewServeMux(),
		log:         log.Named("request"),
		coinHandler: coinHandler,
	}
}

func (r *Router) RunRouter(addr string) error {
	// Настройка обработчиков
	r.mux.HandleFunc("/currency/add", r.coinHandler.AddCoin)
	r.mux.HandleFunc("/currency/get", r.coinHandler.GetCoin)
	r.mux.HandleFunc("/currency/remove", r.coinHandler.DeleteCoin)

	r.server = &http.Server{
		Addr:    addr,
		Handler: r.loggingMiddleware(r.mux),
	}

	serverErr := make(chan error, 1)

	go func() {
		r.log.Info("Starting server", zap.String("addr", addr))
		if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
		close(serverErr)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		r.log.Info("Received signal", zap.String("signal", sig.String()))
	case err := <-serverErr:
		r.log.Error("Server error", zap.Error(err))
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r.log.Info("Shutting down server...")
	if err := r.server.Shutdown(ctx); err != nil {
		r.log.Error("Forced shutdown", zap.Error(err))
		return err
	}

	r.log.Info("Server stopped gracefully")
	return nil
}

func (r *Router) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestLog := r.log.With(
			zap.String("method", req.Method),
			zap.String("path", req.URL.Path),
			zap.String("remote_addr", req.RemoteAddr),
		)

		requestLog.Info("Request started")
		ctx := context.WithValue(req.Context(), "logger", requestLog)
		next.ServeHTTP(w, req.WithContext(ctx))
		requestLog.Info("Request completed")
	})
}
