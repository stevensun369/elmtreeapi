package models

type Subject struct {
	SubjectID string `json:"subjectID,omitempty" bson:"subjectID,omitempty"`
	Name      string `json:"name,omitempty" bson:"name,omitempty"`
	Grade     Grade  `json:"grade,omitempty" bson:"grade,omitempty"`
}
