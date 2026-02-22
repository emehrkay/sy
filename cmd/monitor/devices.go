package monitor

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/emehrkay/sy/service"
	"github.com/spf13/cobra"
)

func seedDevices(csvPath string) error {
	fh, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf(`unable to open file %s -- %w`, csvPath, err)
	}

	defer fh.Close()

	reader := csv.NewReader(fh)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf(`unable to read file %s -- %w`, csvPath, err)
	}

	deviceIDs := []string{}
	for i := 1; i < len(records); i++ {
		deviceIDs = append(deviceIDs, records[i][0])
	}

	return serviceLayer.SeedDevices(context.Background(), service.SeedDeviceRequest{
		DeviceIDs: deviceIDs,
	})
}

func init() {
	var (
		csvPath string
	)
	var seedInitialDevices = &cobra.Command{
		Use:   "seeddevices",
		Short: "adds initial devices to the db",
		Run: func(cmd *cobra.Command, args []string) {
			err := seedDevices(csvPath)
			if err != nil {
				log.Fatalf("error seeding devices -- %v", err)
				return
			}

			fmt.Println("Devices seeded")
		},
	}

	seedInitialDevices.Flags().StringVarP(&csvPath, "csv", "c", "../__/meta/devices.csv", "")
	RootCmd.AddCommand(seedInitialDevices)
}
