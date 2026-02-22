package service

import (
	"context"
	"fmt"
	"time"

	"github.com/emehrkay/sy/storage"
)

type Heartbeat struct {
	DeviceID string
	SentAt   time.Time
}

func (m *monitor) RegiserHeatbeat(ctx context.Context, heartbeat Heartbeat) error {
	err := m.store.AddHeartBeat(storage.Heartbeat{
		DeviceID: heartbeat.DeviceID,
		SentAt:   heartbeat.SentAt,
	})
	if err != nil {
		return fmt.Errorf(`unable to add heartbeat -- %w`, err)
	}

	return nil
}
