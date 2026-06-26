package main

import (
	"context"
	"fmt"
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
	var mainCmd = &cobra.Command{
		Use:   "mysql-ps",
		Short: "MySQL Process List",
		Long:  `Show MySQL Process List, with refreshing it every N seconds`,
		Run: func(_ *cobra.Command, _ []string) {
			helpers.LoadConfig()
			ui.CreateUIGrid()
			ui.SetGlobalHandler(ui.KeyHandler)

			dsn := os.Getenv("MYSQL_DSN")

			dbStore, err := db.ConnectDB(
				os.Getenv("MYSQL_USER"),
				os.Getenv("MYSQL_PASSWORD"),
				dsn,
			)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			ui.IsRunningParam.Store(true)
			if ui.TimerSecParam <= 0 {
				ui.TimerSecParam = DefaultRefreshInterval
			}

			var wg sync.WaitGroup

			config := ui.WorkerConfig{
				TimerSec:   ui.TimerSecParam,
				ShowSystem: &ui.ShowSystem,
				IsRunning:  &ui.IsRunningParam,
				StatusBar:  ui.UIStatusBar,
				ListView:   ui.UIListView,
				SQLView:    ui.UISQLView,
				DSN:        dsn,
				Databases:  databases,
				App:        ui.UIApp,
				WG:         &wg,
			}

			ctx, cancel := context.WithCancel(context.Background())

			wg.Add(1)
			go ui.PSWorker(ctx, dbStore.GetProcessList, config)
			ui.Run()

			cancel()
			wg.Wait()
			dbStore.Close()
		},
	}

	mainCmd.Flags().Float32VarP(&ui.TimerSecParam, "interval", "i", DefaultRefreshInterval, "Refresh interval in seconds")
	mainCmd.Flags().BoolVarP(&ui.UseMouse, "mouse", "m", false, "Enable mouse interaction")
	mainCmd.Flags().StringArrayVarP(&databases, "database", "d", []string{}, "Databases list to filter by; example - -d b1 -d db2")

	if err := mainCmd.Execute(); err != nil {
		mainCmd.PrintErrln(err)
		os.Exit(1)
	}
}
