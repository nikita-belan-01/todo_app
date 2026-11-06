package http

import (
	"context"
	"net/http"
	"sync"

	"github.com/nikita-belan-01/todo_app/internal/core/delivery/http/response"
	"github.com/nikita-belan-01/todo_app/pkg/logger"
	"go.uber.org/zap"
)

func (h HealthHTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke Health handler")

	responseHandler.JSONResponse(struct {
		Status string `json:"status"`
	}{
		Status: "healthy",
	}, http.StatusOK)
}

func (h HealthHTTPHandler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.Config.RequestTimeout)
	defer cancel()

	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke Ready handler")

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results = make(map[string]checkResult, len(h.Checks))
		ready   = true
	)

	for _, c := range h.Checks {
		wg.Add(1)

		go func() {
			defer wg.Done()

			res := checkResult{Status: upCheckStatus}
			if err := c.Pinger.Ping(ctx); err != nil {
				res = checkResult{
					Status: downCheckStatus,
					Error:  unavailableCheckError}

				log.Warn("readiness check failed", zap.String("check", c.Name), zap.Error(err))
			}

			mu.Lock()
			defer mu.Unlock()

			results[c.Name] = res
			if res.Status != upCheckStatus {
				ready = false
			}
		}()
	}

	wg.Wait()

	resp := readyResponse{Status: okReadyStatus, Checks: results}
	status := http.StatusOK

	if !ready {
		resp.Status = degradedReadyStatus
		status = http.StatusServiceUnavailable
	}

	responseHandler.JSONResponse(resp, status)
}
