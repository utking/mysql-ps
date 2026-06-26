package ui

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/utking/mysql-ps/helpers"
)

type WorkerConfig struct {
	UI             *UIComponents
	DSN            string
	Databases      []string
	WG             *sync.WaitGroup
	OptionalUpdate func(func())
}

func (c *WorkerConfig) Update(fn func()) {
	if c.OptionalUpdate != nil {
		c.OptionalUpdate(fn)
	} else if c.UI != nil && c.UI.App != nil {
		c.UI.App.QueueUpdateDraw(fn)
	} else {
		fn()
	}
}

func (c *UIComponents) Run() {
	c.ListView.SetSelectedFunc(c.OpenSQLQuery)

	if err := c.App.
		SetRoot(c.Flex, true).
		EnableMouse(c.UseMouse).
		Run(); err != nil {
		panic(err)
	}
}

func PSWorker(
	ctx context.Context,
	listFn func([]string, []any) ([]helpers.ProcessItem, error),
	config WorkerConfig,
) {
	if config.WG != nil {
		defer config.WG.Done()
	}

	ticker := time.NewTicker(time.Duration(float64(time.Second) * float64(config.UI.TimerSec)))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			performUpdate(&config, listFn)
		case <-config.UI.updateTriggerChan:
			performUpdate(&config, listFn)
		}
	}
}

func performUpdate(
	config *WorkerConfig,
	listFn func([]string, []any) ([]helpers.ProcessItem, error),
) {
	ui := config.UI
	var listFilters []string
	if !ui.ShowSystem.Load() {
		listFilters = []string{"DB != 'sys'"}
	} else {
		listFilters = []string{}
	}

	if ui.IsRunning.Load() == false {
		status := "Paused"
		listLen := 0

		config.Update(func() {
			ui.StatusBar.SetBorderColor(tcell.ColorYellow)
			UpdateStatusBar(
				ui.StatusBar,
				status,
				listLen,
				ui.TimerSec,
				ui.ShowSystem.Load(),
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
		ui.SQLView.SetText(err.Error())
		ui.IsRunning.Store(false)
		return
	}

	status := "Running"
	listLen := len(itemsList)

	type label struct {
		Name    string
		Content string
	}
	var labels []label

	for _, item := range itemsList {
		if strings.Contains(item.Info.String, "INFORMATION_SCHEMA.PROCESSLIST") {
			listLen--
			continue
		}
		labels = append(labels, label{
			Name: fmt.Sprintf("%d: %s (%ds) from %s@%s - %s",
				item.ID,
				item.DB.String,
				item.Time,
				item.User,
				helpers.HostDropPort(item.Host),
				item.State.String),
			Content: item.Info.String,
		})
	}

	config.Update(func() {
		ui.StatusBar.SetBorderColor(tcell.ColorWhite)
		ui.ListView.Clear()
		for _, l := range labels {
			ui.ListView.AddItem(l.Name, l.Content, 0, nil)
		}
		UpdateStatusBar(
			ui.StatusBar,
			status,
			listLen,
			ui.TimerSec,
			ui.ShowSystem.Load(),
			config.DSN,
			getMemUsage())
	})
}
