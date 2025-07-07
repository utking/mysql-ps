package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/utking/mysql-ps/helpers"
)

var (
	TimerSec    float32
	IsRunning   bool
	ShowSystem  bool
	UseMouse    bool
	ListLengh   int
	status      string
	listFilters []string
	Databases   []string
)

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
	listFn func([]string, []interface{}) ([]helpers.ProcessItem, error),
	databases []interface{},
) {
	for range time.Tick(time.Millisecond * time.Duration(1000*TimerSec)) {
		if !ShowSystem {
			listFilters = []string{"DB != 'sys'"}
		} else {
			listFilters = []string{}
		}

		if !IsRunning {
			status = "Paused"

			UIStatusBar.SetBorderColor(tcell.ColorYellow)
			UIApp.QueueUpdateDraw(func() {
				UpdateStatusBar(status, ListLengh)
			})

			continue
		}

		var (
			err       error
			itemsList []helpers.ProcessItem
		)

		status = "Running"

		UIStatusBar.SetBorderColor(tcell.ColorWhite)
		UIListView.Clear()

		if itemsList, err = listFn(listFilters, databases); err != nil {
			UISQLView.SetText(err.Error())

			IsRunning = false

			continue
		}

		ListLengh = len(itemsList)

		for i := range itemsList {
			if strings.Contains(
				itemsList[i].Info.String,
				"INFORMATION_SCHEMA.PROCESSLIST",
			) {
				ListLengh--
			}
		}

		UpdateStatusBar(status, ListLengh)
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

				UIListView.AddItem(lineName, itemsList[i].Info.String, 0, nil)
			}
		})
	}
}
