package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Content struct {
	widget.BaseWidget
}

func NewContent() *Content {
	return &Content{}
}

func (*Content) CreateRenderer() fyne.WidgetRenderer {
	l := widget.NewLabel("Hello World!")
	c := container.NewBorder(
		nil, nil, nil, nil,
		container.NewCenter(l),
	)

	return widget.NewSimpleRenderer(c)
}
