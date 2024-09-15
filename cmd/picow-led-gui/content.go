package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func newContent() *fyne.Container {
	l := widget.NewLabel("Hello World!")

	return container.New(
		layout.NewBorderLayout(
			nil, nil, nil, nil,
		),
		container.NewCenter(l),
	)
}
