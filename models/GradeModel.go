package models

type Grade struct {
	GradeID     string `json:"gradeID,omitempty"`
	SchoolID    string `json:"schoolID,omitempty"`
	YearFrom    int    `json:"yearFrom,omitempty"`
	YearTo      int    `json:"yearTo,omitempty"`
	GradeNumber int    `json:"gradeNumber,omitempty"`
	GradeLetter string `json:"gradeLetter,omitempty"`
}