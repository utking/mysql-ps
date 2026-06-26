package main

import (
	"context"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/utking/mysql-ps/db"
	"github.com/utking/mysql-ps/helpers"
	"github.com/utking/mysql-ps/ui"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DefaultRefreshInterval = float32(2.0)
)

var (
	databases []string
)

func main() {
	uiComponents := ui.NewUI()

	var mainCmd = &cobra.Command{
		Use:   "mysql-ps",
		Short: "MySQL Process List",
		Long:  `Show MySQL Process List, with refreshing it every N seconds`,
		RunE: func(_ *cobra.Command, _ []string) error {
			helpers.LoadConfig()
			uiComponents.SetGlobalHandler()

			dsn := os.Getenv("MYSQL_DSN")
			user := os.Getenv("MYSQL_USER")
			password := os.Getenv("MYSQL_PASSWORD")

			dbStore, err := db.ConnectDB(user, password, dsn)
			if err != nil {
				return err
			}

			uiComponents.IsRunning.Store(true)
			if uiComponents.TimerSec <= 0 {
				uiComponents.TimerSec = DefaultRefreshInterval
			}

			var wg sync.WaitGroup

			config := ui.WorkerConfig{
				UI:        uiComponents,
				DSN:       dsn,
				Databases: databases,
				WG:        &wg,
			}

			ctx, cancel := context.WithCancel(context.Background())

			wg.Add(1)
			go ui.PSWorker(ctx, dbStore.GetProcessList, config)
			uiComponents.Run()

			cancel()
			wg.Wait()
			dbStore.Close()

			return nil
		},
	}

	mainCmd.Flags().Float32VarP(&uiComponents.TimerSec, "interval", "i", DefaultRefreshInterval, "Refresh interval in seconds")
	mainCmd.Flags().BoolVarP(&uiComponents.UseMouse, "mouse", "m", false, "Enable mouse interaction")
	mainCmd.Flags().StringArrayVarP(&databases, "database", "d", []string{}, "Databases list to filter by; example - -d b1 -d db2")

	if err := mainCmd.Execute(); err != nil {
		mainCmd.PrintErrln(err)
	}
}
