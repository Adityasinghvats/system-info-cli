package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
	app := tview.NewApplication()

	//new view for showing data
	var cpuView, memView, loadView, uptimeView *tview.TextView
	cpuView = tview.NewTextView()
	memView = tview.NewTextView()
	loadView = tview.NewTextView()
	uptimeView = tview.NewTextView()

	cpuView.SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() { app.Draw() }).
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetTitle(" CPU Usage: ")
	memView.SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() { app.Draw() }).
		SetBorder(true).
		SetBorderColor(tcell.ColorYellow).
		SetTitle(" Memory Usage: ")
	loadView.SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() { app.Draw() }).
		SetBorder(true).
		SetBorderColor(tcell.ColorIsRGB).
		SetTitle(" Processing Load: ")
	uptimeView.SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() { app.Draw() }).
		SetBorder(true).
		SetBorderColor(tcell.ColorWhiteSmoke).
		SetTitle(" Uptime: ")

	//create flex layout
	topView := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(cpuView, 0, 1, false).
		AddItem(memView, 0, 1, false)
	bottomView := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(loadView, 0, 1, false).
		AddItem(uptimeView, 0, 1, false)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topView, 0, 1, false).
		AddItem(bottomView, 0, 1, false)

	//function to get data into texview
	updateInfo := func() {
		//get stats using gopsutil
		cpuPercents, errCpu := cpu.Percent(0, false)
		cpuStr := "Fetching..."
		if errCpu == nil && len(cpuPercents) > 0 {
			color := "green"
			if cpuPercents[0] > 80.0 {
				color = "red"
			} else if cpuPercents[0] > 50.0 {
				color = "yellow"
			}
			cpuStr = fmt.Sprintf("[%s]Overall CPU Usage: %.2f%%[-]", color, cpuPercents[0])
		} else if errCpu != nil {
			cpuStr = fmt.Sprintf("[red]Error: %v[-]", errCpu)
		} else {
			cpuStr = "[orange]CPU Error: No data[-]"
		}
		//memory stats
		vmStat, errMem := mem.VirtualMemory()
		memStr := "Fetching..."
		if errMem == nil {
			memStr = fmt.Sprintf("[blue]Total: %s | Used: %s (%.2f%%)[-]",
				formatBytes(vmStat.Total),
				formatBytes(vmStat.Used),
				vmStat.UsedPercent)
		} else {
			memStr = fmt.Sprintf("[red]Error: %v[-]", errMem)
		}
		//load status
		loadStat, errLoad := load.Avg()
		loadStr := "Fethcing..."
		if errLoad == nil {
			loadStr = fmt.Sprintf("[green]1 min: %.2f | 5 min: %.2f | 15 min: %.2f[-]", loadStat.Load1, loadStat.Load5, loadStat.Load15)
		} else {
			loadStr = fmt.Sprintf("[grey]Not Available or Error: %v[-]", errLoad)
		}
		//uptime
		uptimeSec, errUptime := host.Uptime()
		uptimeStr := "Fetching..."
		if errUptime == nil {
			uptimeDuration := time.Duration(uptimeSec) * time.Second
			// Format duration for better readability (e.g., "1h2m3s")
			uptimeStr = fmt.Sprintf("[blue]%s[-]", uptimeDuration.String())
		} else {
			uptimeStr = fmt.Sprintf("[red]Error: %v[-]", errUptime)
		}
		//update Ui
		app.QueueUpdateDraw(func() {
			cpuView.SetText(cpuStr)
			memView.SetText(memStr)
			loadView.SetText(loadStr)
			uptimeView.SetText(uptimeStr)
		})
	}
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		updateInfo()
		for range ticker.C {
			updateInfo()
		}
	}()
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
