package models

type Subject struct {
	SubjectID string `json:"subjectID,omitempty"`
	Name      string `json:"name,omitempty"`
	Grade     Grade  `json:"grade,omitempty"`
}
