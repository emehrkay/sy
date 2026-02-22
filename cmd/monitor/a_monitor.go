package monitor

import (
	"github.com/spf13/cobra"

	"github.com/emehrkay/sy/service"
	"github.com/emehrkay/sy/storage"
)

var (
	memoryStorage storage.Storage
	serviceLayer  service.Monitor
	RootCmd       = &cobra.Command{
		Use:   "monitor",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {
	memoryStorage = storage.NewMemory()
	serviceLayer = service.New(memoryStorage)
}
