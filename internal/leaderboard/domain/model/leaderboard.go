package model

// Score
type Score struct {
	ClientID  string  `json:"clientId"`
	Score     float64 `json:"score"`
	CreatedAt int64   `json:"createdAt,omitempty"`
}
