package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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

	gui.w.SetContent(widget.NewLabel("Ready to backup?"))

	startBtn := widget.NewButtonWithIcon("Start backup!", theme.DownloadIcon(), func() {
		log.Println("starting backup...")
		gui.w.SetContent(widget.NewLabel("Running..."))
		gui.startCh <- struct{}{}
	})

	gui.w.SetContent(startBtn)

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
