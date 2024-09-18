package api

type API struct {
	Colors  APIColors  `json:"colors"`
	Devices APIDevices `json:"devices"`
}

func NewAPI() *API {
	return &API{
		Devices: make(APIDevices, 0),
		Colors:  make(APIColors, 0),
	}
}
