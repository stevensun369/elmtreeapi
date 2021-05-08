package models

type TermMark struct {
	TermMarkID string  `json:"termMarkID,omitempty" bson:"termMarkID,omitempty"`
	Value      float64 `json:"value,omitempty" bson:"value,omitempty"`
	StudentID  string  `json:"studentID,omitempty" bson:"studentID,omitempty"`
	Grade      Grade   `json:"grade,omitempty" bson:"grade,omitempty"`
	Term       int     `json:"term,omitempty" bson:"term,omitempty"`
}