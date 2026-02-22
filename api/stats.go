package api

import (
	"net/http"
	"time"

	"github.com/emehrkay/sy/service"
)

type CreateStatsRequest struct {
	SentAt     time.Time `json:"sent_at"`
	UploadTime uint64    `json:"upload_time"`
}

type GetStatsResponse struct {
	Uptime            float64 `json:"uptime"`
	AverageUploadTime string  `json:"avg_upload_time"`
}

func (s *server) createStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	deviceID := r.PathValue("deviceID")
	req, err := requestBody[CreateStatsRequest](r)
	if err != nil {
		s.respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.monitorService.AddStats(ctx, service.Stats{
		DeviceID:   deviceID,
		SentAt:     req.SentAt,
		UploadTime: req.UploadTime,
	})
	s.respondJson(w, nil, http.StatusNoContent)
}

func (s *server) getStats(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("deviceID")
	stats, err := s.monitorService.GetStats(r.Context(), service.GetStatsRequest{
		DeviceID: deviceID,
	})
	if err != nil {
		s.respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := GetStatsResponse{
		Uptime:            stats.Uptime,
		AverageUploadTime: stats.AverageUploadTimeHumanReadable,
	}

	s.respondJson(w, resp, http.StatusOK)
}
