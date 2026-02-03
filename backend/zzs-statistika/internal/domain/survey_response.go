package domain

type SurveyResponse struct {
	SurveyID string            `json:"survey_id" db:"survey_id"`
	Answers  map[string]string `json:"answers" db:"answers"`
}
