package models

type Mark struct {
	MarkID    string       `json:"markID,omitempty" bson:"markID,omitempty"`
	Value     int          `json:"value,omitempty" bson:"value,omitempty"`
	DateDay   string       `json:"dateDay,omitempty" bson:"dateDay,omitempty"`
	DateMonth string       `json:"dateMonth,omitempty" bson:"dateMonth,omitempty"`
	Subject   ShortSubject `json:"subject,omitempty" bson:"subject,omitempty"`
	StudentID string       `json:"studentID,omitempty" bson:"studentID,omitempty"`
	Grade     Grade        `json:"grade,omitempty" bson:"grade,omitempty"`
	Term      int          `json:"term,omitempty" bson:"term,omitempty"`
}