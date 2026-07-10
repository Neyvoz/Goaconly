package httpserver

import (
	"net/http"
	"time"
)

const (
	// ReadHeaderTimeout защищает от Slowloris-атаки: клиент открывает
	// соединение и присылает заголовки по одному байту, держа воркер
	// сервера занятым бесконечно. 5s достаточно для нормального клиента
	// даже на плохом мобильном соединении.
	defaultReadHeaderTimeout = 5 * time.Second
	// ReadTimeout — общий лимit на чтение всего запроса (заголовки + тело).
	defaultReadTimeout = 10 * time.Second
	// WriteTimeout — лимит на запись ответа клиенту.
	// ВАЖНО: этот таймаут считается от момента, когда сервер принял
	// соединение, а не от начала записи — учитывай это при долгих хендлерах.
	defaultWriteTimeout = 15 * time.Second
	// IdleTimeout — сколько keep-alive соединение может простаивать
	// между запросами до закрытия. 60s — стандартная практика для API,
	// балансирует между переиспользованием TCP-соединений (экономия на
	// хендшейках) и утечкой файловых дескрипторов на простаивающих клиентах.
	defaultIdleTimeout = 60 * time.Second
)

type Config struct {
	Addr string
}

type Server struct {
	httpServer *http.Server
}

func New(cfg Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              cfg.Addr,
			Handler:           handler,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
			ReadTimeout:       defaultReadTimeout,
			WriteTimeout:      defaultWriteTimeout,
			IdleTimeout:       defaultIdleTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
