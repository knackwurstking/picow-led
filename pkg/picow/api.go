package picow

type Api struct {
	Devices []*Device `json:"devices"`
}

func NewApi() *Api {
	return &Api{
		Devices: make([]*Device, 0),
	}
}

type Device struct{}
