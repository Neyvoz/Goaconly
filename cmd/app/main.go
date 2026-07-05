package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sitepulse/internal/infrastructure/netclient"
	"sitepulse/internal/repository/postgres"
	usercase "sitepulse/internal/usecase"
	"sitepulse/internal/worker"

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

	targetRepo := postgres.NewTargetRepo(db)
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

	time.AfterFunc(20*time.Second, func() {
		if err := sched.UpdateInterval(1, 1*time.Minute); err != nil {
			logger.Error("update interval failed", "error", err)
		}
		logger.Info("interval updated for target 1")
	})

	if err := uc.EnqueueAllTargets(ctx); err != nil {
		logger.Error("initial enqueue failed", "error", err)
	}

	<-ctx.Done()
	logger.Info("shutdown signal received, waiting for workers...")

	pool.Wait()
	logger.Info("all workers done, goodbye")
}
