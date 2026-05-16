// Package models defines transport and domain DTOs (no business logic).
package models

// MatrixRequest is the payload for POST /api/v1/matrices/qr.
type MatrixRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

// MatrixResponse is returned to the client after orchestration with API 2.
type MatrixResponse struct {
	Original MatrixSummary `json:"original"`
	QR       QRResult        `json:"qr"`
	Stats    interface{}       `json:"stats,omitempty"` // populated from API 2 response
}

// MatrixSummary describes input dimensions (placeholder for richer metadata).
type MatrixSummary struct {
	Rows int `json:"rows"`
	Cols int `json:"cols"`
}

// QRResult holds Q and R factors from decomposition (implementation pending).
type QRResult struct {
	Q [][]float64 `json:"q"`
	R [][]float64 `json:"r"`
}

// API2StatsPayload is sent to API 2 after QR decomposition.
// OriginalMatrix is included so Node can compute stats on the user's input.
type API2StatsPayload struct {
	Q              [][]float64 `json:"q"`
	R              [][]float64 `json:"r"`
	OriginalMatrix [][]float64 `json:"originalMatrix"`
}

// QRResponse is the final response assembled by api-go and returned to the client.
type QRResponse struct {
	OriginalMatrix [][]float64            `json:"originalMatrix"`
	QR             QRResult               `json:"qr"`
	Statistics     map[string]interface{} `json:"statistics"`
}
