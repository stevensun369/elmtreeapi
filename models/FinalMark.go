package models

type FinalMark struct {
	FinalMarkID string       `json:"finalMarkID,omitempty" bson:"finalMarkID,omitempty"`
	Value       int          `json:"value,omitempty" bson:"value,omitempty"`
	Subject     ShortSubject `json:"subject,omitempty" bson:"subject,omitempty"`
	StudentID   string       `json:"studentID,omitempty" bson:"studentID,omitempty"`
	Grade       Grade        `json:"grade,omitempty" bson:"grade,omitempty"`
	Term        int          `json:"term,omitempty" bson:"term,omitempty"`
}