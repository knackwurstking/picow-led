package api

import (
	"fmt"
	"sync"
)

var (
	apiColorsMutex = &sync.Mutex{}
)

type APIColors map[string]Color

func NewAPIColors() APIColors {
	return make(map[string]Color, 0)
}

func (ac APIColors) Add(n string, c Color) error {
	defer apiColorsMutex.Unlock()
	apiColorsMutex.Lock()

	if _, ok := ac[n]; ok {
		return fmt.Errorf("color \"%s\" already exists", n)
	}

	ac[n] = c

	return nil
}

func (ac APIColors) Replace(n string, c Color) error {
	defer apiColorsMutex.Unlock()
	apiColorsMutex.Lock()

	if _, ok := ac[n]; !ok {
		return fmt.Errorf("color \"%s\" not exists", n)
	}

	ac[n] = c

	return nil
}

func (ac APIColors) Remove(n string) {
	defer apiColorsMutex.Unlock()
	apiColorsMutex.Lock()

	_, ok := ac[n]
	if ok {
		delete(ac, n)
	}
}

type Color []uint

func NewColor() Color {
	return make(Color, 0)
}
