package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	app     fyne.App
	w       fyne.Window
	stopCh  chan struct{}
	startCh chan struct{}
	logs    binding.StringList
	info    map[string]string
}

const (
	InfoVersionKey = "version"
)

var UI *GUI

func init() {
	UI = &GUI{
		app:     app.New(),
		stopCh:  make(chan struct{}),
		startCh: make(chan struct{}),
		logs:    binding.NewStringList(),
		info:    make(map[string]string),
	}
}

func (gui *GUI) AddInfo(key, value string) {
	gui.info[key] = value
}

func (gui *GUI) Run() {
	gui.w = gui.app.NewWindow("SmugMug Backup")
	gui.w.Resize(fyne.NewSize(800, 600))

	label := widget.NewLabel("Ready to backup?")

	var startBtn *widget.Button
	startFn := func() {
		gui.logs.Prepend("Starting backup...")
		gui.startCh <- struct{}{}
		if startBtn != nil {
			startBtn.Disable()
		}
	}
	startBtn = widget.NewButtonWithIcon("Start backup!", theme.DownloadIcon(), startFn)
	top := container.New(layout.NewCenterLayout(), container.New(layout.NewGridLayout(1), label, startBtn))

	main := widget.NewListWithData(gui.logs, func() fyne.CanvasObject {
		return widget.NewLabel("template")
	},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	bottom := container.NewBorder(nil, nil, nil, widget.NewLabelWithStyle(fmt.Sprintf("Version: %s", gui.info[InfoVersionKey]), fyne.TextAlignTrailing, fyne.TextStyle{}))

	content := container.NewBorder(top, bottom, nil, nil, main)

	gui.w.SetContent(content)

	gui.w.SetCloseIntercept(gui.closeIntercept)

	gui.w.ShowAndRun()
}

func (gui *GUI) AddLog(line string) {
	gui.logs.Prepend(line)
}

func (gui *GUI) StartBtnTapped() <-chan struct{} {
	return gui.startCh
}

func (gui *GUI) Stop() {
	gui.app.Quit()
}

func (gui *GUI) closeIntercept() {
	gui.w.Close()
	gui.stopCh <- struct{}{}
}

func (gui *GUI) Quit() <-chan struct{} {
	return gui.stopCh
}
