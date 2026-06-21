package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/utking/mysql-ps/helpers"
)

var (
	TimerSec   float32
	IsRunning  bool
	ShowSystem bool
	UseMouse   bool
	ListLengh  int
)

type WorkerConfig struct {
	TimerSec   float32
	ShowSystem bool
	IsRunning  *bool
	StatusBar  *tview.TextView
	ListView   *tview.List
	SQLView    *tview.TextView
	DSN        string
	Databases  []string
}

func Run() {
	UIListView.SetSelectedFunc(OpenSQLQuery)

	if err := UIApp.
		SetRoot(UIFlex, true).
		EnableMouse(UseMouse).
		Run(); err != nil {
		panic(err)
	}
}

func PSWorker(
	ctx context.Context,
	listFn func([]string, []interface{}) ([]helpers.ProcessItem, error),
	databases []interface{},
	config WorkerConfig,
) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(1000*config.TimerSec))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var listFilters []string
			if !config.ShowSystem {
				listFilters = []string{"DB != 'sys'"}
			} else {
				listFilters = []string{}
			}

			if *config.IsRunning == false {
				status := "Paused"
				listLen := 0 // Initialize local listLen
				UIApp.QueueUpdateDraw(func() {
					config.StatusBar.SetBorderColor(tcell.ColorYellow)
					UpdateStatusBar(
						config.StatusBar,
						status,
						listLen,
						config.TimerSec,
						config.ShowSystem, config.DSN, getMemUsage())
				})
				continue
			}

			var (
				err       error
				itemsList []helpers.ProcessItem
			)

			status := "Running"

			UIApp.QueueUpdateDraw(func() {
				config.StatusBar.SetBorderColor(tcell.ColorWhite)
				config.ListView.Clear()
			})

			// Convert databases from config to []interface{} for listFn
			dbInterfaces := make([]interface{}, len(config.Databases))
			for i, v := range config.Databases {
				dbInterfaces[i] = v
			}

			if itemsList, err = listFn(listFilters, dbInterfaces); err != nil {
				config.SQLView.SetText(err.Error())
				*config.IsRunning = false
				continue
			}

			listLen := len(itemsList)

			for i := range itemsList {
				if strings.Contains(
					itemsList[i].Info.String,
					"INFORMATION_SCHEMA.PROCESSLIST",
				) {
					listLen--
				}
			}

			UpdateStatusBar(
				config.StatusBar,
				status,
				listLen,
				config.TimerSec,
				config.ShowSystem, config.DSN, getMemUsage())
			UIApp.QueueUpdateDraw(func() {
				for i := range itemsList {
					if strings.Contains(
						itemsList[i].Info.String,
						"INFORMATION_SCHEMA.PROCESSLIST",
					) {
						continue
					}
					lineName := fmt.Sprintf("%d: %s (%ds) from %s@%s - %s",
						itemsList[i].ID,
						itemsList[i].DB.String,
						itemsList[i].Time,
						itemsList[i].User,
						helpers.HostDropPort(itemsList[i].Host),
						itemsList[i].State.String)

					config.ListView.AddItem(lineName, itemsList[i].Info.String, 0, nil)
				}
			})
		}
	}
}
