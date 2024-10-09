package picow

type (
	None    any
	Temp    float64
	Version string
)

type DiskUsage struct {
	Used int `json:"used"`
	Free int `json:"free"`
}

type Response[T None | Pins | Color | Temp | Version] struct {
	Data  T      `json:"data"`
	Error string `json:"error"`
	ID    ID     `json:"id"`
}
