package service

import (
	"context"

	"github.com/emehrkay/sy/storage"
)

type Monitor interface {
	SeedDevices(ctx context.Context, request SeedDeviceRequest) error
	RegiserHeatbeat(ctx context.Context, heartbeat Heartbeat) error
	AddStats(ctx context.Context, stats Stats) error
	GetStats(ctx context.Context, request GetStatsRequest) (*StatsResponse, error)
}

func New(store storage.Storage) Monitor {
	return &monitor{
		store: store,
	}
}

type monitor struct {
	store storage.Storage
}
