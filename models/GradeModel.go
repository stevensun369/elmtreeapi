package models

type Grade struct {
	GradeID     string `json:"gradeID,omitempty" bson:"gradeID,omitempty"`
	SchoolID    string `json:"schoolID,omitempty" bson:"schoolID,omitempty"`
	YearFrom    int    `json:"yearFrom,omitempty" bson:"yearFrom,omitempty"`
	YearTo      int    `json:"yearTo,omitempty" bson:"yearTo,omitempty"`
	GradeNumber int    `json:"gradeNumber,omitempty" bson:"gradeNumber,omitempty"`
	GradeLetter string `json:"gradeLetter,omitempty" bson:"gradeLetter,omitempty"`
}