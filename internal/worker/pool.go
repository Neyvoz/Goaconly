package worker

import (
	"context"
	"log/slog"
	"sync"
)

// Job — обобщённый тип задачи.
type Job[T any] struct {
	Payload T
}

// Pool — пул воркеров.
type Pool[T any] struct {
	workerCount int
	jobs        chan Job[T]
	handler     func(ctx context.Context, job Job[T])
	wg          sync.WaitGroup
	logger      *slog.Logger
}

// NewPool создаёт пул. bufferSize — размер буфера канала jobs.
// Правило: bufferSize >= workerCount * 2, чтобы воркеры не голодали.
func NewPool[T any](
	workerCount int,
	bufferSize int,
	handler func(ctx context.Context, job Job[T]),
	logger *slog.Logger,
) *Pool[T] {
	return &Pool[T]{
		workerCount: workerCount,
		jobs:        make(chan Job[T], bufferSize),
		handler:     handler,
		logger:      logger,
	}
}

// Start запускает workerCount горутин и блокируется до отмены ctx.
// После отмены ctx дожидается завершения всех горутин.
func (p *Pool[T]) Start(ctx context.Context) {
	for i := range p.workerCount {
		p.wg.Add(1)
		go p.runWorker(ctx, i)
	}
}

// Submit отправляет задачу в пул. Неблокирующий вариант:
// если канал переполнен — логируем и дропаем (backpressure).
func (p *Pool[T]) Submit(job Job[T]) bool {
	select {
	case p.jobs <- job:
		return true
	default:
		p.logger.Warn("worker pool: job queue is full, dropping job")
		return false
	}
}

// Wait ожидает завершения всех воркеров. Вызывай после отмены ctx.
func (p *Pool[T]) Wait() {
	close(p.jobs)
	p.wg.Wait()
}

func (p *Pool[T]) runWorker(ctx context.Context, id int) {
	defer p.wg.Done()
	p.logger.Info("worker started", "id", id)
	for {
		select {
		case job, ok := <-p.jobs:
			if !ok {
				p.logger.Info("worker stopped", "id", id)
				return
			}
			p.handler(ctx, job)
		case <-ctx.Done():
			p.logger.Info("worker context cancelled, draining...", "id", id)
			for job := range p.jobs {
				p.handler(ctx, job)
			}
			return
		}
	}
}
