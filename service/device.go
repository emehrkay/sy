package service

import "context"

type SeedDeviceRequest struct {
	DeviceIDs []string
}

func (m *monitor) SeedDevices(ctx context.Context, request SeedDeviceRequest) error {
	return m.store.AddDevice(request.DeviceIDs...)
}
