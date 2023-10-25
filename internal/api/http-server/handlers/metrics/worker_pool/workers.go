package worker_pool

import (
	"analytics_collector/internal/api/storage"
	sl "analytics_collector/internal/logging"
	"context"
	"fmt"
	"log/slog"
	"sync"
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
		case userInfo := <-jobs:
			err := saver.Save(ctx, userInfo)
			if err != nil {
				log.ErrorContext(ctx,
					fmt.Sprintf("worker %d failed task", workerNumber),
					slog.String("Job", fmt.Sprintf("%+v", userInfo)),
					sl.Err(err),
				)

				continue
			} else {
				log.InfoContext(ctx,
					fmt.Sprintf("worker %d done task", workerNumber),
					slog.String("Job", fmt.Sprintf("%+v", userInfo)),
				)
			}
		}
	}
}
