package micro

import "errors"

func GetColor() error {
	pkg := NewCommand(IDDefault, TypeGet, "led", "color")

	// TODO: Send request and wait for response

	return errors.New("under construction")
}

func SetColor() error {
	pkg := NewCommand(IDDefault, TypeSet, "led", "color")

	// TODO: Send request and wait for response

	return errors.New("under construction")
}
