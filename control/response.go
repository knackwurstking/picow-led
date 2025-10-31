package control

type Response[T any] struct {
	ID    RequestID `json:"id"`
	Error string    `json:"error"`
	Data  T         `json:"data"`
}

type (
	GetPinsResponse Response[[]uint8]
	SetPinsResponse Response[struct{}]

	GetColorResponse Response[[]uint8]
	SetColorResponse Response[struct{}]

	GetTemperatureResponse Response[float32]
	GetDiskUsageResponse   Response[*DiskUsage]
	GetVersionResponse     Response[string]
)

type DiskUsage struct {
	Total int64 `json:"total"`
	Used  int64 `json:"used"`
}
