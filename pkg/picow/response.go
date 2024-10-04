package picow

import "strconv"

type (
	None    any
	Pins    []uint
	Temp    float64
	Version string
)

type DiskUsage struct {
	Used int `json:"used"`
	Free int `json:"free"`
}

type Color []uint

func (c *Color) StringArray() []string {
	cS := make([]string, 0)
	for _, n := range *c {
		cS = append(cS, strconv.Itoa(int(n)))
	}
	return cS
}

type Response[T None | Pins | Color | Temp | Version] struct {
	Data  T      `json:"data"`
	Error string `json:"error"`
	ID    ID     `json:"id"`
}
