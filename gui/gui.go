package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	app     fyne.App
	w       fyne.Window
	stopCh  chan struct{}
	startCh chan struct{}
}

func New() *GUI {
	return &GUI{
		app:     app.New(),
		stopCh:  make(chan struct{}),
		startCh: make(chan struct{}),
	}
}

func (gui *GUI) Run() {
	gui.w = gui.app.NewWindow("SmugMug Backup")
	gui.w.Resize(fyne.NewSize(800, 600))

	label := widget.NewLabel("Ready to backup?")
	logging := widget.NewMultiLineEntry()

	var startBtn *widget.Button
	startFn := func() {
		log.Println("starting backup...")
		logging.SetText("Starting backup...")
		gui.startCh <- struct{}{}
		if startBtn != nil {
			startBtn.Disable()
		}
	}

	startBtn = widget.NewButtonWithIcon("Start backup!", theme.DownloadIcon(), startFn)

	topSection := container.New(layout.NewCenterLayout(), container.New(layout.NewGridLayout(1), label, startBtn))
	bottomSection := container.NewScroll(logging)
	mainLayout := layout.NewGridLayoutWithRows(2)
	content := container.New(mainLayout, topSection, bottomSection)
	gui.w.SetContent(content)

	gui.w.SetCloseIntercept(gui.closeIntercept)

	gui.w.ShowAndRun()
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
