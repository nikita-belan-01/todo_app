package http

type checkStatus string

type checkError string

const (
	upCheckStatus   checkStatus = "up"
	downCheckStatus checkStatus = "down"
)

const unavailableCheckError checkError = "unavailable"

type checkResult struct {
	Status checkStatus `json:"status"`
	Error  checkError  `json:"error,omitempty"`
}

type readyStatus string

const (
	okReadyStatus       readyStatus = "ok"
	degradedReadyStatus readyStatus = "degraded"
)

type readyResponse struct {
	Status readyStatus            `json:"status"`
	Checks map[string]checkResult `json:"checks"`
}
