# TODOs

- [x] Implement device types, for now only the colors for RGB can be changed
    - RGB [DEFAULT]
    - RGBW
    - RGBWW
    - W

> ~For now this has some kind of auto W handling for my RGBW devices, this needs to be removed later, see "htmx-devices.go" line 209~

- [x] Fix the color range sliders, values not submitted?

Also remove the zero value, min and max checks on the slider component props

- [ ] Create a device control page
- [ ] Action button for each device, route to the device control page
- [ ] Color setup, brightness control
- [ ] Remove the previous color control from the device edit page
