package domain

type AnalyticsResult struct {
	Metric string  `json:"metric" db:"metric"`
	Value  float64 `json:"value" db:"value"`
}
