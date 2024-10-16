package main

import (
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

			if err := db.ConnectDB(
				os.Getenv("MYSQL_USER"),
				os.Getenv("MYSQL_PASSWORD"),
				os.Getenv("MYSQL_DSN"),
			); err != nil {
				log.Println(err)
				os.Exit(1)
			}

			ui.IsRunning = true
			if ui.TimerSec <= 0 {
				ui.TimerSec = DefaultRefreshInterval
			}

			// expand slice of strings to slice of interfaces
			databaseList := make([]interface{}, len(databases))
			for i, v := range databases {
				databaseList[i] = v
			}

			go ui.PSWorker(db.GetProcessList, databaseList)
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
