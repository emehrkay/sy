package api

import (
	"net/http"
	"time"

	"github.com/emehrkay/sy/service"
)

type CreateHeartbeatRequest struct {
	SentAt time.Time `json:"sent_at"`
}

func (s *server) createHeartbeat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	deviceID := r.PathValue("deviceID")
	req, err := requestBody[CreateHeartbeatRequest](r)
	if err != nil {
		s.respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.monitorService.RegiserHeatbeat(ctx, service.Heartbeat{
		DeviceID: deviceID,
		SentAt:   req.SentAt,
	})

	s.respondJson(w, nil, http.StatusNoContent)
}
