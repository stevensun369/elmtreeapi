package models

type Mark struct {
	MarkID    string  `json:"markID,omitempty"`
	Value     int     `json:"value,omitempty"`
	DateDay   string  `json:"dateDay,omitempty"`
	DateMonth string  `json:"dateMonth,omitempty"`
	Subject   Subject `json:"subject,omitempty"`
	StudentID string  `json:"studentID,omitempty"`
	Grade     Grade   `json:"grade,omitempty"`
	Term      int     `json:"term,omitempty"`
}