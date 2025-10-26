package micro

const (
	IDNoResponse ID = -1
	IDDefault    ID = 0

	TypeGet Type = "get"
	TypeSet Type = "set"
)

type (
	ID   int
	Type string
)

type Command struct {
	ID      ID     `json:"id"`
	Type    Type   `json:"type"`
	Group   string `json:"group"`
	Command string `json:"command"`

	CommandArgs []string `json:"args"`
}

// Mappings:
//
//	`TypeSet` ("set"):
//		- group: "config" command: "pins"
//	  	- group: "led"    command: "color"
//
//	`TypeGet` ("get"):
//	 	- group: "config" command: "pins"
//	  	- group: "led"    command: "color"
//	  	- group: "info"   command: "temp"
//	  	- group: "info"   command: "disk-usage"
//	  	- group: "info"   command: "version"
//
// Args (optional):
//
//	if `Type` is "set" and `Group` is "config" and `Command` is "pins":
//		[]uint8 - range between 0-28 converted to a slice with strings
//			https://i0.wp.com/randomnerdtutorials.com/wp-content/uploads/2024/02/Raspberry-Pi-Pico-W-RP2040-Rev3-Board-Pinout-GPIOs.png?quality=100&strip=all&ssl=1
//
//	elif `Type` is "set" and `Group` is `led and `Command` is "color":
//		[]uint8 - range between 0-255 converted to a slice with strings
//
//	else
//		nil
//
// Examples:
//
//	`...(IDDefault, "get", "led", "color")`
//	`...(IDDefault, "set", "led", "color")`
func NewCommand(id ID, t Type, group string, command string, args ...string) *Command {
	return &Command{
		ID:          id,
		Type:        t,
		Group:       group,
		Command:     command,
		CommandArgs: args,
	}
}
