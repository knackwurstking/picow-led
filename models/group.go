package models

type ColorBinding struct {
	DeviceID DeviceID `json:"device_id"`
	ColorID  ColorID  `json:"color_id"`
}

type GroupSetup []ColorBinding

type Group struct {
	ID        GroupID    `json:"id"`
	Name      string     `json:"name"`
	Devices   []DeviceID `json:"devices"`
	CreatedAt string     `json:"created_at"`
}

func (g *Group) Validate() bool {
	return g.Name != ""
}
