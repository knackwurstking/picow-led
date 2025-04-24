package picow

import "strconv"

type Pins []uint

func (p *Pins) StringArray() []string {
	pS := make([]string, 0)
	for _, n := range *p {
		pS = append(pS, strconv.Itoa(int(n)))
	}
	return pS
}
