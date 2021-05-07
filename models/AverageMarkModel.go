package models

type AverageMark struct {
	AverageMarkID string  `json:"averageMarkID,omitempty"`
	Value         int     `json:"value,omitempty"`
	Subject       Subject `json:"subject,omitempty"`
	StudentID     string  `json:"studentID,omitempty"`
	Grade         Grade   `json:"grade,omitempty"`
	Term          int     `json:"term,omitempty"`
}