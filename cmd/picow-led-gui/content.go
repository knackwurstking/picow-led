package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func newContent() *fyne.Container {
	toolbar := widget.NewToolbar()

	aside := container.NewBorder(
		nil, nil, nil, nil,
		// ...
	)

	pages := container.NewBorder(
		nil, nil, nil, nil,
		// ...
	)

	// TODO: Limit drag to resize width somehow
	drawer := container.NewHSplit(aside, pages)
	drawer.SetOffset(0.75)

	return container.New(
		layout.NewBorderLayout(
			toolbar, nil, nil, nil,
		),
		drawer,
	)
}
