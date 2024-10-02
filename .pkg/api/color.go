package api

type Color []uint

func NewColor() Color {
	return make(Color, 0)
}

type ColorEvent struct {
	Name  string `json:"name"`
	Color Color  `json:"color"`
}

func NewColorEvent(name string, color Color) ColorEvent {
	return ColorEvent{name, color}
}
