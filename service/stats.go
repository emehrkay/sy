package service

import (
	"context"
	"fmt"
	"time"

	"github.com/emehrkay/sy/storage"
)

type GetStatsRequest struct {
	DeviceID string
}

type Stats struct {
	DeviceID   string
	SentAt     time.Time
	UploadTime uint64
}

type StatsCollection []Stats

func (sc *StatsCollection) FromStore(stats storage.StatsCollection) {
	for _, st := range stats {
		*sc = append(*sc, Stats{
			DeviceID:   st.DeviceID,
			SentAt:     st.SentAt,
			UploadTime: st.UploadTime,
		})
	}
}

type StatsResponse struct {
	Uptime                         float64
	AverageUploadTime              float64
	AverageUploadTimeHumanReadable string
	heartbeats                     storage.HeartbeatCollection
	stats                          storage.StatsCollection
}

func (st *StatsResponse) caculateTotals() {
	st.Uptime = float64(len(st.heartbeats)) / st.heartbeats.MinuteRange() * 100
	st.AverageUploadTime = st.stats.AverageUploadTime()
	st.AverageUploadTimeHumanReadable = st.stats.AverageUploadTimeHumanReadable()
}

func (m *monitor) AddStats(ctx context.Context, stats Stats) error {
	err := m.store.AddStats(storage.Stats{
		DeviceID:   stats.DeviceID,
		SentAt:     stats.SentAt,
		UploadTime: stats.UploadTime,
	})
	if err != nil {
		return fmt.Errorf(`unable to add stats -- %w`, err)
	}

	return nil
}

func (m *monitor) GetStats(ctx context.Context, request GetStatsRequest) (*StatsResponse, error) {
	heartbeats, stats, err := m.store.GetDevice(request.DeviceID)
	if err != nil {
		return nil, fmt.Errorf(`unable to get stats for device: %v -- %w`, request.DeviceID, err)
	}

	resp := &StatsResponse{
		heartbeats: *heartbeats,
		stats:      *stats,
	}
	resp.caculateTotals()

	return resp, nil
}
