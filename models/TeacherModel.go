package models

// the password and the subjectList can be empty
type Teacher struct {
	TeacherID     string    `json:"teacherID,omitempty"`
	FirstName     string    `json:"firstName,omitempty"`
	LastName      string    `json:"lastName,omitempty"`
	CNP           string    `json:"cnp"`
	Password      string    `json:"password,omitempty"`
	HomeroomGrade Grade     `json:"homeroomGrade,omitempty"`
	SubjectList   []Subject `json:"subjectList"`
}
