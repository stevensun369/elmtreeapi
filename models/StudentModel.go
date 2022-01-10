package models

// the password and the subjectList can be empty
type Student struct {
	StudentID   string           `json:"studentID,omitempty" bson:"studentID,omitempty"`
	FirstName   string           `json:"firstName,omitempty" bson:"firstName,omitempty"`
	DadInitials string           `json:"dadInitials,omitempty" bson:"dadInitials,omitempty"`
	LastName    string           `json:"lastName,omitempty" bson:"lastName,omitempty"`
	CNP         string           `json:"cnp,omitempty" bson:"cnp,omitempty"`
	Password    string           `json:"password" bson:"password"`
	Grade       Grade            `json:"grade" bson:"grade"`
	SubjectList []ShortSubject `json:"subjectList" bson:"subjectList"`
}

// we actually define ShortSubject because it doesn't have the grade field.
type ShortSubject struct {
	SubjectID string `json:"subjectID,omitempty" bson:"subjectID,omitempty"`
	Name      string `json:"name,omitempty" bson:"name,omitempty"`
}
