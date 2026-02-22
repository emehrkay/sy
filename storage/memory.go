package storage

import (
	"sort"
	"sync"
)

func NewMemory() Storage {
	return &Memory{
		mu:         &sync.Mutex{},
		heartbeats: map[string]HeartbeatCollection{},
		stats:      map[string]StatsCollection{},
	}
}

type Memory struct {
	mu         *sync.Mutex
	heartbeats map[string]HeartbeatCollection
	stats      map[string]StatsCollection
}

func (m *Memory) AddDevice(deviceIDs ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, dID := range deviceIDs {
		m.heartbeats[dID] = HeartbeatCollection{}
		m.stats[dID] = StatsCollection{}
	}

	return nil
}

func (m *Memory) GetDevice(deviceID string) (*HeartbeatCollection, *StatsCollection, error) {
	var hb HeartbeatCollection
	var stats StatsCollection
	var ok bool

	if hb, ok = m.heartbeats[deviceID]; !ok {
		return nil, nil, ErrDeviceNotFound
	}

	if stats, ok = m.stats[deviceID]; !ok {
		return nil, nil, ErrDeviceNotFound
	}

	// mimic sorting that would be apart of a db query
	sort.Slice(hb, func(i, j int) bool {
		return hb[i].SentAt.Before(hb[j].SentAt)
	})
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].SentAt.Before(stats[j].SentAt)
	})

	return &hb, &stats, nil
}

func (m *Memory) AddHeartBeat(heartbeat Heartbeat) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, _, err := m.GetDevice(heartbeat.DeviceID)
	if err != nil {
		return err
	}

	m.heartbeats[heartbeat.DeviceID] = append(m.heartbeats[heartbeat.DeviceID], heartbeat)

	return nil
}

func (m *Memory) AddStats(stats Stats) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, _, err := m.GetDevice(stats.DeviceID)
	if err != nil {
		return err
	}

	m.stats[stats.DeviceID] = append(m.stats[stats.DeviceID], stats)

	return nil
}

func (m *Memory) GetStats(deviceID string) (StatsCollection, error) {
	_, stats, err := m.GetDevice(deviceID)
	if err != nil {
		return nil, err
	}

	return *stats, nil
}
