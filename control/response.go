package control

// Response is a generic response structure used by all control commands.
type Response[T any] struct {
	ID    RequestID `json:"id"`    // ID of the request, used for matching responses.
	Error string    `json:"error"` // Error message if an error occurred during the command execution.
	Data  T         `json:"data"`  // Data returned by the command, or nil if no data is returned.
}

// GetPinsResponse represents the response to a GET request for GPIO pin states.
type GetPinsResponse Response[[]uint8]

// SetPinsResponse represents the response to a SET request for GPIO pin states. No response data is returned.
type SetPinsResponse Response[struct{}]

// GetColorResponse represents the response to a GET request for LED color settings.
type GetColorResponse Response[[]uint8]

// SetColorResponse represents the response to a SET request for LED color settings. No response data is returned.
type SetColorResponse Response[struct{}]

// GetTemperatureResponse represents the response to a GET request for the current temperature.
type GetTemperatureResponse Response[float32]

// GetDiskUsageResponse represents the response to a GET request for disk usage information.
type GetDiskUsageResponse Response[*DiskUsage]

// GetVersionResponse represents the response to a GET request for the firmware version.
type GetVersionResponse Response[string]

// DiskUsage struct contains disk usage statistics.
type DiskUsage struct {
	Total int64 `json:"total"` // Total disk space in bytes.
	Used  int64 `json:"used"`  // Used disk space in bytes.
}
