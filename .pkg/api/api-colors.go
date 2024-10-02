package api

import (
	"sync"
)

var (
	apiColorsMutex = &sync.Mutex{}
)

type APIColors map[string]Color

func NewAPIColors() APIColors {
	return make(map[string]Color, 0)
}
