package main

import (
	"log"
	"os"
	"strconv"

	"github.com/utking/mysql-ps/db"
	"github.com/utking/mysql-ps/helpers"
	"github.com/utking/mysql-ps/ui"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	helpers.LoadConfig()
	ui.CreateUIGrid()
	ui.SetGlobalHandler(ui.KeyHandler)

	if err := db.ConnectDB(os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DSN")); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ui.IsRunning = true
	if ui.TimerSec, _ = strconv.Atoi(os.Getenv("REFRESH_INTERVAL")); ui.TimerSec <= 0 {
		ui.TimerSec = 2
	}

	go ui.PSWorker(db.GetProcessList)
	ui.Run()
}
