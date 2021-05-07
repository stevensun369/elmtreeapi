package models

// the password and the subjectList can be empty
type Student struct {
	StudentID   string           `json:"studentID,omitempty"`
	FirstName   string           `json:"firstName,omitempty"`
	DadInitials string           `json:"dadInitials,omitempty"`
	LastName    string           `json:"lastName,omitempty"`
	CNP         string           `json:"cnp,omitempty"`
	Password    string           `json:"password"`
	Grade       Grade            `json:"homeroomGrade"`
	SubjectList []StudentSubject `json:"subjectList"`
}

// we actually define StudentSubject because it doesn't have the grade field.
type StudentSubject struct {
	SubjectID string `json:"subjectID,omitempty"`
	Name      string `json:"name,omitempty"`
}