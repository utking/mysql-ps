package main

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/utking/mysql-ps/db"
	"github.com/utking/mysql-ps/helpers"
	"github.com/utking/mysql-ps/ui"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DefaultRefreshInterval = 2.0
)

var (
	databases []string
)

func main() {
	var mainCmd = &cobra.Command{
		Use:   "",
		Short: "MySQL Process List",
		Long:  `Show MySQL Process List, with refreshing it every N seconds`,
		Run: func(cmd *cobra.Command, args []string) {
			helpers.LoadConfig()
			ui.CreateUIGrid()
			ui.SetGlobalHandler(ui.KeyHandler)

			dbStore, err := db.ConnectDB(
				os.Getenv("MYSQL_USER"),
				os.Getenv("MYSQL_PASSWORD"),
				os.Getenv("MYSQL_DSN"),
			)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			ui.IsRunning = true
			if ui.TimerSec <= 0 {
				ui.TimerSec = DefaultRefreshInterval
			}

			config := ui.WorkerConfig{
				TimerSec:   ui.TimerSec,
				ShowSystem: ui.ShowSystem,
				IsRunning:  &ui.IsRunning,
				StatusBar:  ui.UIStatusBar,
				ListView:   ui.UIListView,
				SQLView:    ui.UISQLView,
				DSN:        os.Getenv("MYSQL_DSN"),
				Databases:  databases,
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go ui.PSWorker(ctx, dbStore.GetProcessList, nil, config)
			ui.Run()
		},
	}

	mainCmd.Flags().Float32VarP(&ui.TimerSec, "interval", "i", DefaultRefreshInterval, "Refresh interval in seconds")
	mainCmd.Flags().BoolVarP(&ui.UseMouse, "mouse", "m", false, "Enable mouse interaction")
	mainCmd.Flags().StringArrayVarP(&databases, "database", "d", []string{}, "Databases list to filter by; example - -d b1 -d db2")

	if err := mainCmd.Execute(); err != nil {
		mainCmd.Println(err)
		os.Exit(1)
	}
}
