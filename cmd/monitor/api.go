package monitor

import (
	"log"
	"net/http"

	"github.com/emehrkay/sy/api"
	"github.com/spf13/cobra"
)

func init() {
	var (
		port    string
		csvPath string
	)

	var startAPIServer = &cobra.Command{
		Use:   "apistart",
		Short: "starts the api server",
		Run: func(cmd *cobra.Command, args []string) {
			// seed the devices before starting the server
			err := seedDevices(csvPath)
			if err != nil {
				log.Fatalf(`unable to seed devices: %v`, err)
			}

			router := http.NewServeMux()
			server := api.New(port, serviceLayer, router)
			log.Fatal(server.Run())
		},
	}

	startAPIServer.Flags().StringVarP(&csvPath, "csv", "c", "../__/meta/devices.csv", "")
	startAPIServer.Flags().StringVarP(&port, "port", "p", ":6733", "")
	RootCmd.AddCommand(startAPIServer)
}
