package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/utking/mysql-ps/helpers"
)

var (
	TimerSec  int
	IsRunning bool
	ListLengh int
	status    string
)

func Run() {
	UIListView.SetSelectedFunc(OpenSQLQuery)

	if err := UIApp.SetRoot(UIFlex, true).Run(); err != nil {
		panic(err)
	}
}

func PSWorker(listFn func() ([]helpers.ProcessItem, error)) {
	for range time.Tick(time.Second * time.Duration(TimerSec)) {
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

		if itemsList, err = listFn(); err != nil {
			UISQLView.SetText(err.Error())

			IsRunning = false

			continue
		}

		ListLengh = len(itemsList)

		UpdateStatusBar(status, ListLengh)
		UIApp.QueueUpdateDraw(func() {
			for i := range itemsList {
				lineName := fmt.Sprintf("%s (%ds) from %s@%s",
					itemsList[i].DB.String,
					itemsList[i].Time,
					itemsList[i].User,
					helpers.HostDropPort(itemsList[i].Host))

				UIListView.AddItem(lineName, itemsList[i].Info.String, 0, nil)
			}
		})
	}
}
