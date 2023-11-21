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
	app       fyne.App
	w         fyne.Window
	stopCh    chan struct{}
	analyzeCh chan struct{}
	startCh   chan struct{}
	logs      binding.StringList
	info      map[string]any
}

const (
	InfoVersionKey    = "version"
	AnalysisDoneKey   = "analysis_done"
	AnalysisAlbumsKey = "analysis_albums"
	AnalysisImagesKey = "analysis_images"
)

var UI *GUI

func init() {
	info := map[string]any{
		InfoVersionKey:    binding.NewString(),
		AnalysisDoneKey:   binding.NewBool(),
		AnalysisAlbumsKey: binding.NewInt(),
		AnalysisImagesKey: binding.NewInt(),
	}

	UI = &GUI{
		app:       app.New(),
		stopCh:    make(chan struct{}),
		analyzeCh: make(chan struct{}),
		startCh:   make(chan struct{}),
		logs:      binding.NewStringList(),
		info:      info,
	}
}

func (gui *GUI) AddInfo(key string, value any) {
	switch value.(type) {
	case string:
		gui.info[key].(binding.String).Set(value.(string))
	case int:
		gui.info[key].(binding.Int).Set(value.(int))
	case bool:
		gui.info[key].(binding.Bool).Set(value.(bool))
	}
}

func (gui *GUI) Run() {
	gui.w = gui.app.NewWindow("SmugMug Backup")
	gui.w.Resize(fyne.NewSize(800, 600))

	gui.w.SetContent(gui.buildLayout())

	gui.w.SetCloseIntercept(gui.closeIntercept)

	gui.w.ShowAndRun()
}

func (gui *GUI) buildLayout() *fyne.Container {
	var startBtn *widget.Button
	startFn := func() {
		gui.logs.Prepend("Starting backup...")
		gui.startCh <- struct{}{}
		if startBtn != nil {
			startBtn.Disable()
		}
	}
	startBtn = widget.NewButtonWithIcon("Start backup!", theme.DownloadIcon(), startFn)

	var analyzeBtn *widget.Button
	analyzeFn := func() {
		gui.logs.Prepend("Analyzing...")
		gui.info[AnalysisDoneKey].(binding.Bool).Set(false)
		gui.analyzeCh <- struct{}{}
		if analyzeBtn != nil {
			analyzeBtn.Disable()
		}

		check := gui.info[AnalysisDoneKey].(binding.Bool)
		check.AddListener(binding.NewDataListener(
			func() {
				if v, _ := check.Get(); v {
					analyzeBtn.Enable()
					check.Set(false)
				}
			},
		))
	}
	analyzeBtn = widget.NewButtonWithIcon("Analyze", theme.DownloadIcon(), analyzeFn)

	top := container.NewHBox(
		container.New(layout.NewCenterLayout(), container.New(
			layout.NewGridLayout(1), widget.NewLabel("Analyze your account"), analyzeBtn),
		),
		container.New(layout.NewCenterLayout(), container.New(
			layout.NewGridLayout(1), widget.NewLabel("Ready to backup?"), startBtn),
		),
	)

	main := widget.NewListWithData(gui.logs, func() fyne.CanvasObject {
		return widget.NewLabel("template")
	},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	// TODO: deal with possibile errors
	nAlbums := gui.info[AnalysisAlbumsKey].(binding.Int)
	nImages := gui.info[AnalysisImagesKey].(binding.Int)
	versionStr, _ := gui.info[InfoVersionKey].(binding.String).Get()

	bottom := container.NewBorder(
		nil,
		nil,
		container.NewHBox(
			widget.NewLabelWithData(binding.IntToStringWithFormat(nAlbums, "Albums: %d")),
			widget.NewLabelWithData(binding.IntToStringWithFormat(nImages, "Images: %d")),
		),
		// widget.NewLabelWithStyle(fmt.Sprintf("Albums: %d, Images: %d", nAlbums, nImages), fyne.TextAlignTrailing, fyne.TextStyle{}),
		widget.NewLabelWithStyle(fmt.Sprintf("Version: %s", versionStr), fyne.TextAlignTrailing, fyne.TextStyle{}),
	)

	return container.NewBorder(top, bottom, nil, nil, main)
}

func (gui *GUI) AddLog(line string) {
	gui.logs.Prepend(line)
}

func (gui *GUI) StartBtnTapped() <-chan struct{} {
	return gui.startCh
}

func (gui *GUI) AnalyzeBtnTapped() <-chan struct{} {
	return gui.analyzeCh
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
