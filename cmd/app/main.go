package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpserver "goaconly/internal/delivery/http"
	"goaconly/internal/delivery/http/handler"
	"goaconly/internal/infrastructure/netclient"
	"goaconly/internal/repository/postgres"
	usercase "goaconly/internal/usecase"
	"goaconly/internal/worker"

	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Сборка HTTP-слоя: repository -> usecase -> handler -> router
	targetRepo := postgres.NewTargetRepo(db)
	targetUsecase := usercase.NewTargetUsecase(targetRepo)
	targetHandler := handler.NewTargetHandler(targetUsecase)

	router := httpserver.NewRouter(httpserver.Dependencies{
		TargetHandler: targetHandler,
	})
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}

	srv := httpserver.New(httpserver.Config{Addr: httpAddr}, router)

	// Сборка конвейера мониторинга: checker + worker pool + scheduler
	checkResultRepo := postgres.NewCheckResultRepo(db)

	checker := netclient.NewChecker(logger)

	uc := usercase.NewMonitoringUseCase(targetRepo, checkResultRepo, checker, nil, logger)

	pool := worker.NewPool(
		10, 50,
		uc.HandleCheckJob,
		logger,
	)
	uc.SetPool(pool)

	sched := usercase.NewScheduler(pool, checkResultRepo)

	// Загрузка активных целей из БД и постановка на расписание
	targets, err := targetRepo.GetAllActive(ctx)
	if err != nil {
		logger.Error("failed to load active targets", "error", err)
	}

	go sched.Run(ctx)

	for _, t := range targets {
		if err := sched.AddTarget(ctx, t); err != nil {
			logger.Warn("failed to add target to scheduler",
				"target_id", t.ID,
				"error", err,
			)
		}
	}

	pool.Start(ctx)

	// Тестовое обновление интервала (временно, для проверки динамического планировщика)
	time.AfterFunc(20*time.Second, func() {
		if err := sched.UpdateInterval(1, 1*time.Minute); err != nil {
			logger.Error("update interval failed", "error", err)
		}
		logger.Info("interval updated for target 1")
	})

	if err := uc.EnqueueAllTargets(ctx); err != nil {
		logger.Error("initial enqueue failed", "error", err)
	}

	// Запуск сервера
	go func() {
		logger.Info("starting HTTP server", "addr", httpAddr)
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", "error", err)
		}
	}()

	// Ожидание сигнала завершения
	<-ctx.Done()
	logger.Info("shutdown signal received, waiting for workers...")

	pool.Wait()
	logger.Info("all workers done, goodbye")
}
