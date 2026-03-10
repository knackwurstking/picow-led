package control

// Response is a generic response structure used by all control commands.
type Response[T any] struct {
	ID    RequestID `json:"id"`    // ID of the request, used for matching responses.
	Error string    `json:"error"` // Error message if an error occurred during the command execution.
	Data  T         `json:"data"`  // Data returned by the command, or nil if no data is returned.
}

// DiskUsage struct contains disk usage statistics.
type DiskUsage struct {
	Total int64 `json:"total"` // Total disk space in bytes.
	Used  int64 `json:"used"`  // Used disk space in bytes.
}
