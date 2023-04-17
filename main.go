package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/utking/mysql-ps/db"
	"github.com/utking/mysql-ps/ui"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file", err)
	}

	ui.CreateUIGrid()
	ui.SetGlobalHandler(ui.KeyHandler)

	if err := db.ConnectDB(os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DSN")); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ui.IsRunnin = true
	if ui.TimerSec, _ = strconv.Atoi(os.Getenv("REFRESH_INTERVAL")); ui.TimerSec <= 0 {
		ui.TimerSec = 1
	}

	go ui.PSWorker(db.GetProcessList)
	ui.Run()
}
