package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"analytics_collector/internal/api/storage"
	"analytics_collector/internal/lib/response"
	sl "analytics_collector/internal/logging"
)

const (
	AuthHeader = "X-Tantum-Authorization"
)

func HandleAnalytics(ctx context.Context, log *slog.Logger, jobs chan<- storage.UserActionInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.metrics.Analytics"
		log.With(slog.String("op", op))

		// check correct method of request
		if r.Method != "POST" {
			log.ErrorContext(ctx, "Method must be POST",
				slog.String("Current method", r.Method),
			)

			http.Error(w, "Method must be POST", http.StatusBadRequest)
			return
		}

		preparedDBInfo, err := prepareInfo(r)
		if err != nil {
			log.ErrorContext(ctx, "prepare info for saving failed", sl.Err(err))

			http.Error(w, fmt.Sprintf("Uncorrect request: %s", err.Error()), http.StatusBadRequest)
			return
		}

		jobs <- preparedDBInfo

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		jsonResp, err := json.Marshal(response.OK())
		if err != nil {
			log.ErrorContext(ctx, "bad json response", sl.Err(err))
			return
		}

		_, err = w.Write(jsonResp)
		if err != nil {
			log.ErrorContext(ctx, "bad json response write", sl.Err(err))
			return
		}

		return
	}
}

func prepareInfo(r *http.Request) (storage.UserActionInfo, error) {
	// check auth header
	userId, err := getAuthHeader(r.Header)
	if err != nil {
		return storage.UserActionInfo{}, fmt.Errorf("auth header error: %w", err)
	}

	// get body of request
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return storage.UserActionInfo{}, fmt.Errorf("reading body error: %w", err)
	}
	if !json.Valid(bodyBytes) {
		return storage.UserActionInfo{}, fmt.Errorf("bad json format body")
	}
	body := string(bodyBytes)

	// get headers from request
	jsonHeadersBytes, err := json.Marshal(r.Header)
	if err != nil {
		return storage.UserActionInfo{}, fmt.Errorf("reading headers error: %w", err)
	}
	jsonHeaders := string(jsonHeadersBytes)

	result := storage.UserActionInfo{
		Time:   time.Now(),
		UserID: userId,
		Data:   storage.RequestInfo{Body: body, Headers: jsonHeaders},
	}

	return result, nil
}

func getAuthHeader(headers http.Header) (string, error) {
	userId := headers.Get(AuthHeader)
	if userId == "" {
		return "", fmt.Errorf("auth header is empty")
	}

	return userId, nil
}
