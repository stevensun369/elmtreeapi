package models

type Truancy struct {
	TruancyID string  `json:"truancyID,omitempty"`
	DateDay   string  `json:"dateDay,omitempty"`
	DateMonth string  `json:"dateMonth,omitempty"`
	Subject   Subject `json:"subject,omitempty"`
	StudentID string  `json:"studentID,omitempty"`
	Grade     Grade   `json:"grade,omitempty"`
	Term      int     `json:"term,omitempty"`
	Motivated bool    `json:"motivated,omitempty"`
}