package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	app    fyne.App
	w      fyne.Window
	stopCh chan struct{}
}

func New() *GUI {
	return &GUI{
		app:    app.New(),
		stopCh: make(chan struct{}),
	}
}

func (gui *GUI) Run() {
	gui.w = gui.app.NewWindow("Hello World")

	gui.w.SetContent(widget.NewLabel("Hello World!"))

	gui.w.SetCloseIntercept(gui.closeIntercept)

	gui.w.ShowAndRun()
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
