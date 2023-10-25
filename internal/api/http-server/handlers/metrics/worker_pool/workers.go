package worker_pool

import (
	"analytics_collector/internal/api/storage"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Saver interface {
	Save(ctx context.Context, info storage.UserActionInfo) error
}

// New create channel for reading task for workers pool. Create `count` number workers
// and return open channel for writing.
func New(ctx context.Context, log *slog.Logger, count int, saver Saver) chan storage.UserActionInfo {
	jobsChannel := make(chan storage.UserActionInfo, count)

	var wg sync.WaitGroup

	for i := 1; i <= count; i++ {
		i := i

		wg.Add(1)
		go func() {
			defer wg.Done()
			startWorker(ctx, log, i, jobsChannel, saver)
		}()
	}

	go func() {
		wg.Wait()
		close(jobsChannel)
		log.Info("jobs channel closed")
	}()

	return jobsChannel
}

func startWorker(ctx context.Context, log *slog.Logger, workerNumber int, jobs <-chan storage.UserActionInfo, saver Saver) {
	for {
		select {
		case <-ctx.Done():
			log.Info(fmt.Sprintf("worker %d is ended", workerNumber))
			return
		case job := <-jobs:
			// TODO: implements workers
			time.Sleep(3 * time.Second)
			log.Info(
				fmt.Sprintf("worker %d done task", workerNumber),
				slog.String("Job", fmt.Sprintf("%+v", job)),
			)
		}
	}
}
