package domain

type AnalyticsResult struct {
	Metric string  `json:"metric" db:"metric"`
	Value  float64 `json:"value" db:"value"`
	Type   string  `json:"type" db:"type"`
	SurveyID string `json:"survey_id" db:"survey_id"`
	Count  int     `json:"count" db:"count"`
}
