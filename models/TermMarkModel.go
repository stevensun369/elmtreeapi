package models

type TermMark struct {
	TermMarkID string `json:"termMarkID,omitempty"`
	Value      int    `json:"value,omitempty"`
	StudentID  string `json:"studentID,omitempty"`
	Grade      Grade  `json:"grade,omitempty"`
	Term       int    `json:"term,omitempty"`
}