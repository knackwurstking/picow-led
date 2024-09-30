package api

type Pins []uint

func NewPins() Pins {
	return make(Pins, 0)
}
