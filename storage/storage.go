package storage

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrDeviceNotFound = errors.New("device not found")
)

type Storage interface {
	AddDevice(deviceIDs ...string) error
	GetDevice(deviceID string) (*HeartbeatCollection, *StatsCollection, error)
	AddHeartBeat(heartbeat Heartbeat) error
	AddStats(stats Stats) error
	GetStats(deviceID string) (StatsCollection, error)
}

type Heartbeat struct {
	DeviceID string
	SentAt   time.Time
}

type HeartbeatCollection []Heartbeat

func (hbc HeartbeatCollection) MinuteRange() float64 {
	if len(hbc) < 2 {
		return 0
	}

	return hbc[len(hbc)-1].SentAt.Sub(hbc[0].SentAt).Minutes()
}

type Stats struct {
	DeviceID   string
	SentAt     time.Time
	UploadTime uint64
}

type StatsCollection []Stats

func (sc StatsCollection) AverageUploadTime() float64 {
	if len(sc) == 0 {
		return 0
	}

	var totalUploads float64
	for _, up := range sc {
		totalUploads += float64(up.UploadTime)
	}

	return totalUploads / float64(len(sc))
}

func (sc StatsCollection) AverageUploadTimeHumanReadable() string {
	avgTime := sc.AverageUploadTime()
	seconds := int64(avgTime) / int64(time.Second)
	rem := int64(avgTime) % int64(time.Second)
	avg := time.Unix(seconds, rem)
	resp := strings.Builder{}

	// hours are unnecessary
	// if avg.Hour() > 0 {
	// 	fmt.Fprintf(&resp, "%dh", avg.Hour()))
	// }

	fmt.Fprintf(&resp, "%dm", avg.Minute())
	fmt.Fprintf(&resp, "%d.%ds", avg.Second(), avg.Nanosecond())

	return resp.String()
}
