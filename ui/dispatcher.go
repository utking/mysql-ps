package ui

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/utking/mysql-ps/helpers"
)

var (
	TimerSecParam  float32
	IsRunningParam atomic.Bool
	ShowSystem     atomic.Bool
	UseMouse       bool
)

var updateTriggerChan chan struct{}
var lastManualUpdate time.Time

func init() {
	updateTriggerChan = make(chan struct{}, 1)
}

type WorkerConfig struct {
	TimerSec       float32
	ShowSystem     *atomic.Bool
	IsRunning      *atomic.Bool
	StatusBar      *tview.TextView
	ListView       *tview.List
	SQLView        *tview.TextView
	DSN            string
	Databases      []string
	App            *tview.Application
	OptionalUpdate func(func()) // Changed from Update to OptionalUpdate
}

func (c *WorkerConfig) Update(fn func()) {
	if c.OptionalUpdate != nil {
		c.OptionalUpdate(fn)
	} else if c.App != nil {
		c.App.QueueUpdateDraw(fn)
	} else {
		fn()
	}
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
	listFn func([]string, []any) ([]helpers.ProcessItem, error),
	config WorkerConfig,
) {
	ticker := time.NewTicker(time.Duration(float64(time.Second) * float64(config.TimerSec)))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			performUpdate(&config, listFn)
		case <-updateTriggerChan:
			performUpdate(&config, listFn)
		}
	}
}

func performUpdate(
	config *WorkerConfig,
	listFn func([]string, []any) ([]helpers.ProcessItem, error),
) {
	var listFilters []string
	if !ShowSystem.Load() {
		listFilters = []string{"DB != 'sys'"}
	} else {
		listFilters = []string{}
	}

	if config.IsRunning.Load() == false {
		status := "Paused"
		listLen := 0

		config.Update(func() {
			config.StatusBar.SetBorderColor(tcell.ColorYellow)
			UpdateStatusBar(
				config.StatusBar,
				status,
				listLen,
				config.TimerSec,
				ShowSystem.Load(),
				config.DSN,
				getMemUsage())
		})
		return
	}

	var (
		err       error
		itemsList []helpers.ProcessItem
	)

	dbInterfaces := make([]any, len(config.Databases))
	for i, v := range config.Databases {
		dbInterfaces[i] = v
	}

	if itemsList, err = listFn(listFilters, dbInterfaces); err != nil {
		config.SQLView.SetText(err.Error())
		config.IsRunning.Store(false)
		return
	}

	status := "Running"
	listLen := len(itemsList)

	for i := range itemsList {
		if strings.Contains(
			itemsList[i].Info.String,
			"INFORMATION_SCHEMA.PROCESSLIST",
		) {
			listLen--
		}
	}

	config.Update(func() {
		config.StatusBar.SetBorderColor(tcell.ColorWhite)
		config.ListView.Clear()
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
		UpdateStatusBar(
			config.StatusBar,
			status,
			listLen,
			config.TimerSec,
			config.ShowSystem.Load(),
			config.DSN,
			getMemUsage())
	})
}
