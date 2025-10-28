package models

// Duty contains rgb like values (0-255)
type Duty []uint8

type Color struct {
	ID        ColorID `json:"id"`
	Name      string  `json:"name"`
	Duty      Duty    `json:"duty"`
	CreatedAt string  `json:"created_at"`
}

func NewColor(name string, duty Duty) *Color {
	return &Color{
		Name: name,
		Duty: duty,
	}
}

func (c *Color) Validate() bool {
	return c.Name != "" && len(c.Duty) > 0
}
