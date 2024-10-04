package picow

import "encoding/json"

type DeviceData_Server struct {
	Name   string `json:"name"`
	Addr   string `json:"addr"`
	Online bool   `json:"online"`
}

type DeviceData struct {
	Server DeviceData_Server `json:"server"`
	Pins   []uint            `json:"pins"`
	Color  []uint            `json:"color"`
}

type Device struct {
	data DeviceData `json:"-"`
}

func (d *Device) MarshalJSON() ([]byte, error) {
	// TODO: Get "pins" and "color" from server

	return json.Marshal(d.data)
}

func (d *Device) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &d.data)
	if err != nil {
		return err
	}

	// TODO: Sync device "pins" and "color"

	return nil
}
