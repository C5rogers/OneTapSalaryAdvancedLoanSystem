package payloads

type RatingBreakdown struct {
	CountScore     float64 `json:"count_score"`
	VolumeScore    float64 `json:"volume_score"`
	DurationScore  float64 `json:"duration_score"`
	StabilityScore float64 `json:"stability_score"`
	FinalScore     float64 `json:"final_score"`
}
