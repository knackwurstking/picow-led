package models

type Color struct {
	ID        ColorID `json:"id"`
	Name      string  `json:"name"`
	Duty      []uint8 `json:"duty"`
	CreatedAt string  `json:"created_at"`
}

func NewColor(name string, duty []uint8) *Color {
	return &Color{
		Name: name,
		Duty: duty,
	}
}

func (c *Color) Validate() bool {
	return c.Name != "" && len(c.Duty) > 0
}
