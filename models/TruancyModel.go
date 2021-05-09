package models

type Truancy struct {
	TruancyID string       `json:"truancyID,omitempty" bson:"truancyID,omitempty"`
	DateDay   string       `json:"dateDay,omitempty" bson:"dateDay,omitempty"`
	DateMonth string       `json:"dateMonth,omitempty" bson:"dateMonth,omitempty"`
	Subject   ShortSubject `json:"subject,omitempty" bson:"subject,omitempty"`
	StudentID string       `json:"studentID,omitempty" bson:"studentID,omitempty"`
	Grade     Grade        `json:"grade,omitempty" bson:"grade,omitempty"`
	Term      int          `json:"term,omitempty" bson:"term,omitempty"`
	Motivated bool         `json:"motivated,omitempty" bson:"motivated,omitempty"`
}