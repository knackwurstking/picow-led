package picow

import "strconv"

type Color []uint

func (c *Color) StringArray() []string {
	cS := make([]string, 0)
	for _, n := range *c {
		cS = append(cS, strconv.Itoa(int(n)))
	}
	return cS
}
